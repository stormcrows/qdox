package watcher

import (
	"fmt"
	"log"
	"regexp"

	"github.com/radovskyb/watcher"
	"github.com/stormcrows/qdox/pkg/nlp"
)

// FileHandler updates corpus and model on observed folder's changes
func FileHandler(folder string, pattern *regexp.Regexp, c *nlp.Corpus, m *nlp.Model) func(e *watcher.Event) error {
	return func(e *watcher.Event) (err error) {
		if e.IsDir() {
			return
		}

		log.Println(fmt.Sprintf("%s: %s", e.Op, e.Path))

		if err = c.Load(folder, pattern); err != nil {
			return
		}

		if err = m.Train(c); err != nil {
			return
		}

		return
	}
}
