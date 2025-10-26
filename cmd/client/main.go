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
	log.Println("bootstrapped")

	if *config.Channel == "" {
		log.Println("channel: " + app.Host.ID().String())
		err = app.Serve(*config.DirName)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		if *config.IsListing {
			dirList, err := app.ListDir()
			if err != nil {
				log.Panicln(err)
			}
			for _, dir := range dirList {
				log.Println(dir)
			}
		} else if *config.GetFile != "" {
			app.GetFile(*config.GetFile, *config.OutPath)
			// TODO
		}

	}

	log.Println("Done")
}
