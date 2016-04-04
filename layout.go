package formats

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	parsers = map[string]func(*os.File) *parse.ParsedLayout{
		"arj": parse.ARJ,
		"bmp": parse.BMP,
		"gif": parse.GIF,
	}
)

func matchParser(file *os.File) *parse.ParsedLayout {
	for name, parse := range parsers {
		parsed := parse(file)
		if parsed != nil {
			parsed.FormatName = name
			parsed.FileSize = fileSize(file)
			return parsed
		}
	}
	return nil
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) *parse.ParsedLayout {

	return matchParser(file)
}
