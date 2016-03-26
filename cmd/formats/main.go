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
	fileLayout = formats.ParsedLayout{}
	hexPar     *termui.Par
	boxPar     *termui.Par
	asciiPar   *termui.Par
	helpPar    *termui.Par
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

	layout, err := formats.ParseLayout(file)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

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
	hexPar.Height = formats.HexView.VisibleRows + 2
	hexPar.Width = 56
	hexPar.Y = 0
	hexPar.BorderLabel = "hex"
	hexPar.BorderFg = termui.ColorCyan

	asciiPar = termui.NewPar("")
	asciiPar.Height = formats.HexView.VisibleRows + 2
	asciiPar.Width = 18
	asciiPar.X = 55
	asciiPar.Y = 0
	asciiPar.TextFgColor = termui.ColorWhite
	asciiPar.BorderLabel = "ascii"
	asciiPar.BorderFg = termui.ColorCyan

	boxPar = termui.NewPar(formats.HexView.CurrentFieldInfo(file, fileLayout))
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

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		formats.HexView.Next(len(fileLayout.Layout))
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		formats.HexView.Prev()
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		formats.HexView.StartingRow--
		if formats.HexView.StartingRow < 0 {
			formats.HexView.StartingRow = 0
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		formats.HexView.StartingRow++
		if formats.HexView.StartingRow > (fileLen / 16) {
			formats.HexView.StartingRow = fileLen / 16
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<previous>", func(termui.Event) {
		// pgup jump a whole screen
		formats.HexView.StartingRow -= int64(formats.HexView.VisibleRows)
		if formats.HexView.StartingRow < 0 {
			formats.HexView.StartingRow = 0
		}
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<next>", func(termui.Event) {
		// pgdown, jump a whole screen
		formats.HexView.StartingRow += int64(formats.HexView.VisibleRows)
		if formats.HexView.StartingRow > (fileLen / 16) {
			formats.HexView.StartingRow = fileLen / 16
		}
		refreshUI(file)
	})

	refreshUI(file)
	termui.Loop() // block until StopLoop is called
}

func refreshUI(file *os.File) {

	hexView := parse.HexViewState{
		StartingRow:  0,
		VisibleRows:  11,
		RowWidth:     16,
		CurrentField: 0,
	}

	hexPar.Text = fileLayout.PrettyHexView(file, hexView)
	asciiPar.Text = fileLayout.PrettyASCIIView(file, hexView)
	boxPar.Text = formats.HexView.CurrentFieldInfo(file, fileLayout)
	termui.Render(hexPar, asciiPar, boxPar, helpPar)
}
