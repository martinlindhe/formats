package formats

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"
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

// ParseLayout returns a Layout for the file
func ParseLayout(file *os.File) ([]Layout, error) {

	// XXX we fake result from structToFlatStruct() to test presentation
	return []Layout{
		Layout{0x0000, 2, Uint16le, "magic"},
		Layout{0x0002, 4, Uint32le, "width"},
		Layout{0x0006, 4, Uint32le, "height"},
		Layout{0x000a, 9, ASCIIZ, "NAME.EXT"},
		Layout{0x000a + 9, 2, Uint16le, "tag"},
	}, nil
}

func structToFlatStruct(obj interface{}) []Layout { // XXX implement

	res := []Layout{}

	//	spew.Dump(x)

	// XXX iterate over struct, create a 2d rep of the structure mapping

	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		// XXX is it a struct ?
		// fmt.Println(f.Tag)

		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	return res
}

func PrettyHexView(file *os.File, fileLayout []Layout) string {

	hex := ""

	base := HexView.StartingRow * int64(HexView.RowWidth)
	ceil := base + int64(HexView.VisibleRows*HexView.RowWidth)

	layout := fileLayout[HexView.CurrentField]
	// fmt.Printf("Using field %v, field %d\n", val, currentField)

	lineCount := 0
	for i := base; i < ceil; i += int64(HexView.RowWidth) {

		file.Seek(i, os.SEEK_SET)
		line, err := GetHex(file, layout)

		hex += fmt.Sprintf("[[%04x]](fg-yellow) %s\n", i, line)
		lineCount++
		if err != nil {
			fmt.Println("got err", err)
			break
		}
	}

	fmt.Println("XXX lineCount:", lineCount)
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
