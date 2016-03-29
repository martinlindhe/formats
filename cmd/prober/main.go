package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/martinlindhe/formats"
)

var (
	inFile = kingpin.Arg("file", "Input file").Required().String()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	f, err := os.Open(*inFile)
	defer f.Close()
	if err != nil {
		fmt.Printf("error reading file: %s\n", err)
		os.Exit(1)
	}

	layout := formats.ParseLayout(f)

	fmt.Println(layout.PrettyPrint())
}
