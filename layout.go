package formats

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	parsers = map[string]func(*os.File) *parse.ParsedLayout{
		"arj": parse.ARJ,
		"bmp": parse.BMP,
	}
)

func matchParser(file *os.File) *parse.ParsedLayout {
	for name, parse := range parsers {
		x := parse(file)
		if x != nil {
			fmt.Println("XXX matched", name)
			return x
		}
	}
	return nil
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) *parse.ParsedLayout {

	return matchParser(file)
}
