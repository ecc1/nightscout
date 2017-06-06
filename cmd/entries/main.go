package main

import (
	"fmt"
	"log"

	"github.com/ecc1/nightscout"
)

func main() {
	nightscout.SetVerbose(true)
	var entries []nightscout.Entry
	err := nightscout.Get("entries", &entries)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nightscout.JSON(entries))
}
