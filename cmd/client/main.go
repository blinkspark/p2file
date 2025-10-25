package main

import (
	"context"
	"flag"
	"log"

	"github.com/blinkspark/p2file"
	"github.com/blinkspark/p2file/config"
	"github.com/libp2p/go-libp2p/core/peer"
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
		id, err := peer.Decode(*config.Channel)
		if err != nil {
			log.Panicln(err)
		}
		pi, err := app.Dht.FindPeer(context.Background(), id)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(pi)
		err = app.Host.Connect(context.Background(), pi)
		if err != nil {
			log.Panicln(err)
		}
	}

	log.Println("Done")
}
