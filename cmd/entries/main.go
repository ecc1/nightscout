package main

import (
	"log"

	"github.com/ecc1/nightscout"
)

func main() {
	entries, err := nightscout.DownloadEntries(10)
	if err != nil {
		log.Fatal(err)
	}
	entries.Print()
}
