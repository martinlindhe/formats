package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/martinlindhe/formats"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	path  = kingpin.Arg("path", "Input path").Required().String()
	short = kingpin.Flag("short", "Short mode").Short('s').Bool()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	err := filepath.Walk(*path, func(fileName string, f os.FileInfo, err error) error {

		file, err := os.Open(fileName)
		if err != nil {
			return nil
		}
		fi, _ := file.Stat()
		if fi.IsDir() {
			return nil
		}

		fmt.Printf("%s:", fileName)

		layout, err := formats.ParseLayout(file)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		if layout == nil {
			fmt.Println("ERR: layout is nil")
			os.Exit(1)
		}

		fmt.Printf("%s\n", layout.ShortPrint())

		return nil
	})
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}
}
