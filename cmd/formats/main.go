package main

import (
	"fmt"
	"io"
	"os"

	"github.com/martinlindhe/formats"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile = kingpin.Arg("file", "Input file").Required().String()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	formats.Formatting(formats.HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      2,
	})

	file, _ := os.Open(*inFile)
	defer file.Close()

	reader := io.Reader(file)

	// XXX get console screen height
	hex, _ := formats.GetHex(&reader, 3)

	fmt.Println(hex)
}
