package formats

import (
	"encoding/binary"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
		ASCIIZ:   "ASCIIZ",
		Byte:     "byte",
		Uint16le: "uint16-le",
		Uint32le: "uint32-le",
		Int16le:  "int16-le",
		Int32le:  "int32-le",
	}

	if val, ok := m[dt]; ok {
		return val
	}

	// NOTE should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(dt)
}

// ...
const (
	_               = iota
	ASCIIZ DataType = iota
	Byte
	Uint16le
	Uint32le
	Int16le
	Int32le
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
func ParseLayout(file *os.File) ([]Layout, error) {

	parsed, err := parseFileByDescription(file, fileExt(file))

	if parsed == nil {
		fmt.Println("XXX if find by extension fails, search all for magic id")
	}

	return parsed, err
}

func parseFileByDescription(file *os.File, formatName string) ([]Layout, error) {

	format, err := ReadFormatDescription(formatName)
	if err != nil {
		return nil, err
	}

	fmt.Println(format)

	// XXX parse yml
	fmt.Println("XXX parse", formatName)
	return nil, nil
}

// PrettyHexView ...
func PrettyHexView(file *os.File, fileLayout []Layout) string {

	hex := ""

	base := HexView.StartingRow * int64(HexView.RowWidth)
	ceil := base + int64(HexView.VisibleRows*HexView.RowWidth)

	layout := fileLayout[HexView.CurrentField]
	// fmt.Printf("Using field %v, field %d\n", val, currentField)

	for i := base; i < ceil; i += int64(HexView.RowWidth) {

		file.Seek(i, os.SEEK_SET)
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
func GetHex(file *os.File, layout Layout) (res string, err error) {

	reader := io.Reader(file)

	symbols := []string{}

	base, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return "", err
	}

	for w := int64(0); w < 16; w++ {
		var b byte
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			res = combineHexRow(symbols)
			return
		}

		groupFmt := "%02x"
		ceil := base + w

		if ceil >= layout.Offset && ceil < layout.Offset+int64(layout.Length) {
			groupFmt = "[%02x](fg-blue)"
		}

		group := fmt.Sprintf(groupFmt, b)
		symbols = append(symbols, group)
	}
	res = combineHexRow(symbols)
	return
}
