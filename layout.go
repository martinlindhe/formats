package formats

import (
	"encoding/binary"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FormatDescription ...
type FormatDescription struct {
	Format Format `json:"format"`
}

// Format ...
type Format struct {
	Name   string   `json:"name"`
	Mime   string   `json:"mime"`
	Struct []string `json:"struct"`
}

// ReadFormatDescription ...
func ReadFormatDescription(formatName string) (*Format, error) {

	formatFile := "./formats/" + formatName + ".yml"

	if !exists(formatFile) {
		return nil, fmt.Errorf("Unknown format %s", formatFile)
	}

	data, err := ioutil.ReadFile(formatFile)
	if err != nil {
		return nil, err
	}

	desc := FormatDescription{}
	err = yaml.Unmarshal(data, &desc)
	return &desc.Format, err
}

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

// ParseLayout returns a Layout for the file
func ParseLayout(file *os.File) (*ParsedLayout, error) {

	parsed, err := parseFileByDescription(file, fileExt(file))
	if parsed == nil {
		panic("XXX if find by extension fails, search all for magic id")
	}

	return parsed, err
}

func parseFileByDescription(file *os.File, formatName string) (*ParsedLayout, error) {

	format, err := ReadFormatDescription(formatName)
	if err != nil {
		return nil, err
	}

	reader := io.Reader(file)

	res := ParsedLayout{
		FormatName: formatName,
	}

	for _, step := range format.Struct {

		// params: name | data type and size | type-dependant
		params := strings.Split(step, "|")

		layout := Layout{}

		layout.Offset, _ = file.Seek(0, os.SEEK_CUR)
		layout.Info = params[0]

		p1 := strings.Split(params[1], ":")

		if p1[0] == "byte" && len(p1) == 2 {
			// "byte:3", params[2] holds the bytes

			expectedLen, err := strconv.ParseInt(p1[1], 10, 64)
			if err != nil {
				return nil, err
			}
			if expectedLen > 255 {
				return nil, fmt.Errorf("byte:len too big (max 255)")
			}
			if expectedLen <= 0 {
				return nil, fmt.Errorf("byte:len len must be at least 1")
			}

			layout.Length = byte(expectedLen)
			layout.Type = ASCII

			buf := make([]byte, expectedLen)

			_, err = reader.Read(buf)
			if err != nil {
				return nil, err
			}

			// split expected forms on comma
			expectedForms := strings.Split(params[2], ",")
			found := false
			for _, expectedForm := range expectedForms {

				expectedBytes := []byte(expectedForm)
				// fmt.Println("expects to find", expectedLen, "bytes:", string(expectedBytes))
				if !found && string(buf) == string(expectedBytes) {
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("didnt find expected bytes %s", params[2])
			}

		} else if params[1] == "uint8" || params[1] == "byte" {
			// "byte", params[2] describes a bit field

			layout.Length = 1
			layout.Type = Uint8

			var b byte
			if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
				fmt.Println(b) // XXX make use of+!
			}

		} else if params[1] == "uint16le" {
			layout.Length = 2
			layout.Type = Uint16le

			var b uint16
			if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
				fmt.Println(b) // XXX make use of+!
			}

		} else {
			return nil, fmt.Errorf("dunno how to handle %s", params[1])
		}

		res.Layout = append(res.Layout, layout)
	}

	return &res, nil
}

// PrettyHexView ...
func (pl *ParsedLayout) PrettyHexView(file *os.File) string {

	hex := ""

	base := HexView.StartingRow * int64(HexView.RowWidth)
	ceil := base + int64(HexView.VisibleRows*HexView.RowWidth)

	layout := pl.Layout[HexView.CurrentField]
	// fmt.Printf("Using field %v, field %d\n", val, currentField)

	for i := base; i < ceil; i += int64(HexView.RowWidth) {

		ofs, err := file.Seek(i, os.SEEK_SET)
		if i != ofs {
			fmt.Printf("err: unexpected offset %04x, expected %04x\n", ofs, i)
		}
		line, err := GetHex(file, layout)

		hex += fmt.Sprintf("[[%04x]](fg-yellow) %s\n", i, line)
		if err != nil {
			fmt.Println("got err", err)
			break
		}
	}

	return hex
}

// GetHex dumps a row of hex from io.Reader
func GetHex(file *os.File, layout Layout) (string, error) {

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

		groupFmt := "%02x"
		ceil := base + w

		if ceil >= layout.Offset && ceil < layout.Offset+int64(layout.Length) {
			groupFmt = "[%02x](fg-blue)"
		}

		group := fmt.Sprintf(groupFmt, b)
		symbols = append(symbols, group)
	}

	return combineHexRow(symbols), nil
}
