package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// ParsedLayout ...
type ParsedLayout struct {
	FormatName string
	FileSize   int64
	Layout     []Layout
}

// Layout represents a parsed file structure
type Layout struct { // XXX aka Chunk in cs code
	Offset int64
	Length int64
	Type   DataType
	Info   string
	Childs []Layout // XXX make use of + display. parent is a layout group
}

// DataType ...
type DataType int

func (dt DataType) String() string {

	m := map[DataType]string{
		Int8:     "int8",
		Uint8:    "uint8",
		Int16le:  "int16-le",
		Uint16le: "uint16-le",
		Int32le:  "int32-le",
		Uint32le: "uint32-le",
		ASCII:    "ASCII",
		ASCIIZ:   "ASCIIZ",
		Group:    "Group",
	}

	if val, ok := m[dt]; ok {
		return val
	}

	// NOTE should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(int(dt))
}

// ...
const (
	Group DataType = 1 + iota
	Int8
	Uint8
	Int16le
	Uint16le
	Int32le
	Uint32le
	ASCII
	ASCIIZ
)

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

// the output of cmd/prober
func (parsedLayout *ParsedLayout) PrettyPrint() string {

	res := ""
	for _, layout := range parsedLayout.Layout {
		res += layout.Info + fmt.Sprintf(" (%04x)", layout.Offset) + ", " + layout.Type.String() + "\n"

		for _, child := range layout.Childs {
			res += "  " + child.Info + fmt.Sprintf(" (%04x)", child.Offset) + ", " + child.Type.String() + "\n"
		}
	}

	return res
}

func (pl *ParsedLayout) readUint32leFromInfo(file *os.File, info string) uint32 {

	layout := pl.findInfoField(info)
	if layout == nil {
		fmt.Println("ERROR didnt find field", info)
		return 0
	}

	file.Seek(layout.Offset, os.SEEK_SET)

	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b
}
