package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var (
	bitmaskMap = map[int]uint32{
		1:  1,
		2:  3,
		3:  7,
		4:  0xf,
		5:  0x1f,
		6:  0x3f,
		7:  0x7f,
		8:  0xff,
		9:  0x1ff,
		10: 0x3ff,
		11: 0x7ff,
		12: 0xfff,
	}
)

// CurrentFieldInfo renders info of current field
func (state *HexViewState) CurrentFieldInfo(f *os.File, pl ParsedLayout) string {

	if len(pl.Layout) == 0 {
		fmt.Println("CurrentFieldInfo: pl.Layout is empty")
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

	// decode bit mask
	if len(field.Masks) > 0 {
		var b uint32

		switch field.Type {
		case Uint8:
			val, _ := ReadUint8(f, field.Offset)
			b = uint32(val)

		default:
			panic("unknown bitmask size " + field.Type.String())
		}

		for _, mask := range field.Masks {

			if bitmask, ok := bitmaskMap[mask.Length]; ok {

				tmp := bitmask << uint32(mask.Low)
				val := (b & tmp) >> uint32(mask.Low)

				res += fmt.Sprintf("%d: %s:%d = ", mask.Low, mask.Info, mask.Length) +
					fmt.Sprintf("%d", val) + "\n"

			} else {
				panic("need mask for len " + fmt.Sprintf("%d", mask.Length))
			}
		}
		return res
	}

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
