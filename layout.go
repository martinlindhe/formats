package formats

import (
	"github.com/martinlindhe/formats/parse"
	"github.com/martinlindhe/formats/parse/archive"
	"os"
)

var (
	parsers = map[string]func(*os.File) (*parse.ParsedLayout, error){
		// archive
		"7z":    archive.SEVENZIP,
		"arj":   archive.ARJ,
		"bzip2": archive.BZIP2,
		"cab":   archive.CAB,
		"gzip":  archive.GZIP,
		"iso":   archive.ISO,
		"lzma":  archive.LZMA,
		"rar":   archive.RAR,
		// "tar": archive.TAR, // XXX
		"td2":    archive.TD2,
		"winimg": archive.WINIMG,
		"xz":     archive.XZ,
		"zip":    archive.ZIP,

		// image
		"bmp":  parse.BMP,
		"gif":  parse.GIF,
		"ico":  parse.ICO,
		"jpeg": parse.JPEG,
		"pcx":  parse.PCX,
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
		"chm": parse.CHM,
		"hlp": parse.HLP,
		"pdf": parse.PDF,
		"rtf": parse.RTF,
		"wri": parse.WRI,

		// font
		"eot":   parse.EOT,
		"otf":   parse.OTF,
		"pfb":   parse.PFB,
		"ttc":   parse.TTC,
		"ttf":   parse.TTF,
		"woff":  parse.WOFF,
		"woff2": parse.WOFF2,

		// exe
		"dex":        parse.DEX,
		"elf":        parse.ELF,
		"java class": parse.JAVA,
		"lua":        parse.LUA,
		"macho":      parse.MachO,
		"mz":         parse.MZ,
		"python":     parse.PYTHON,
		"swf":        parse.SWF,
		"vbe":        parse.VBE,

		// bin
		"gba-rom": parse.GBAROM,
		"n64-rom": parse.N64ROM,
		"sqlite3": parse.SQLITE3,
		"pdb":     parse.PDB, // visual studio debug info

		// os-windows
		"pif":    parse.PIF,
		"lnk":    parse.LNK,
		"ole-cf": parse.OLECF,

		// os-macos
		"ds_store": parse.DSSTORE,
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
