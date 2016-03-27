package main

import (
	"fmt"
	"os"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gizak/termui"
	"github.com/martinlindhe/formats"
	"github.com/martinlindhe/formats/parse"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile     = kingpin.Arg("file", "Input file").Required().String()
	fileLayout = parse.ParsedLayout{}
	hexPar     *termui.Par
	boxPar     *termui.Par
	asciiPar   *termui.Par
	helpPar    *termui.Par
	statsPar   *termui.Par
	hexView    = parse.HexViewState{
		StartingRow:  0,
		VisibleRows:  11,
		RowWidth:     16,
		CurrentField: 0,
	}
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	file, err := os.Open(*inFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer file.Close()

	layout := formats.ParseLayout(file)
	fileLayout = *layout

	// XXX get console screen height

	uiLoop(file)
}

func uiLoop(file *os.File) {

	fileLen, _ := file.Seek(0, os.SEEK_END)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hexPar = termui.NewPar("")
	hexPar.Height = hexView.VisibleRows + 2
	hexPar.Width = 56
	hexPar.Y = 0
	hexPar.BorderLabel = "hex"
	hexPar.BorderFg = termui.ColorCyan

	asciiPar = termui.NewPar("")
	asciiPar.Height = hexView.VisibleRows + 2
	asciiPar.Width = 18
	asciiPar.X = 55
	asciiPar.Y = 0
	asciiPar.BorderRight = false
	asciiPar.TextFgColor = termui.ColorWhite
	asciiPar.BorderLabel = "ascii"
	asciiPar.BorderFg = termui.ColorCyan

	boxPar = termui.NewPar(hexView.CurrentFieldInfo(file, fileLayout))
	boxPar.Height = 6
	boxPar.Width = 28
	boxPar.X = 72
	boxPar.TextFgColor = termui.ColorWhite
	boxPar.BorderLabel = fileLayout.FormatName
	boxPar.BorderFg = termui.ColorCyan

	helpPar = termui.NewPar("navigate with arrow keys,\nquit with q")
	helpPar.Height = 8
	helpPar.Width = 28
	helpPar.X = 72
	helpPar.Y = 5
	helpPar.TextFgColor = termui.ColorWhite
	helpPar.BorderLabel = "help"

	statsPar = termui.NewPar("")
	statsPar.Border = false
	statsPar.Height = 1
	statsPar.Width = 50
	statsPar.X = 0
	statsPar.Y = hexView.VisibleRows + 2

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		hexView.Next(len(fileLayout.Layout))
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		hexView.Prev()
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		hexView.StartingRow--
		if hexView.StartingRow < 0 {
			hexView.StartingRow = 0
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		hexView.StartingRow++
		if hexView.StartingRow > (fileLen / 16) {
			hexView.StartingRow = fileLen / 16
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<previous>", func(termui.Event) {
		// pgup jump a whole screen
		hexView.StartingRow -= int64(hexView.VisibleRows)
		if hexView.StartingRow < 0 {
			hexView.StartingRow = 0
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<next>", func(termui.Event) {
		// pgdown, jump a whole screen
		hexView.StartingRow += int64(hexView.VisibleRows)
		if hexView.StartingRow > (fileLen / 16) {
			hexView.StartingRow = fileLen / 16
		}
		refreshUI(file)
	})

	refreshUI(file)
	termui.Loop() // block until StopLoop is called
}

func refreshUI(file *os.File) {

	hexPar.Text = fileLayout.PrettyHexView(file, hexView)
	asciiPar.Text = fileLayout.PrettyASCIIView(file, hexView)
	boxPar.Text = hexView.CurrentFieldInfo(file, fileLayout)
	statsPar.Text = prettyStatString()

	termui.Render(statsPar, hexPar, asciiPar, boxPar, helpPar)
}

func prettyStatString() string {

	field := fileLayout.Layout[hexView.CurrentField]
	return fmt.Sprintf("selected: %d bytes, offset %04x", field.Length, field.Offset)
}
