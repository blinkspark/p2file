package main

import (
	"flag"
	"log"

	"github.com/blinkspark/p2file"
	"github.com/blinkspark/p2file/config"
)

func main() {
	flag.Parse()

	app, err := p2file.NewApp()
	if err != nil {
		log.Panicln(err)
	}
	defer app.Close()

	err = app.WaitBootstrap(0)
	if err != nil {
		log.Panicln(err)
	}

	if *config.Channel == "" {
		err = app.Serve(*config.DirName)
		if err != nil {
			log.Panicln(err)
		}
	} else {

	}

	log.Println("Done")
}
