package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// ...
const (
	Group DataType = 1 + iota // container

	// single bytes
	Int8
	Uint8

	// little endian
	Int16le
	Int32le
	Uint16le
	Uint32le
	Uint64le

	// big endian
	Uint16be
	Uint32be

	// version
	MajorMinor8    // high nibble = major, low = minor
	MajorMinor16le // first byte = major, last = minor
	MajorMinor16be // first byte = major, last = minor
	MajorMinor32le // first word = major, last = minor
	MinorMajor16le // first byte = minor, last = major

	// timestamps
	DOSDateTime

	// groups of bytes
	Bytes
	ASCII
	ASCIIC
	ASCIIZ
	RGB
)

const (
	Image FileKind = 1 + iota
	Archive
	AudioVideo
	Binary
	Executable
	Document
	Font
	WindowsResource
	MacOSResource
)

const (
	None TextEncoding = iota
	UTF8
	UTF16le
	UTF16be
	UTF32le
	UTF32be
)

var (
	dataTypes = map[DataType]string{
		Group:          "group",
		Int8:           "int8",
		Uint8:          "uint8",
		Int16le:        "int16-le",
		Uint16le:       "uint16-le",
		Int32le:        "int32-le",
		Uint32le:       "uint32-le",
		Uint64le:       "uint64-le",
		Uint16be:       "uint16-be",
		Uint32be:       "uint32-be",
		MajorMinor8:    "major.minor-8",
		MajorMinor16le: "major.minor-16le",
		MajorMinor16be: "major.minor-16be",
		MajorMinor32le: "major.minor-32le",
		MinorMajor16le: "minor.major-16le",
		DOSDateTime:    "datetime.dos-32le",
		Bytes:          "bytes",
		ASCII:          "ASCII",
		ASCIIC:         "ASCIIC",
		ASCIIZ:         "ASCIIZ",
		RGB:            "RGB",
	}
	textEncodings = map[TextEncoding]string{
		None:    "none",
		UTF8:    "utf8",
		UTF16le: "utf16le",
		UTF32le: "utf32le",
		UTF16be: "utf16be",
		UTF32be: "utf32be",
	}
	dataTypeBitsizes = map[DataType]int{
		Uint8:    8,
		Uint16le: 16,
		Uint32le: 32,
	}
	FileKinds = map[FileKind]string{
		Image:           "image",
		Archive:         "archive",
		AudioVideo:      "a/v",
		Binary:          "binary",
		Executable:      "executable",
		Document:        "document",
		Font:            "font",
		WindowsResource: "os-windows",
		MacOSResource:   "os-macos",
	}
)

// DataType ...
type DataType int

// FileKind ...
type FileKind int

// ParsedLayout ...
type ParsedLayout struct {
	FormatName   string
	FileName     string
	FileSize     int64
	FileKind     FileKind
	TextEncoding TextEncoding
	Layout       []Layout
}

type TextEncoding int

// Layout represents a parsed file structure
type Layout struct {
	Offset int64
	Length int64
	Type   DataType
	Info   string
	Childs []Layout
	Masks  []Mask
}

// Mask represents how to decode a bit field
type Mask struct {
	Low    int
	Length int
	Info   string
}

func (e TextEncoding) String() string {

	if val, ok := textEncodings[e]; ok {
		return val
	}

	// NOTE: this should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(int(e))
}

func (dt DataType) String() string {

	if val, ok := dataTypes[dt]; ok {
		return val
	}

	// NOTE: this should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(int(dt))
}

func (l *Layout) GetBitSize() int {

	if val, ok := dataTypeBitsizes[l.Type]; ok {
		return val
	}

	panic("GetBitSize: dont know size of " + l.Type.String())
}

func (l *Layout) parseByteN(reader io.Reader, expectedLen int64) ([]byte, error) {

	buf := make([]byte, expectedLen)

	readLen, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}

	if int64(readLen) != expectedLen {
		return nil, fmt.Errorf("Expected %d bytes, got %d", expectedLen, readLen)
	}
	return buf, nil
}

