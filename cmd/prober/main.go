package main

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile = kingpin.Arg("file", "Input file").Required().String()
	short  = kingpin.Flag("short", "Short mode").Short('s').Bool()
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

	layout, err := formats.ParseLayout(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *short {
		fmt.Println(layout.ShortPrint())
		os.Exit(0)
	}

	fmt.Println(layout.PrettyPrint())
}
