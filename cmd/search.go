package cmd

import (
	"fmt"
	"path"
	"regexp"

	"github.com/urfave/cli"
)

// Search command loads the corpus, trains the model and returns with the results on the terminal
var Search = cli.Command{
	Name:  "search",
	Usage: "qdox search [command options] [folder] [query]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "pattern, P",
			Usage:       "only parse files matching regular expression",
			Destination: &pattern,
			Value:       "\\.txt$",
		},
		cli.IntFlag{
			Name:        "n",
			Usage:       "maximum number of results to return",
			Destination: &n,
			Value:       5,
		},
		cli.Float64Flag{
			Name:        "threshold, t",
			Usage:       "required minimum similarity per document",
			Destination: &threshold,
			Value:       0.3,
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) < 2 {
			fatal(fmt.Errorf("please provide source folder and query"))
		}

		patternr = regexp.MustCompile(pattern)
		folder := path.Clean(c.Args().Get(0))
		query := c.Args().Get(1)

		fatal(corpus.Load(folder, patternr))
		fatal(model.Train(&corpus))

		result := model.Query(query, n, threshold)
		fatal(result.Err)

		for i, v := range result.Matched {
			fmt.Fprintf(c.App.Writer, "%.0f%% %q\n", result.Similarities[i]*100.0, corpus.GetPath(v))
		}
	},
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
