package formats

import (
	"github.com/martinlindhe/formats/parse"
	"github.com/martinlindhe/formats/parse/archive"
	"github.com/martinlindhe/formats/parse/av"
	"github.com/martinlindhe/formats/parse/bin"
	"github.com/martinlindhe/formats/parse/doc"
	"github.com/martinlindhe/formats/parse/exe"
	"github.com/martinlindhe/formats/parse/font"
	"github.com/martinlindhe/formats/parse/image"
	"github.com/martinlindhe/formats/parse/macos"
	"github.com/martinlindhe/formats/parse/windows"
	"os"
)

var (
	parsers = map[string]func(*os.File, [0xffff]byte, parse.ParsedLayout) (*parse.ParsedLayout, error){
		"7z":    archive.SEVENZIP,
		"arj":   archive.ARJ,
		"bzip2": archive.BZIP2,
		"cab":   archive.CAB,
		"gzip":  archive.GZIP,
		"iso":   archive.ISO,
		"lzma":  archive.LZMA,
		"rar":   archive.RAR,
		// "tar": archive.TAR, // XXX
		"td2":      archive.TD2,
		"winimg":   archive.WINIMG,
		"xz":       archive.XZ,
		"zip":      archive.ZIP,
		"aiff":     av.AIFF,
		"asf":      av.ASF,
		"caf":      av.CAF,
		"flac":     av.FLAC,
		"flv":      av.FLV,
		"midi":     av.MIDI,
		"mkv":      av.MKV,
		"mp3":      av.MP3,
		"mp4":      av.MP4,
		"ogg":      av.OGG,
		"riff":     av.RIFF,
		"gba-rom":  bin.GBAROM,
		"n64-rom":  bin.N64ROM,
		"sqlite3":  bin.SQLITE3,
		"pdb":      bin.PDB,
		"chm":      doc.CHM,
		"hlp":      doc.HLP,
		"pdf":      doc.PDF,
		"rtf":      doc.RTF,
		"wri":      doc.WRI,
		"dex":      exe.DEX,
		"elf":      exe.ELF,
		"java":     exe.JAVA, // java class
		"lua":      exe.LUA,
		"macho":    exe.MachO,
		"mz":       exe.MZ,
		"python":   exe.PYTHON,
		"swf":      exe.SWF,
		"vbe":      exe.VBE,
		"eot":      font.EOT,
		"otf":      font.OTF,
		"pfb":      font.PFB,
		"ttc":      font.TTC,
		"ttf":      font.TTF,
		"woff":     font.WOFF,
		"woff2":    font.WOFF2,
		"bmp":      image.BMP,
		"gif":      image.GIF,
		"ico":      image.ICO,
		"jpeg":     image.JPEG,
		"pcx":      image.PCX,
		"png":      image.PNG,
		"tiff":     image.TIFF,
		"ds_store": macos.DSSTORE,
		"pif":      windows.PIF,
		"lnk":      windows.LNK,
		"ole-cf":   windows.OLECF,
	}
)

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	return matchParser(file)
}

func matchParser(file *os.File) (*parse.ParsedLayout, error) {

	var b [0xffff]byte
	// XXX read into arr ...

	for name, parser := range parsers {

		layout := parse.ParsedLayout{
			FormatName: name,
			FileName:   fileGetName(file),
			FileSize:   fileSize(file)}

		pl, err := parser(file, b, layout)
		if pl != nil {
			return pl, err
		}
	}

	raw, _ := parse.RAW(file)
	return raw, nil
}

func fileGetName(file *os.File) string {
	stat, _ := file.Stat()
	return stat.Name()
}
