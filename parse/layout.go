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
	Uint16le
	Int32le
	Uint32le

	// big endian
	Uint16be
	Uint32be

	// version
	MajorMinor8    // high nibble = major, low = minor
	MajorMinor16le // first byte = major, last = minor
	MajorMinor16be // first byte = major, last = minor
	MajorMinor32le // first word = major, last = minor
	MinorMajor16le // first byte = minor, last = major

	// groups of bytes
	Bytes
	ASCII
	ASCIIC
	ASCIIZ
	RGB
)

var (
	dataTypes = map[DataType]string{
		Group:          "Group",
		Int8:           "int8",
		Uint8:          "uint8",
		Int16le:        "int16-le",
		Uint16le:       "uint16-le",
		Int32le:        "int32-le",
		Uint32le:       "uint32-le",
		Uint16be:       "uint16-be",
		Uint32be:       "uint32-be",
		MajorMinor8:    "major.minor-8",
		MajorMinor16le: "major.minor-16le",
		MajorMinor16be: "major.minor-16be",
		MajorMinor32le: "major.minor-32le",
		MinorMajor16le: "minor.major-16le",
		Bytes:          "bytes",
		ASCII:          "ASCII",
		ASCIIC:         "ASCIIC",
		ASCIIZ:         "ASCIIZ",
		RGB:            "RGB",
	}
)

// ParsedLayout ...
type ParsedLayout struct {
	FormatName string
	FileSize   int64
	Layout     []Layout
}

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

// DataType ...
type DataType int

func (dt DataType) String() string {

	if val, ok := dataTypes[dt]; ok {
		return val
	}

	// NOTE should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(int(dt))
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

// the output of cmd/prober
func (pl *ParsedLayout) PrettyPrint() string {

	res := ""
	for _, layout := range pl.Layout {
		res += layout.Info + fmt.Sprintf(" (%04x)", layout.Offset) + ", " + layout.Type.String() + "\n"

		for _, child := range layout.Childs {
			res += "  " + child.Info + fmt.Sprintf(" (%04x)", child.Offset) + ", " + child.Type.String() + "\n"
		}
	}

	return res
}

func (pl *ParsedLayout) updateLabel(label string, newLabel string) {

	for layoutIdx, layout := range pl.Layout {
		if layout.Info == label {
			pl.Layout[layoutIdx].Info = newLabel
			return
		}
		for childIdx, child := range layout.Childs {
			if child.Info == label {
				pl.Layout[layoutIdx].Childs[childIdx].Info = newLabel
				return
			}
		}
	}

	panic("label not found: " + label)
}

// NOTE: went public for testing
func (pl *ParsedLayout) DecodeBitfieldFromInfo(file *os.File, info string) uint32 {

	layout := pl.findBitfieldLayout(info)
	if layout == nil {
		fmt.Println("ERROR: layout", info, "not found")
		return 0
	}

	mask := pl.findBitfieldMask(info)
	if mask == nil {
		fmt.Println("ERROR: mask", info, "not found")
		return 0
	}

	file.Seek(layout.Offset, os.SEEK_SET)

	var b byte
	binary.Read(file, binary.LittleEndian, &b)

	// XXX mask bits accordingly ....

	m := uint32(0)
	if mask.Low == 0 {
		switch mask.Length {
		case 3:
			m = 7
		default:
			panic("FIXME unhandled mask len!")
		}

		return uint32(b) & m
	}

	panic("XXX fixme handle bit shifts stuff and tests")
}

func (pl *ParsedLayout) readUint32leFromInfo(file *os.File, info string) (uint32, error) {

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

	return readBytesFrom(file, layout.Offset, layout.Length), nil
}

type ByLayout []Layout

func (slice ByLayout) Len() int {
	return len(slice)
}

func (slice ByLayout) Less(i, j int) bool {
	return slice[i].Offset < slice[j].Offset
}

func (slice ByLayout) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
