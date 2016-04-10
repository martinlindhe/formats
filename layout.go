package formats

import (
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	parsers = map[string]func(*os.File) (*parse.ParsedLayout, error){
		// compression
		"7z":     parse.SEVENZIP,
		"arj":    parse.ARJ,
		"bzip2":  parse.BZIP2,
		"cab":    parse.CAB,
		"gzip":   parse.GZIP,
		"iso":    parse.ISO,
		"td2":    parse.TD2,
		"winimg": parse.WINIMG, // XXX correct name of the format?

		// image
		"bmp":  parse.BMP,
		"gif":  parse.GIF,
		"ico":  parse.ICO,
		"jpeg": parse.JPEG,
		"png":  parse.PNG,
		"tiff": parse.TIFF,
	}
)

func matchParser(file *os.File) (*parse.ParsedLayout, error) {
	for name, parse := range parsers {
		parsed, err := parse(file)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			parsed.FormatName = name
			parsed.FileSize = fileSize(file)
			return parsed, nil
		}
	}
	return nil, fmt.Errorf("no parser found")
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	return matchParser(file)
}
