package cmd

import (
	"regexp"

	"github.com/stormcrows/qdox/pkg/nlp"
)

var (
	port           = 8080
	corpus         = nlp.NewCorpus()
	model          = nlp.NewLSIModel()
	n              = 5
	threshold      = 0.3
	serveFiles     = false
	interact       = false
	watcherEnabled = false
	interval       = int64(1000)
	pattern        = "\\.txt$"
	patternr       = regexp.MustCompile(pattern)
)
