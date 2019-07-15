package cmd

import (
	"github.com/urfave/cli"
)

// NewApp initializes the app
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "qdoc"
	app.Usage = "query documents for given phrases"
	app.UsageText = "qdoc [global options] command [command options] [arguments...]"
	app.Author = "Stormcrows"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{Search, Serve}

	return app
}
