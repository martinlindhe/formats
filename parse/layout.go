package parse

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Layout represents a parsed file structure
type Layout struct {
	Offset int64
	Length byte
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
	}

	if val, ok := m[dt]; ok {
		return val
	}

	// NOTE should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(int(dt))
}

// ParsedLayout ...
type ParsedLayout struct {
	FormatName string
	FileSize   int64
	Layout     []Layout
}

// ...
const (
	Int8 DataType = 1 + iota
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

// transforms a part of file into a Layout, according to `step`
func (pl *ParsedLayout) intoLayout(file *os.File, step string) (*Layout, error) {
	// XXX unused
	reader := io.Reader(file)

	// params: name | data type and size | type-dependant
	params := strings.Split(step, "|")

	layout := Layout{}

	layout.Offset, _ = file.Seek(0, os.SEEK_CUR)
	layout.Info = params[0]

	param1 := ""
	param2 := ""
	if len(params) > 1 {
		param1 = params[1]
	}
	if len(params) > 2 {
		param2 = params[2]
	}

	if b, err := parseExpectedBytes(&layout, reader, param1, param2); err == nil {
		layout.Length = byte(len(b))
		layout.Type = ASCII
	} else if _, err := parseExpectedByte(reader, param1, param2); err == nil {
		layout.Length = 1
		layout.Type = Uint8
	} else if _, err := parseExpectedUint16le(reader, param1, param2); err == nil {
		layout.Length = 2
		layout.Type = Uint16le
	} else if _, err := parseExpectedUint32le(reader, param1, param2); err == nil {
		layout.Length = 4
		layout.Type = Uint32le
	} else {
		return nil, fmt.Errorf("dunno how to handle %s, %s, %s", params[0], param1, param2)
	}

	return &layout, nil
}

func (pl *ParsedLayout) isOffsetKnown(ofs int64) bool {

	for _, layout := range pl.Layout {
		if ofs >= layout.Offset && ofs < layout.Offset+int64(layout.Length) {
			return true
		}
	}
	return false
}
