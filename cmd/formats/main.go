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
	mappedPct  float64
	offsetsPar *termui.Par
	hexPar     *termui.Par
	boxPar     *termui.Par
	asciiPar   *termui.Par
	statsPar   *termui.Par
	boxFooter  *termui.Par
	hexView    = parse.HexViewState{
		BrowseMode: parse.ByGroup,
		RowWidth:   16,
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

	parsers, err := formats.MatchAll(file)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if len(parsers) == 0 {
		panic("no parsers returned, should not happen")
	}

	var chosenParser *parse.ParsedLayout

	if len(parsers) == 1 {

		chosenParser = parsers.First()

	} else if len(parsers) > 1 {

		chosenParser, err = parsers.ChoseOne(file)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}

	fileLayout = *chosenParser
	mappedPct = fileLayout.PercentMapped(fileLayout.FileSize)

	uiLoop(file)
}

func uiLoop(file *os.File) {

	fileLen, _ := file.Seek(0, os.SEEK_END)

	err := termui.Init()
	defer termui.Close()
	if err != nil {
		fmt.Println("termui error:", err)
		return
	}

	calcVisibleRows()
	createUIComponents()
	updateUIPositions()
	focusAtCurrentField()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		hexView.BrowseMode = parse.ByFieldInGroup
		focusAtCurrentField()
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

	termui.Handle("/sys/kbd/o", func(termui.Event) {
		// home. TODO: map to cmd-UP on osx, "HOME" button otherwise
		hexView.StartingRow = 0
		refreshUI(file)
	})

	termui.Handle("/sys/kbd/p", func(termui.Event) {
		// end. TODO: map to cmd-DOWN on osx, "END" button otherwise
		hexView.StartingRow = (fileLen / 16) - int64(hexView.VisibleRows) + 1
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

	termui.Handle("/sys/kbd/?", func(termui.Event) {
		// search
		// XXX show input for search field

		//q := "0x06 0x0E 0x1F"

		// XXX find pos of sequence, and focus on it...
		// XXX for this to work, we need a new browsing mode, which is "browse matches",
		// where right-left button browse between the matches...
		// XXX so, we need to create a ParsedLayout of the matches ..!
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		// XXX resize is bugged on some heights...
		calcVisibleRows()
		updateUIPositions()
		refreshUI(file)
	})

	refreshUI(file)
	termui.Loop() // block until StopLoop is called
}

func calcVisibleRows() {

	hexView.VisibleRows = termui.TermHeight() - 2
}

func createUIComponents() {

	offsetsPar = termui.NewPar("")
	offsetsPar.BorderLeft = false
	offsetsPar.Width = 10
	offsetsPar.BorderLabel = "offset"

	hexPar = termui.NewPar("")
	hexPar.Width = 49
	hexPar.X = 8
	hexPar.Y = 0
	hexPar.BorderLabel = "hex"
	hexPar.BorderFg = termui.ColorCyan

	asciiPar = termui.NewPar("")
	asciiPar.Width = 18
	asciiPar.X = 56
	asciiPar.Y = 0
	asciiPar.BorderRight = false
	asciiPar.TextFgColor = termui.ColorWhite
	asciiPar.BorderLabel = "ascii"
	asciiPar.BorderFg = termui.ColorCyan

	boxPar = termui.NewPar("")
	boxPar.Height = 30
	boxPar.Width = 34
	boxPar.X = 73
	boxPar.TextFgColor = termui.ColorWhite
	boxPar.BorderLabel = fileLayout.TypeSummary()
	boxPar.BorderFg = termui.ColorCyan

	boxFooter = termui.NewPar("")
	boxFooter.Border = false
	boxFooter.Height = 1
	boxFooter.X = 75
	boxFooter.Y = boxPar.Height - 1

	statsPar = termui.NewPar("")
	statsPar.Border = false
	statsPar.Height = 1
	statsPar.X = 9
}

func focusAtCurrentField() {

	var offset int64
	field := fileLayout.Layout[hexView.CurrentGroup]

	switch hexView.BrowseMode {
	case parse.ByGroup:
		offset = field.Offset
	case parse.ByFieldInGroup:
		offset = field.Childs[hexView.CurrentField].Offset
	}

	rowWidth := int64(hexView.RowWidth)
	base := hexView.StartingRow * rowWidth
	ceil := base + int64(hexView.VisibleRows*hexView.RowWidth)

	if offset >= base && offset < ceil {
		// we are in view
		return
	}

	// make scrolling more natural by doing smaller adjustments if possible
	for i := int64(1); i < 10; i++ {
		newOffset := offset + (i * rowWidth)
		if newOffset >= base && newOffset < ceil {
			hexView.StartingRow -= i
			return
		}

		newOffset = offset - (i * rowWidth)
		if newOffset >= base && newOffset < ceil {
			hexView.StartingRow += i
			return
		}
	}

	hexView.StartingRow = int64(offset / rowWidth)
}

func updateUIPositions() {

	statsPar.Y = hexView.VisibleRows + 1
	asciiPar.Height = hexView.VisibleRows + 2
	offsetsPar.Height = hexView.VisibleRows + 2
	hexPar.Height = hexView.VisibleRows + 2
}

func refreshUI(file *os.File) {

	// recalc, to work with resizing of terminal window
	hexView.VisibleRows = termui.TermHeight() - 2

	offsetsPar.Text = fileLayout.PrettyOffsetView(file, hexView)
	hexPar.Text = fileLayout.PrettyHexView(file, hexView)
	asciiPar.Text = fileLayout.PrettyASCIIView(file, hexView)

	boxPar.Text = hexView.CurrentFieldInfo(file, fileLayout)

	if mappedPct < 100.0 {
		boxFooter.Text = fmt.Sprintf("%.1f%%", mappedPct) + " mapped"
	}
	boxFooter.Width = len(boxFooter.Text)

	statsPar.Text = prettyStatString()
	statsPar.Width = len(statsPar.Text)

	termui.Render(offsetsPar, hexPar, asciiPar, boxPar, boxFooter, statsPar)
}

func prettyStatString() string {

	if len(fileLayout.Layout) == 0 {
		return ""
	}

	group := fileLayout.Layout[hexView.CurrentGroup]

	warn := "" // useful while debugging
	end := int64(0)

	// if in sub field view
	if hexView.BrowseMode == parse.ByFieldInGroup {
		field := group.Childs[hexView.CurrentField]

		end = field.Offset + field.Length - 1
		if end >= fileLayout.FileSize {
			warn = " [PAST EOF-f](fg-red)"
		}
		return fmt.Sprintf("selected %d bytes (%x) from %04x to %04x", field.Length, field.Length, field.Offset, end) + warn
	}

	end = group.Offset + group.Length - 1
	if end >= fileLayout.FileSize {
		warn = " [PAST EOF-g](fg-red)"
	}

	return fmt.Sprintf("selected %d bytes (%x) from %04x to %04x", group.Length, group.Length, group.Offset, end) + warn
}
