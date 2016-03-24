package main

import (
	"fmt"
	"os"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gizak/termui"
	"github.com/martinlindhe/formats"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile     = kingpin.Arg("file", "Input file").Required().String()
	fileLayout = formats.ParsedLayout{}
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

	hex := fileLayout.PrettyHexView(file)

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hexPar := termui.NewPar(hex)
	hexPar.Height = formats.HexView.VisibleRows + 2
	hexPar.Width = 56
	hexPar.Y = 0
	hexPar.BorderLabel = "hex"
	hexPar.BorderFg = termui.ColorCyan

	boxText := formats.HexView.CurrentFieldInfo(file, fileLayout)
	box := termui.NewPar(boxText)
	box.Height = 8
	box.Width = 30
	box.X = 56
	box.TextFgColor = termui.ColorWhite
	box.BorderLabel = fileLayout.FormatName
	box.BorderFg = termui.ColorCyan

	help := termui.NewPar("navigate with arrow keys,\nquit with q")
	help.Height = 5
	help.Width = 30
	help.X = 56
	help.Y = 8
	help.TextFgColor = termui.ColorWhite
	help.BorderLabel = "help"
	//help.BorderFg = termui.ColorCyan

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		formats.HexView.Next(len(fileLayout.Layout))
		hexPar.Text = fileLayout.PrettyHexView(file)
		box.Text = formats.HexView.CurrentFieldInfo(file, fileLayout)
		termui.Render(hexPar, box)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		formats.HexView.Prev()
		hexPar.Text = fileLayout.PrettyHexView(file)
		box.Text = formats.HexView.CurrentFieldInfo(file, fileLayout)
		termui.Render(hexPar, box)
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		formats.HexView.StartingRow--
		if formats.HexView.StartingRow < 0 {
			formats.HexView.StartingRow = 0
		}
		hexPar.Text = fileLayout.PrettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		formats.HexView.StartingRow++
		if formats.HexView.StartingRow > (fileLen / 16) {
			formats.HexView.StartingRow = fileLen / 16
		}
		hexPar.Text = fileLayout.PrettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<previous>", func(termui.Event) {
		// pgup jump a whole screen
		formats.HexView.StartingRow -= int64(formats.HexView.VisibleRows)
		if formats.HexView.StartingRow < 0 {
			formats.HexView.StartingRow = 0
		}
		hexPar.Text = fileLayout.PrettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Handle("/sys/kbd/<next>", func(termui.Event) {
		// pgdown, jump a whole screen
		formats.HexView.StartingRow += int64(formats.HexView.VisibleRows)
		if formats.HexView.StartingRow > (fileLen / 16) {
			formats.HexView.StartingRow = fileLen / 16
		}
		hexPar.Text = fileLayout.PrettyHexView(file)
		termui.Render(hexPar)
	})

	termui.Render(help, hexPar, box)

	termui.Loop() // block until StopLoop is called
}
