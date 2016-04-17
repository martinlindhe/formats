package formats

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	parsers = map[string]func(*os.File) (*parse.ParsedLayout, error){
		// archive
		"7z":    parse.SEVENZIP,
		"arj":   parse.ARJ,
		"bzip2": parse.BZIP2,
		"cab":   parse.CAB,
		"gzip":  parse.GZIP,
		"iso":   parse.ISO,
		"lzma":  parse.LZMA,
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
		"otf":   parse.OTF,
		"ttf":   parse.TTF,
		"woff":  parse.WOFF,
		"woff2": parse.WOFF2,

		// exe
		"mz":         parse.MZ,
		"macho":      parse.MachO,
		"lua":        parse.LUA,
		"python":     parse.PYTHON,
		"java class": parse.JAVA,
		"swf":        parse.SWF,

		// bin
		"gba-rom": parse.GBAROM,
		"n64-rom": parse.N64ROM,
		"sqlite3": parse.SQLITE3,
		"pif":     parse.PIF, // windows
		"lnk":     parse.LNK, // windows
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
			parsed.FileName = fileGetName(file)
			parsed.FileSize = fileSize(file)
			return parsed, nil
		}
	}

	raw, _ := parse.RAW(file)
	return raw, nil
}

func fileGetName(file *os.File) string {
	stat, _ := file.Stat()
	return stat.Name()
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	return matchParser(file)
}
