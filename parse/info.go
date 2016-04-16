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

	res := group.Info

	if state.BrowseMode == ByGroup {
		return res
	}

	if state.CurrentField >= len(group.Childs) {
		return "CHILD OUT OF RANGE"
	}

	field := group.Childs[state.CurrentField]

	res += "\n └ " + field.fieldInfoByType(f) + " (" + field.Type.String() + ")"

	return res
}

func (field *Layout) prettyDecimalAndHex(i int64) string {

	dec := fmt.Sprintf("%d", i)
	hex := fmt.Sprintf("%x", i)
	if dec == hex {
		return dec
	}
	return dec + " (" + hex + ")"
}

func (field *Layout) fieldInfoByType(f *os.File) string {

	f.Seek(field.Offset, os.SEEK_SET)

	res := field.Info + "\n\n"

	// decode data based on type and show

	switch field.Type {
	case Int8:
		var i int8
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint8:
		var i uint8
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Int16le:
		var i int16
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint16le:
		var i uint16
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Int32le:
		var i int32
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint32le:
		var i uint32
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint64le:
		var i uint64
		if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint16be:
		var i uint16
		if err := binary.Read(f, binary.BigEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case Uint32be:
		var i uint32
		if err := binary.Read(f, binary.BigEndian, &i); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += field.prettyDecimalAndHex(int64(i))

	case MajorMinor16le:
		var b [2]uint8
		if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d.%d", b[0], b[1])

	case MinorMajor16le:
		var b [2]uint8
		if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d.%d", b[1], b[0])

	case MajorMinor32le:
		var b [2]uint16
		if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d.%d", b[0], b[1])

	case Bytes:
		res += fmt.Sprintf("chunk of bytes")

	case ASCII, ASCIIZ:
		buf := make([]byte, field.Length)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += string(buf)

	case ASCIIC:
		// len (byte) + ASCII
		var len byte
		if err := binary.Read(f, binary.LittleEndian, &len); err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}

		buf := make([]byte, len)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += string(buf)

	case RGB:
		buf := make([]byte, field.Length)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Sprintf("%v", err)
		}
		res += fmt.Sprintf("%d, %d, %d", buf[0], buf[1], buf[2])
		//val := uint64(buf[0])<<16 | uint64(buf[1])<<8 | uint64(buf[2])
		//res += fmt.Sprintf("%06x", val)

	default:
		res += "XXX unhandled " + field.Type.String()
	}

	return res
}
