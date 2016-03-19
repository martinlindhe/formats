package main

import (
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
	inFile      = kingpin.Arg("file", "Input file").Required().String()
	startingRow = int64(0)
	visibleRows = 10
	rowWidth    = 16
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

	// XXX we fake result from structToFlatStruct() to test presentation
	fileLayout := map[uint64]formats.Layout{
		0x0000:     formats.Layout{2, formats.Uint16le, "magic"}, // XXX also data type
		0x0002:     formats.Layout{4, formats.Uint32le, "width"},
		0x0006:     formats.Layout{4, formats.Uint32le, "height"},
		0x000a:     formats.Layout{9, formats.ASCIIZ, "NAME.EXT"}, // XXX asciiz
		0x000a + 9: formats.Layout{2, formats.Uint16le, "tag"},
	}

	// XXX get console screen height

	uiLoop(&fileLayout, file)
}

func prettyHexView(file *os.File) string {

	reader := io.Reader(file)

	hex := ""

	base := startingRow * int64(rowWidth)
	ceil := base + int64(visibleRows*rowWidth)

	for i := base; i < ceil; i += int64(rowWidth) {

		file.Seek(i, os.SEEK_SET)
		line, err := formats.GetHex(&reader)
		hex += fmt.Sprintf("[[%04x]](fg-yellow) %s\n", i, line)
		if err != nil {
			fmt.Println("got err", err)
			break
		}
	}
	return hex
}

func uiLoop(layout *map[uint64]formats.Layout, file *os.File) {

	fileLen, _ := file.Seek(0, os.SEEK_END)

	hex := prettyHexView(file)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	//hex := "Simple colored text\nwith label. It [can be](fg-red) multilined with \\n or [break automatically](fg-red,fg-bold)"
	/*
	   	hex := `[[0000]](fg-yellow) [0a bb](fg-blue) 2c ff 8e 88 00 0a 01 02 00 ff ff 3f 17 fe
	   [0020] 0a bb 2c ff 8e 88 00 0a 01 02 00 ff ff 3f 17 fe`
	*/
	hexPar := termui.NewPar(hex)
	hexPar.Height = visibleRows + 2
	hexPar.Width = 56
	hexPar.Y = 0
	hexPar.BorderLabel = "Hex"
	// hexPar.BorderFg = termui.ColorYellow

	p := termui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.Y = 15
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = termui.ColorCyan

	// handle key q pressing
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
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

	termui.Render(p, hexPar) // feel free to call Render, it's async and non-block
	termui.Loop()            // block until StopLoop is called

}
