package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// CurrentFieldInfo renders info of current field
func (f *HexViewState) CurrentFieldInfo(file *os.File, pl ParsedLayout) string {

	if len(pl.Layout) == 0 {
		fmt.Println("pl.Layout is empty")
		return ""
	}

	field := pl.Layout[f.CurrentGroup]

	res := "field: " + field.Info

	res += "\nvalue: "

	file.Seek(field.Offset, os.SEEK_SET)

	// decode data based on type and show
	r := io.Reader(file)

	switch field.Type {
	case Int8:
		var i int8
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint8:
		var i uint8
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Int16le:
		var i int16
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint16le:
		var i uint16
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Int32le:
		var i int32
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint32le:
		var i uint32
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case ASCII, ASCIIZ:
		buf := make([]byte, field.Length)
		_, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		res += string(buf)

	default:
		res += "XXX unhandled " + field.Type.String()
	}

	res += " (" + field.Type.String() + ")"

	return res
}
