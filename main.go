package main

import (
	"log"
	"os"

	"github.com/stormcrows/qdoc/cmd"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
