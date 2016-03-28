package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinlindhe/formats"
)

func main() {

	if len(flag.Args()) < 1 {
		fmt.Printf("Not enough parameters, use -h for usage\n")
		return
	}

	fileName := flag.Args()[0]

	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		fmt.Printf("error reading file: %s\n", err)
		return
	}

	// XXX array with function pointers to probers?

	formats.Probe(f)
}
