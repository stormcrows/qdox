package main

import (
	"log"
	"os"

	"github.com/stormcrows/qdox/cmd"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
