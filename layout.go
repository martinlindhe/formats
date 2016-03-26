package formats

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	parsers = map[string]func(*os.File) []parse.Layout{
		"arj": parse.ARJ,
		"bmp": parse.BMP,
	}
)

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	parsed, err := parseFileByExtension(file)
	if parsed == nil {
		fmt.Println(err)
		panic("XXX if find by extension fails, search all for magic id")
	}

	return parsed, err
}

func parseFileByExtension(
	file *os.File) (*parse.ParsedLayout, error) {

	res := parse.ParsedLayout{
		FileSize: getFileSize(file),
	}

	ext := fileExt(file)

	res.FormatName = "XXX some name"

	if parser, ok := parsers[ext]; ok {
		res.Layout = parser(file)
	} else {
		fmt.Println("error: no match for", ext)
	}

	// XXX

	return &res, nil
}
