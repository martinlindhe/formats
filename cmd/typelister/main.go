package main

import (
	"fmt"
	"os"
	"path/filepath"

	//	"github.com/martinlindhe/formats"
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

	fileList := []os.FileInfo{}
	err := filepath.Walk(*path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, f)
		return nil
	})
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		fmt.Println(file.Name(), ": xxx check format!")
	}

}
