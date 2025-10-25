package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/blinkspark/p2file"
)

func main() {
	payload := p2file.Payload{
		Type: p2file.PL_GET_RES_DONE,
	}
	res, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	t := string(res)
	t = t + "\n"
	log.Print(t)
	log.Print(strings.HasSuffix(t, "\n"))
}
