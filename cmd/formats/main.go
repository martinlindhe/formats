package main

import (
	"fmt"
	"os"

	"github.com/gizak/termui"
	"github.com/martinlindhe/formats"
	"github.com/martinlindhe/formats/parse"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile     = kingpin.Arg("file", "Input file").Required().String()
	fileLayout = parse.ParsedLayout{}
	offsetsPar *termui.Par
	hexPar     *termui.Par
	boxPar     *termui.Par
	asciiPar   *termui.Par
	statsPar   *termui.Par
	hexView    = parse.HexViewState{
		BrowseMode:   parse.ByGroup,
		StartingRow:  0,
		VisibleRows:  11,
		RowWidth:     16,
		CurrentGroup: 0,
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
	if layout == nil {
		fmt.Println("error: file not recognized")
		os.Exit(1)
	}

	fileLayout = *layout

	uiLoop(file)
}

func uiLoop(file *os.File) {

	fileLen, _ := file.Seek(0, os.SEEK_END)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hexView.VisibleRows = termui.TermHeight() - 3

	offsetsPar = termui.NewPar("")
	offsetsPar.BorderLeft = false
	offsetsPar.Width = 10
	offsetsPar.Height = hexView.VisibleRows + 2
	offsetsPar.BorderLabel = "offset"

	hexPar = termui.NewPar("")
	hexPar.Height = hexView.VisibleRows + 2
	hexPar.Width = 49
	hexPar.X = 8
	hexPar.Y = 0
	hexPar.BorderLabel = "hex"
	hexPar.BorderFg = termui.ColorCyan

	asciiPar = termui.NewPar("")
	asciiPar.Height = hexView.VisibleRows + 2
	asciiPar.Width = 18
	asciiPar.X = 56
	asciiPar.Y = 0
	asciiPar.BorderRight = false
	asciiPar.TextFgColor = termui.ColorWhite
	asciiPar.BorderLabel = "ascii"
	asciiPar.BorderFg = termui.ColorCyan

	boxPar = termui.NewPar("")
	boxPar.Height = 8
	boxPar.Width = 28
	boxPar.X = 73
	boxPar.TextFgColor = termui.ColorWhite
	boxPar.BorderLabel = fileLayout.FormatName
	boxPar.BorderFg = termui.ColorCyan

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

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		hexView.BrowseMode = parse.ByFieldInGroup
		//		hexView.CurrentField = 0
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<escape>", func(termui.Event) {
		hexView.BrowseMode = parse.ByGroup
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		switch hexView.BrowseMode {
		case parse.ByGroup:
			hexView.NextGroup(fileLayout.Layout)

		case parse.ByFieldInGroup:
			hexView.NextFieldInGroup(fileLayout.Layout)
		}
		focusAtCurrentField()
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		switch hexView.BrowseMode {
		case parse.ByGroup:
			hexView.PrevGroup()

		case parse.ByFieldInGroup:
			hexView.PrevFieldInGroup()
		}
		focusAtCurrentField()
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

func focusAtCurrentField() {

	// get current field offset
	offset := fileLayout.Layout[hexView.CurrentGroup].Offset

	base := hexView.StartingRow * int64(hexView.RowWidth)
	ceil := base + int64(hexView.VisibleRows*hexView.RowWidth)

	// see if it is view with current row selection
	if offset >= base && offset <= ceil {
		return
	}

	// if not, change current row
	hexView.StartingRow = int64(offset / 16)
}

func refreshUI(file *os.File) {

	// recalc, to work with resizing of terminal window
	hexView.VisibleRows = termui.TermHeight() - 3

	offsetsPar.Text = fileLayout.PrettyOffsetView(file, hexView)
	hexPar.Text = fileLayout.PrettyHexView(file, hexView)
	asciiPar.Text = fileLayout.PrettyASCIIView(file, hexView)
	boxPar.Text = hexView.CurrentFieldInfo(file, fileLayout)
	statsPar.Text = prettyStatString()

	termui.Render(offsetsPar, statsPar, hexPar, asciiPar, boxPar)
}

func prettyStatString() string {

	group := fileLayout.Layout[hexView.CurrentGroup]

	// if in sub field view
	if hexView.BrowseMode == parse.ByFieldInGroup {
		field := group.Childs[hexView.CurrentField]
		return fmt.Sprintf("selected: %d bytes (%x), offset %04x", field.Length, field.Length, field.Offset)
	}

	return fmt.Sprintf("selected: %d bytes (%x), offset %04x", group.Length, group.Length, group.Offset)
}
