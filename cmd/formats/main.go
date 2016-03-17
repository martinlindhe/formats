package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	ui "github.com/gizak/termui"
	"github.com/martinlindhe/arj"
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

	// ---

	// extract arj struct
	arj, err := arj.ParseARJArchive(file)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	res := arjStructToFlatStruct(&arj)

	spew.Dump(res)

	/*
		reader := io.Reader(file)

		// rewind
		file.Seek(0, os.SEEK_SET)

		// XXX get console screen height
		hex, _ := formats.GetHex(&reader, 3)
		fmt.Println(hex)
	*/

	// uiLoop()
}

func arjStructToFlatStruct(t *arj.Arj) []string { // XXX figure out return type

	res := []string{}

	//	spew.Dump(x)

	// XXX iterate over struct, create a 2d rep of the structure mapping

	s := reflect.ValueOf(t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	return res
}

func uiLoop() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Gauge"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	// handle key q pressing
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})

	ui.Render(p, g) // feel free to call Render, it's async and non-block
	ui.Loop()       // block until StopLoop is called
}
