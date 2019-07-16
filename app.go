package main

import (
	"github.com/stormcrows/qdox/cmd"
	"github.com/urfave/cli"
)

// NewApp initializes the app
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "qdox"
	app.Usage = "query documents for given phrases"
	app.UsageText = "qdox [global options] command [command options] [arguments...]"
	app.Author = "Stormcrows"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{cmd.Search, cmd.Serve}

	return app
}
