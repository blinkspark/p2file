package main

import (
	"log"

	"github.com/blinkspark/p2file"
)

func main() {
	app, err := p2file.NewApp(nil)
	if err != nil {
		log.Panicln(err)
	}
	defer app.Close()

	err = app.WaitBootstrap(0)
	if err != nil {
		log.Panicln(err)
	}

	err = app.Serve()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Done")
}
