package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/nightscout"
)

var (
	verbose = flag.Bool("v", false, "verbose mode")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] glucose.json\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	entries, err := nightscout.ReadEntries(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	wrong := 0
	for i, e := range entries {
		if e.Type != nightscout.SGVType {
			continue
		}
		trend := nightscout.Trend(entries[i:])
		if trend != e.Direction {
			wrong++
			if *verbose {
				fmt.Printf("%s  %-13s  %-13s\n", e.Time().Format(time.Stamp), trend, e.Direction)
			}
		}
		total++
	}
	fmt.Printf("%d / %d wrong (%d%% correct)\n", wrong, total, 100*(total-wrong)/total)
}
