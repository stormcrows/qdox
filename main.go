package main

import (
	"log"
	"os"
)

func main() {
	if err := NewApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
