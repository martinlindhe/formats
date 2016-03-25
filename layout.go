package formats

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Layout represents a parsed file structure layout as a flat list
type Layout struct {
	Offset int64
	Length byte
	Type   DataType
	Info   string
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
	panic(dt)
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

func fileExt(file *os.File) string {

	ext := filepath.Ext(file.Name())
	if len(ext) > 0 {
		// strip leading dot
		ext = ext[1:]
	}
	return ext
}

// ParseLayout returns a ParsedLayout for the file
func ParseLayout(file *os.File) (*ParsedLayout, error) {

	parsed, err := parseFileByDescription(file, fileExt(file))
	if parsed == nil {
		fmt.Println(err)
		panic("XXX if find by extension fails, search all for magic id")
	}

	return parsed, err
}

func getFileSize(file *os.File) int64 {

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}

func parseFileByDescription(
	file *os.File, formatName string) (*ParsedLayout, error) {

	format, err := ReadFormatDescription(formatName)
	if err != nil {
		return nil, err
	}

	res := ParsedLayout{
		FormatName: formatName,
		FileSize:   getFileSize(file),
	}

	for _, step := range format.Details {

		layout, err := res.intoLayout(file, step)
		if err != nil {
			fmt.Println("trouble parsing:", err)
		}

		res.Layout = append(res.Layout, *layout)
	}

	return &res, nil
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

// transforms a part of file into a Layout, according to `step`
func (pl *ParsedLayout) intoLayout(file *os.File, step string) (*Layout, error) {

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

	if expectedLen, err := parseExpectedBytes(&layout, reader, param1, param2); err == nil {
		layout.Length = byte(expectedLen)
		layout.Type = ASCII
	} else if _, err := parseExpectedByte(reader, param1, param2); err == nil {
		layout.Length = 1
		layout.Type = Uint8
	} else if _, err := parseExpectedUint16le(reader, param1, param2); err == nil {
		layout.Length = 2
		layout.Type = Uint16le
	} else {
		return nil, fmt.Errorf("dunno how to handle %s", param1)
	}

	return &layout, nil
}

func parseExpectedUint16le(reader io.Reader, param1 string, param2 string) (uint16, error) {

	if param1 != "uint16le" {
		return 0, fmt.Errorf("wrong type")
	}
	var b uint16
	err := binary.Read(reader, binary.LittleEndian, &b);
	return b, err
}

func parseExpectedByte(reader io.Reader, param1 string, param2 string) (byte, error) {

	if param1 != "uint8" && param1 != "byte" {
		return 0, fmt.Errorf("wrong type")
	}
	// XXX "byte", params[2] describes a bit field
	var b byte
	err := binary.Read(reader, binary.LittleEndian, &b);
	return b, err
}

func parseExpectedBytes(layout *Layout, reader io.Reader, param1 string, param2 string) (int64, error) {

	p1 := strings.Split(param1, ":")

	if p1[0] != "byte" || len(p1) != 2 {
		return 0, fmt.Errorf("wrong type")
	}

	expectedLen, err := parseExpectedLen(p1[1])
	if err != nil {
		return 0, err
	}

	// "byte:3", params[2] holds the bytes
	buf, err := layout.parseByteN(reader, expectedLen)
	if err != nil {
		return 0, err
	}

	// split expected forms on comma
	expectedForms := strings.Split(param2, ",")
	for _, expectedForm := range expectedForms {

		expectedBytes := []byte(expectedForm)
		if int64(len(expectedForm)) == 2*expectedLen {
			// hex string?
			bytes, err := hex.DecodeString(expectedForm)
			if err == nil && byteSliceEquals(buf, bytes) {
				return expectedLen, nil
			}
		}
		if string(buf) == string(expectedBytes) {
			return expectedLen, nil
		}
	}

	return 0, fmt.Errorf("didnt find expected bytes %s", param2)
}

func parseExpectedLen(s string) (int64, error) {
	expectedLen, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if expectedLen > 255 {
		return 0, fmt.Errorf("len too big (max 255)")
	}
	if expectedLen <= 0 {
		return 0, fmt.Errorf("len too small (min 1)")
	}
	return expectedLen, nil
}

// PrettyHexView ...
func (pl *ParsedLayout) PrettyHexView(file *os.File) string {

	ofsFmt := "%08x"
	if pl.FileSize <= 0xffff {
		ofsFmt = "%04x"
	} else if pl.FileSize <= 0xffffff {
		ofsFmt = "%06x"
	}

	hex := ""

	base := HexView.StartingRow * int64(HexView.RowWidth)
	ceil := base + int64(HexView.VisibleRows*HexView.RowWidth)

	for i := base; i < ceil; i += int64(HexView.RowWidth) {

		ofs, err := file.Seek(i, os.SEEK_SET)
		if i != ofs {
			log.Fatalf("err: unexpected offset %04x, expected %04x\n", ofs, i)
		}
		line, err := pl.GetHex(file)

		ofsText := fmt.Sprintf(ofsFmt, i)

		hex += fmt.Sprintf("[[%s]](fg-yellow) %s\n", ofsText, line)
		if err != nil {
			fmt.Println("got err", err)
			break
		}
	}

	return hex
}

func (pl *ParsedLayout) isOffsetKnown(ofs int64) bool {

	for _, layout := range pl.Layout {
		if ofs >= layout.Offset && ofs < layout.Offset+int64(layout.Length) {
			return true
		}
	}
	return false
}

// GetHex dumps a row of hex from io.Reader
func (pl *ParsedLayout) GetHex(file *os.File) (string, error) {

	layout := pl.Layout[HexView.CurrentField]

	reader := io.Reader(file)

	symbols := []string{}

	base, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return "", err
	}

	for w := int64(0); w < 16; w++ {
		var b byte
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			if err == io.EOF {
				return combineHexRow(symbols), nil
			}
			return "", err
		}

		ceil := base + w

		colorName := "fg-white"
		if !pl.isOffsetKnown(base + w) {
			colorName = "fg-red"
		}
		if ceil >= layout.Offset && ceil < layout.Offset+int64(layout.Length) {
			colorName = "fg-blue"
		}

		group := fmt.Sprintf("[%02x](%s)", b, colorName)
		symbols = append(symbols, group)
	}

	return combineHexRow(symbols), nil
}
