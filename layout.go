package formats

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

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

type Parser func(*parse.ParseChecker) (*parse.ParsedLayout, error)

var (
	parsers = map[string]Parser{
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
		"tga":            image.TGA, // XXX has too loose detection
		"tiff":           image.TIFF,
		"xcursor":        image.XCursor,
		"bom_store":      macos.BOMStore,
		"bplist":         macos.BPLIST,
		"code_directory": macos.CodeDirectory,
		"ds_store":       macos.DSStore,
		"keychain":       macos.Keychain,
		"ari8":           windows.ARI8,
		"dxbc":           windows.DXBC,
		"grp":            windows.GRP,
		"hwrs":           windows.HWRS,
		"lnk":            windows.LNK,
		"ole-cf":         windows.OLECF,
		"p7x":            windows.P7X,
		"pif":            windows.PIF,
		"pri":            windows.PRI,
		"rbs":            windows.RBS,
		"rsrc":           windows.RSRC,
		"uce":            windows.UCE,
		"xbf":            windows.XBF,
	}
)

type MatchingParsers map[string]Parser

func (mp *MatchingParsers) First() Parser {

	for _, parser := range *mp {
		return parser
	}
	return nil
}

func (mp *MatchingParsers) ChoseOne() (Parser, error) {

	i := 1
	fmt.Println("multiple parsers matched input file, please choose one:\n")
	for name, _ := range parsers {
		fmt.Printf("%d: %s\n", i, name)
		i++
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	u, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid input")
	}

	i = 1
	for _, parser := range parsers {
		if i == int(u) {
			return parser, nil
		}
		i++
	}

	return nil, fmt.Errorf("selection not in list")
}

// returns all matching parsers
func MatchAll(file *os.File) (MatchingParsers, error) {

	fileSize, _ := fileSize(file)
	layout := parse.ParsedLayout{
		FileName: fileGetName(file),
		FileSize: fileSize}
	checker := parse.ParseChecker{
		File:         file,
		ParsedLayout: layout}
	m := map[string]Parser{}

	var err error
	checker.Header, err = readHeaderChunk(file)
	if err != nil {
		fmt.Println("warning: MatchAll failed reading header chunk")
		return nil, err
	}

	for name, parser := range parsers {

		pl, err2 := parser(&checker)
		if err2 != nil {
			fmt.Println("XXX parser", name, "failed")
			return nil, err2
		}
		if pl == nil {
			continue
		}
		if pl.FormatName == "" {
			pl.FormatName = name
		}
		m[name] = parser
	}
	return m, nil
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) { // XXX deprecate ParseLayout

	return matchParser(file)
}

// slice to expand, new length in bytes
func expandByteSlice(b []byte, newLen int64) []byte {

	i := int64(len(b))
	newLen = newLen - i
	return append(b[:i], append(make([]byte, newLen), b[i:]...)...)
}

func readHeaderChunk(file *os.File) ([]byte, error) {

	fileSize, _ := fileSize(file)
	maxHeaderLen := int64(0xffff)
	len := fileSize
	if len > maxHeaderLen {
		len = maxHeaderLen
	}

	b := make([]byte, len)

	file.Seek(0, os.SEEK_SET)
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		fmt.Println("error: failed to read header!", err)
		return nil, err
	}

	// resize to maxHeaderLen
	return expandByteSlice(b, maxHeaderLen), nil
}

func matchParser(file *os.File) (*parse.ParsedLayout, error) { // XXX deprecate! for MatchAll ...

	fileSize, err := fileSize(file)
	if err != nil {
		return nil, err
	}
	if fileSize == 0 {
		return nil, fmt.Errorf("empty file")
	}

	layout := parse.ParsedLayout{
		FileName: fileGetName(file),
		FileSize: fileSize}

	checker := parse.ParseChecker{
		File:         file,
		ParsedLayout: layout}

	checker.Header, err = readHeaderChunk(file)
	if err != nil {
		fmt.Println("warning: matchParser failed reading header chunk")
		return nil, err
	}

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
