package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gizak/termui"
	// "github.com/martinlindhe/arj"
	"github.com/martinlindhe/formats"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile     = kingpin.Arg("file", "Input file").Required().String()
	fileLayout = []formats.Layout{}
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	formats.Formatting(formats.HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      1,
	})

	file, _ := os.Open(*inFile)
	defer file.Close()

	fileLayout = formats.ParseLayout(file)

	// ---

	/*
		// extract arj struct
		arj, err := arj.ParseARJArchive(file)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		res := structToFlatStruct(&arj)
	*/

	// XXX get console screen height

	uiLoop(&fileLayout, file)
}

func uiLoop(layout *[]formats.Layout, file *os.File) {

	fileLen, _ := file.Seek(0, os.SEEK_END)

	hex := prettyHexView(file)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hexPar := termui.NewPar(hex)
	hexPar.Height = visibleRows + 2
	hexPar.Width = 56
	hexPar.Y = 0
	hexPar.BorderLabel = "Hex"
	// hexPar.BorderFg = termui.ColorYellow

	box := termui.NewPar("info box")
	box.Height = 8
	box.Width = 40
	box.X = 60
	box.TextFgColor = termui.ColorWhite
	box.BorderLabel = "info"
	box.BorderFg = termui.ColorCyan

	p := termui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.Y = 15
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = termui.ColorCyan

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		currentField++
		if currentField >= len(fileLayout) {
			currentField = len(fileLayout) - 1
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		currentField--
		if currentField < 0 {
			currentField = 0
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		startingRow--
		if startingRow < 0 {
			startingRow = 0
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		startingRow++
		if startingRow > (fileLen / 16) {
			startingRow = fileLen / 16
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<previous>", func(termui.Event) {
		// pgup jump a whole screen
		startingRow -= int64(visibleRows)
		if startingRow < 0 {
			startingRow = 0
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<next>", func(termui.Event) {
		// pgdown, jump a whole screen
		startingRow += int64(visibleRows)
		if startingRow > (fileLen / 16) {
			startingRow = fileLen / 16
		}
		hexPar.Text = prettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Render(p, hexPar, box)

	termui.Loop() // block until StopLoop is called
}
