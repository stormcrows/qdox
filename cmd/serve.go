package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"regexp"
	"syscall"
	"text/template"
	"time"

	"github.com/buaazp/fasthttprouter"

	"github.com/stormcrows/qdox/internal"

	"github.com/stormcrows/qdox/pkg/watcher"

	"github.com/urfave/cli"
	"github.com/valyala/fasthttp"
)

// Result defines a single search result
type Result struct {
	Path       string
	Similarity string
}

// QueryResponse is JSON response to /query requests
type QueryResponse struct {
	Query   string
	Results []Result
}

// Tpl holds compiled templates for execution
var Tpl *template.Template

var (
	nSelection      = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tSelection      = []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	defaultResponse = map[string]interface{}{"NSelection": nSelection, "TSelection": tSelection}
)

// Serve command trains on corpus from the provided folder and then serves /query requests via http
var Serve = cli.Command{
	Name:  "serve",
	Usage: "qdox serve [command options] [folder]",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "starts serving at given port",
			Destination: &port,
			Value:       8080,
		},
		cli.StringFlag{
			Name:        "pattern, P",
			Usage:       "parse files matching given regexp pattern",
			Destination: &pattern,
			Value:       "\\.txt$",
		},
		cli.BoolFlag{
			Name:        "watcher, w",
			Usage:       "updates model on observed folder's change",
			Destination: &watcherEnabled,
		},
		cli.Int64Flag{
			Name:        "watcher-interval, wi",
			Usage:       "folder update check interval in ms",
			Destination: &interval,
			Value:       int64(1000),
		},
		cli.BoolFlag{
			Name:        "interact, i",
			Usage:       "simple query ui served at /index level",
			Destination: &interact,
		},
	},
	Action: func(c *cli.Context) (err error) {
		// args
		if len(c.Args()) < 1 {
			return fmt.Errorf("please provide folder path")
		}

		if Tpl == nil {
			Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
		}
		patternr = regexp.MustCompile(pattern)
		folder := c.Args().Get(0)

		// nlp
		corpus.Load(folder, patternr)
		err = model.Train(&corpus)
		if err != nil {
			panic(err)
		}

		// watcher
		if watcherEnabled {
			watcher := &watcher.Watcher{
				MaxEvents: 10,
				Handler:   internal.FileHandler(folder, patternr, &corpus, model),
				Folder:    folder,
				Interval:  time.Millisecond * time.Duration(interval),
				Pattern:   patternr,
			}

			done := make(chan bool, 1)
			go func() {
				err = watcher.Watch(done)
				if err != nil {
					panic(err)
				}
			}()

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				defer func() {
					close(sigs)
					close(done)
				}()

				for {
					select {
					case <-done:
						fmt.Println("done")
						return
					case sig := <-sigs:
						fmt.Println(sig)
						watcher.Stop()
					case <-time.After(time.Second):
						continue
					}
				}
			}()
		}

		// router
		router := fasthttprouter.New()
		router.ServeFiles("/static/*filepath", folder)
		router.GET("/query", QueryHandler)
		if interact {
			router.GET("/", IndexHandler)
		}

		// serve
		fmt.Printf("qdox listening on port: %d\n", port)
		return fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), router.Handler)
	},
}

// IndexHandler displays template in interaction mode
func IndexHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	Tpl.ExecuteTemplate(ctx, "interaction.gohtml", defaultResponse)
}

// QueryHandler handles search queries and responds with JSON
func QueryHandler(ctx *fasthttp.RequestCtx) {
	// args
	query := string(ctx.QueryArgs().Peek("q"))
	if query == "" {
		respondErr(400, "q param is missing, usage: /query?q=find%20me%somebody%20to%20love", ctx)
		return
	}

	n := 5
	if ctx.QueryArgs().GetUintOrZero("n") > 0 {
		n = ctx.QueryArgs().GetUintOrZero("n")
	}

	threshold := 0.3
	if ctx.QueryArgs().GetUfloatOrZero("threshold") > 0.0 {
		threshold = ctx.QueryArgs().GetUfloatOrZero("threshold")
	}

	// nlp query
	result := model.Query(query, n, threshold)
	if result.Err != nil {
		log.Println(fmt.Sprintf("query error: %s", result.Err))
		respondErr(500, "", ctx)
		return
	}

	// response
	resp := QueryResponse{query, make([]Result, len(result.Matched))}

	for i, v := range result.Matched {
		similarity := fmt.Sprintf("%.0f", result.Similarities[i]*100.0)
		path := fmt.Sprintf("static/%s", path.Base(model.Corpus.GetPath(v)))
		resp.Results[i] = Result{path, similarity}
	}

	jresp, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf("json marshalling error: %s", result.Err))
		respondErr(500, "", ctx)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	fmt.Fprintln(ctx, string(jresp))
}

func respondErr(statusCode int, body string, ctx *fasthttp.RequestCtx) {
	if body == "" {
		ctx.Error(fasthttp.StatusMessage(statusCode), statusCode)
	} else {
		ctx.Error(body, statusCode)
	}
}
