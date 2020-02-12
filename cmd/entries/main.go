package main

import (
	"log"

	"github.com/ecc1/nightscout"
)

func main() {
	site, err := nightscout.DefaultSite()
	if err != nil {
		log.Fatal(err)
	}
	entries, err := site.DownloadEntries(10)
	if err != nil {
		log.Fatal(err)
	}
	entries.Print()
}
