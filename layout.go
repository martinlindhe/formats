package formats

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
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

// Parser is a specialized parser for a file format
type Parser func(*parse.Checker) (*parse.ParsedLayout, error)

var (
	parsers = map[string]Parser{
		"7z":             archive.SevenZIP,
		"arj":            archive.ARJ,
		"bzip2":          archive.BZIP2,
		"cab":            archive.CAB,
		"compress-kwaj":  archive.CompressKWAJ,
		"deb":            archive.DEB,
		"ftcomp":         archive.FTCOMP,
		"gzip":           archive.GZIP,
		"iso":            archive.ISO,
		"luks":           archive.LUKS,
		"lzma":           archive.LZMA,
		"mozlz4":         archive.MozillaLZ4,
		"rar":            archive.RAR,
		"td2":            archive.TD2,
		"vdi":            archive.VDI,
		"wim":            archive.WIM,
		"winimg":         archive.WINIMG,
		"xar":            archive.XAR,
		"xz":             archive.XZ,
		"zip":            archive.ZIP,
		"zlib":           archive.ZLib,
		"aiff":           av.AIFF,
		"asf":            av.ASF,
		"au":             av.AU,
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
		"perl-enc":       bin.PerlENC,
		"rdb":            bin.RDB,
		"rifx":           bin.RIFX,
		"terminfo":       bin.Terminfo,
		"timezone":       bin.Timezone,
		"xkm":            bin.XKM,
		"chm":            doc.CHM,
		"hlp":            doc.HLP,
		"pdf":            doc.PDF,
		"rtf":            doc.RTF,
		"wri":            doc.WRI,
		"gb-rom":         exe.GameBoyROM,
		"gba-rom":        exe.GBAROM,
		"n64-rom":        exe.N64ROM,
		"dex":            exe.DEX,
		"elf":            exe.ELF,
		"java":           exe.JavaClass,
		"llvm-bitcode":   exe.LLVMBitcode,
		"lua":            exe.LUA,
		"macho":          exe.MachO,
		"mz":             exe.MZ,
		"python":         exe.PythonBytecode,
		"ruby":           exe.Ruby,
		"swf":            exe.SWF,
		"vbe":            exe.VBE,
		"dfont":          font.DFont,
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
		"bpg":            image.BPG,
		"gif":            image.GIF,
		"icns":           image.ICNS,
		"ico":            image.ICO,
		"jpeg":           image.JPEG,
		"pcx":            image.PCX,
		"png":            image.PNG,
		"psd":            image.PSD,
		"tga":            image.TGA, // XXX has too loose detection. thinks samples/images/ico/icon_003_win2k_cross.cur is tga!
		"tiff":           image.TIFF,
		"xcursor":        image.XCursor,
		"abba":           macos.ABBA,
		"bom_store":      macos.BOMStore,
		"bplist":         macos.BPLIST,
		"code_directory": macos.CodeDirectory,
		"code_signature": macos.CodeSignature,
		"ds_store":       macos.DSStore,
		"keychain":       macos.Keychain,
		"mtlb":           macos.MTLB,
		"rale":           macos.RALE,
		"rals":           macos.RALS,
		"styl":           macos.STYL,
		"vmap4":          macos.VMAP4,
		"ari8":           windows.ARI8,
		"cat":            windows.CAT,
		"dxbc":           windows.DXBC,
		"elst":           windows.ELST,
		"evtx":           windows.EVTX,
		"faxc":           windows.FAXC,
		"grp":            windows.GRP,
		"gres":           windows.GRES,
		"hwrs":           windows.HWRS,
		"imdx":           windows.IMDX,
		"immm":           windows.IMMM,
		"ldmf":           windows.LDMF,
		"lnk":            windows.LNK,
		"mam":            windows.MAM,
		"msft":           windows.MSFT,
		"ole-cf":         windows.OLECF,
		"p7x":            windows.P7X,
		"pcmh":           windows.PCMH,
		"pif":            windows.PIF,
		"pnf":            windows.PNF,
		"pri":            windows.PRI,
		"rbs":            windows.RBS,
		"regf":           windows.REGF,
		"rlfs":           windows.RLFS,
		"rsrc":           windows.RSRC,
		"sdi":            windows.SDI,
		"uce":            windows.UCE,
		"wave":           windows.WAVE,
		"xbf":            windows.XBF,
	}
)

// MatchingParsers is a list parsed layouts with different Parser
type MatchingParsers []parse.ParsedLayout

// First returns the first matching Parser
func (mp *MatchingParsers) First() *parse.ParsedLayout {

	for _, parser := range *mp {
		return &parser
	}
	return nil
}

// ChoseOne asks the user to select one of the matching parsers
func (mp *MatchingParsers) ChoseOne(file *os.File) (*parse.ParsedLayout, error) {

	i := 1
	fmt.Println("multiple parsers matched input file, please choose one:")
	for _, pl := range *mp {
		fmt.Printf("%d: %s\n", i, pl.FormatName)
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
	for _, pl := range *mp {
		if i == int(u) {
			return &pl, nil
		}
		i++
	}

	return nil, fmt.Errorf("selection not in list")
}

// MatchAll returns all matching parsers
func MatchAll(file *os.File) (MatchingParsers, error) {

	fileSize, _ := fileSize(file)
	layout := parse.ParsedLayout{
		FileName: fileGetName(file),
		FileSize: fileSize}
	checker := parse.Checker{
		File:         file,
		ParsedLayout: layout}
	var m MatchingParsers

	var err error
	checker.Header, err = readHeaderChunk(file)
	if err != nil {
		log.Println("warning: MatchAll failed reading header chunk")
		return nil, err
	}

	for name, parser := range parsers {

		pl, err2 := parser(&checker)
		if err2 != nil {
			log.Println("warning: parser", name, "failed: ", err2)
			continue
		}
		if pl == nil {
			continue
		}
		if pl.FormatName == "" {
			pl.FormatName = name
		}
		m = append(m, *pl)
	}

	if len(m) == 0 {
		// try text detector
		text, err := parse.Text(&checker)
		if err == nil {
			m = append(m, *text)
		} else {
			// fall back to raw
			raw, _ := parse.RAW(&checker)
			m = append(m, *raw)
		}
	}

	return m, nil
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*parse.ParsedLayout, error) {

	all, err := MatchAll(file)
	return all.First(), err
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
		log.Println("error: failed to read header!", err)
		return nil, err
	}

	// resize to maxHeaderLen
	return expandByteSlice(b, maxHeaderLen), nil
}

func fileGetName(file *os.File) string {
	stat, _ := file.Stat()
	return stat.Name()
}
