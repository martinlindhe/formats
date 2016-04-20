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

	fileList := []string{}
	err := filepath.Walk(*path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	for _, fileName := range fileList {

		fmt.Printf("%s:", fileName)

		file, _ := os.Open(fileName)
		fi, _ := file.Stat()
		if fi.IsDir() {
			continue
		}

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
	}
}
