package formats

import (
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	parsers = map[string]func(*os.File) (*parse.ParsedLayout, error){
		// compression
		"7z":    parse.SEVENZIP,
		"arj":   parse.ARJ,
		"bzip2": parse.BZIP2,
		"cab":   parse.CAB,
		"gzip":  parse.GZIP,
		"iso":   parse.ISO,
		"rar":   parse.RAR,
		// "tar": parse.TAR, // XXX
		"td2":    parse.TD2,
		"winimg": parse.WINIMG,
		"xz":     parse.XZ,
		"zip":    parse.ZIP,

		// image
		"bmp":  parse.BMP,
		"gif":  parse.GIF,
		"ico":  parse.ICO,
		"jpeg": parse.JPEG,
		"png":  parse.PNG,
		"tiff": parse.TIFF,

		// a/v
		"aiff": parse.AIFF,
		"asf":  parse.ASF,
		"caf":  parse.CAF,
		"flac": parse.FLAC,
		"flv":  parse.FLV,
		"midi": parse.MIDI,
		"mkv":  parse.MKV,
		"mp3":  parse.MP3,
		"mp4":  parse.MP4,
		"ogg":  parse.OGG,
		"riff": parse.RIFF,

		// doc
		"chm":  parse.CHM,
		"hlp":  parse.HLP,
		"pdf":  parse.PDF,
		"rtf":  parse.RTF,
		"word": parse.WORD,
		"wri":  parse.WRI,

		// font
		"otf": parse.OTF,
		"ttf": parse.TTF,

		// exe
		"mz":     parse.MZ,
		"lua":    parse.LUA,
		"python": parse.PYTHON,

		// windows
		"pif": parse.PIF,
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
