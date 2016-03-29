package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/martinlindhe/formats"
	"github.com/martinlindhe/formats/parse"
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

	fmt.Println(prettyLayout(layout))
}

func prettyLayout(parsedLayout *parse.ParsedLayout) string { // XXX move

	res := ""
	for _, layout := range parsedLayout.Layout {
		res += layout.Info + fmt.Sprintf(" (%04x)", layout.Offset) + ", " + layout.Type.String() + "\n"
		// XXX childs
		for _, child := range layout.Childs {
			res += "  " + child.Info + fmt.Sprintf(" (%04x)", child.Offset) + ", " + child.Type.String() + "\n"
		}
	}

	return res
}
