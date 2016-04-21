package formats

import (
	"encoding/binary"
	"os"

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
)

var (
	parsers = map[string]func(*parse.ParseChecker) (*parse.ParsedLayout, error){
		"7z":             archive.SevenZIP,
		"arj":            archive.ARJ,
		"bzip2":          archive.BZIP2,
		"cab":            archive.CAB,
		"deb":            archive.DEB,
		"gzip":           archive.GZIP,
		"iso":            archive.ISO,
		"lzma":           archive.LZMA,
		"rar":            archive.RAR,
		"td2":            archive.TD2,
		"vdi":            archive.VDI,
		"wim":            archive.WIM,
		"winimg":         archive.WINIMG,
		"xz":             archive.XZ,
		"zip":            archive.ZIP,
		"aiff":           av.AIFF,
		"asf":            av.ASF,
		"caf":            av.CAF,
		"flac":           av.FLAC,
		"flv":            av.FLV,
		"midi":           av.MIDI,
		"mkv":            av.MKV,
		"mp3":            av.MP3,
		"mp4":            av.MP4,
		"ogg":            av.OGG,
		"riff":           av.RIFF,
		"dbm":            bin.GnuDBM,
		"gpg":            bin.GnuPG,
		"mo":             bin.GnuMO,
		"mapledb":        bin.MapleDB,
		"sqlite3":        bin.SQLITE3,
		"pdb":            bin.PDB,
		"rdb":            bin.RDB,
		"terminfo":       bin.Terminfo,
		"timezone":       bin.Timezone,
		"xkm":            bin.XKM,
		"chm":            doc.CHM,
		"hlp":            doc.HLP,
		"pdf":            doc.PDF,
		"rtf":            doc.RTF,
		"wri":            doc.WRI,
		"gb-rom":         exe.GameboyROM,
		"gba-rom":        exe.GBAROM,
		"n64-rom":        exe.N64ROM,
		"dex":            exe.DEX,
		"elf":            exe.ELF,
		"java":           exe.JavaClass,
		"lua":            exe.LUA,
		"macho":          exe.MachO,
		"mz":             exe.MZ,
		"python":         exe.PythonBytecode,
		"swf":            exe.SWF,
		"vbe":            exe.VBE,
		"eot":            font.EOT,
		"otf":            font.OTF,
		"pfb":            font.PFB,
		"ps1":            font.PSType1,
		"ttc":            font.TTC,
		"ttf":            font.TTF,
		"woff":           font.WOFF,
		"woff2":          font.WOFF2,
		"x11-snf":        font.X11FontSNF,
		"bmp":            image.BMP,
		"gif":            image.GIF,
		"icns":           image.ICNS,
		"ico":            image.ICO,
		"jpeg":           image.JPEG,
		"pcx":            image.PCX,
		"png":            image.PNG,
		"tiff":           image.TIFF,
		"xcursor":        image.XCursor,
		"bom_store":      macos.BOMStore,
		"bplist":         macos.BPLIST,
		"code_directory": macos.CodeDirectory,
		"ds_store":       macos.DSStore,
		"keychain":       macos.Keychain,
		"dxbc":           windows.DXBC,
		"grp":            windows.GRP,
		"lnk":            windows.LNK,
		"ole-cf":         windows.OLECF,
		"p7x":            windows.P7X,
		"pif":            windows.PIF,
		"pri":            windows.PRI,
		"xbf":            windows.XBF,
		"tga":            image.TGA, // XXX has loose detection, so put it last for now
	}
)

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	return matchParser(file)
}

func matchParser(file *os.File) (*parse.ParsedLayout, error) {

	fileSize, err := fileSize(file)
	if err != nil {
		return nil, err
	}
	maxLen := int64(0xffff)
	len := fileSize
	if len > maxLen {
		len = maxLen
	}

	b := make([]byte, len)

	file.Seek(0, os.SEEK_SET)
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return nil, err
	}

	layout := parse.ParsedLayout{
		FileName: fileGetName(file),
		FileSize: fileSize}

	checker := parse.ParseChecker{
		File:         file,
		ParsedLayout: layout}

	copy(checker.Header[:], b[:len])

	for name, parser := range parsers {

		pl, err2 := parser(&checker)
		if pl != nil {
			if pl.FormatName == "" {
				pl.FormatName = name
			}
			return pl, err2
		}
	}

	// try text detector
	text, err := parse.Text(&checker)
	if err == nil {
		return text, nil
	}

	// fall back to raw
	raw, _ := parse.RAW(&checker)
	return raw, nil
}

func fileGetName(file *os.File) string {
	stat, _ := file.Stat()
	return stat.Name()
}
