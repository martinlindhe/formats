package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// CurrentFieldInfo renders info of current field
func (state *HexViewState) CurrentFieldInfo(f *os.File, pl ParsedLayout) string {

	if len(pl.Layout) == 0 {
		fmt.Println("pl.Layout is empty")
		return ""
	}

	group := pl.Layout[state.CurrentGroup]

	res := "group: " + group.Info

	if state.BrowseMode == ByGroup {
		return res
	}

	if state.CurrentField >= len(group.Childs) {
		return "CHILD OUT OF RANGE"
	}

	field := group.Childs[state.CurrentField]

	res += "\n" + field.fieldInfoByType(f)
	res += " (" + field.Type.String() + ")"

	return res
}

func (field *Layout) fieldInfoByType(f *os.File) string {

	f.Seek(field.Offset, os.SEEK_SET)

	res := "field: " + field.Info + "\n"

	// decode data based on type and show

	switch field.Type {
	case Int8:
		var i int8
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint8:
		var i uint8
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case Int16le:
		var i int16
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint16le:
		var i uint16
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case Int32le:
		var i int32
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint32le:
		var i uint32
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d", i)

	case ASCII, ASCIIZ:
		buf := make([]byte, field.Length)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += string(buf)
	default:
		res += "XXX unhandled " + field.Type.String()
	}

	return res
}