func (pl *ParsedLayout) isOffsetKnown(ofs int64) bool {

	for _, layout := range pl.Layout {
		if ofs >= layout.Offset && ofs < layout.Offset+int64(layout.Length) {
			return true
		}
	}
	return false
}

// returns a layout field with .Info quals `info`
func (pl *ParsedLayout) findInfoField(info string) *Layout {

	for _, layout := range pl.Layout {
		if layout.Info == info {
			return &layout
		}
		for _, childLayout := range layout.Childs {
			if childLayout.Info == info {
				return &childLayout
			}
		}
	}
	return nil
}

func (pl *ParsedLayout) findBitfieldLayout(info string) *Layout {

	for _, layout := range pl.Layout {
		for _, mask := range layout.Masks {
			if mask.Info == info {
				return &layout
			}
		}
		for _, childLayout := range layout.Childs {
			for _, mask := range childLayout.Masks {
				if mask.Info == info {
					return &childLayout
				}
			}
		}
	}
	return nil
}

func (pl *ParsedLayout) findBitfieldMask(info string) *Mask {

	for _, layout := range pl.Layout {
		for _, mask := range layout.Masks {
			if mask.Info == info {
				return &mask
			}
		}
		for _, childLayout := range layout.Childs {
			for _, mask := range childLayout.Masks {
				return &mask
			}
		}
	}
	return nil
}

// the output of cmd/prober --short
func (pl *ParsedLayout) ShortPrint() string {

	return pl.TypeSummary()
}

func (pl *ParsedLayout) TypeSummary() string {

	kindName := ""
	if val, ok := FileKinds[pl.FileKind]; ok {
		kindName = val
	}

	s := pl.FormatName + " " + kindName
	if pl.TextEncoding != None {
		s += " " + pl.TextEncoding.String()
	}
	return s
}

// the output of cmd/prober
func (pl *ParsedLayout) PrettyPrint() string {

	res := "Format: " + pl.FormatName + " (" + pl.FileName +
		", " + fmt.Sprintf("%d", pl.FileSize) + " bytes)\n\n"

	for _, layout := range pl.Layout {
		res += layout.Info + fmt.Sprintf(" (%04x)", layout.Offset) +
			", " + layout.Type.String() + "\n"

		for _, child := range layout.Childs {
			res += "  " + child.Info + fmt.Sprintf(" (%04x)", child.Offset) +
				", " + child.Type.String() + "\n"
		}
	}

	return res
}

// NOTE: went public for testing
func (pl *ParsedLayout) DecodeBitfieldFromInfo(file *os.File, info string) uint32 {

	field := pl.findBitfieldLayout(info)
	if field == nil {
		fmt.Println("ERROR: field", info, "not found")
		return 0
	}

	mask := pl.findBitfieldMask(info)
	if mask == nil {
		fmt.Println("ERROR: mask", info, "not found")
		return 0
	}

	b := ReadUnsignedInt(file, field)
	if bitmask, ok := bitmaskMap[mask.Length]; ok {

		tmp := bitmask << uint32(mask.Low)
		val := (b & tmp) >> uint32(mask.Low)
		return val
	}

	panic("need mask for len " + fmt.Sprintf("%d", mask.Length))
}

func (pl *ParsedLayout) ReadUint32leFromInfo(file *os.File, info string) (uint32, error) {

	layout := pl.findInfoField(info)
	if layout == nil {
		return 0, fmt.Errorf("ERROR didnt find field %v", info)
	}

	file.Seek(layout.Offset, os.SEEK_SET)

	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func (pl *ParsedLayout) readBytesFromInfo(file *os.File, info string) ([]byte, error) {

	layout := pl.findInfoField(info)
	if layout == nil {
		return nil, fmt.Errorf("ERROR didnt find field %v", info)
	}

	return ReadBytesFrom(file, layout.Offset, layout.Length), nil
}
