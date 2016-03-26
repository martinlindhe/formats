package formats

import (
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
		parsed := parse(file)
		if parsed != nil {
			parsed.FormatName = name
			parsed.FileSize = getFileSize(file)
			return parsed
		}
	}
	return nil
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) *parse.ParsedLayout {

	return matchParser(file)
}
