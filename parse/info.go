package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
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

	res += "\n └ " + field.fieldInfoByType(f) + " (" + field.Type.String()
	if field.Type == ASCII || field.Type == ASCIIZ {
		res += ", " + fmt.Sprintf("%d", field.Length) + " bytes"
	}
	res += ")"

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

		b := ReadUnsignedInt(f, field)

		for _, mask := range field.Masks {

			val := calcBitmask(&mask, b)
			res += fmt.Sprintf("%d: %s:%d = ", mask.Low, mask.Info, mask.Length) +
				fmt.Sprintf("%d", val) + "\n"
		}
		return res
	}

	switch field.Type { // XXX use func map
	case Int8:
		res += infoInt8(f, field)

	case Uint8:
		res += infoUint8(f, field)

	case Int16le:
		res += infoInt16le(f, field)

	case Uint16le:
		res += infoUint16le(f, field)

	case Int32le:
		res += infoInt32le(f, field)

	case Uint32le:
		res += infoUint32le(f, field)

	case Uint64le:
		res += infoUint64le(f, field)

	case Uint16be:
		res += infoUint16be(f, field)

	case Uint32be:
		res += infoUint32be(f, field)

	case Uint64be:
		res += infoUint64be(f, field)

	case MajorMinor16le:
		res += infoMajorMinor16le(f, field)

	case MajorMinor16be:
		res += infoMajorMinor16be(f, field)

	case MinorMajor16le:
		res += infoMinorMajor16le(f, field)

	case MajorMinor32le:
		res += infoMajorMinor32le(f, field)

	case DOSDateTime:
		res += infoDOSDateTime(f, field)

	case ArjDateTime:
		res += infoArjDateTime(f, field)

	case DOSOffsetSegment:
		res += infoDOSOffsetSegment(f, field)

	case Bytes:
		res += infoBytes(f, field)

	case ASCII, ASCIIZ:
		res += infoASCIIZ(f, field)

	case ASCIIC:
		res += infoASCIIC(f, field)

	case RGB:
		res += infoRGB(f, field)

	default:
		res += "unhandled type " + field.Type.String()
	}

	return res
}

// len (byte) + ASCII
func infoASCIIC(f *os.File, field *Layout) string {

	var len byte
	if err := binary.Read(f, binary.LittleEndian, &len); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}

	buf := make([]byte, len)
	_, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return string(buf)
}

func infoRGB(f *os.File, field *Layout) string {

	buf := make([]byte, field.Length)
	_, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%d, %d, %d", buf[0], buf[1], buf[2])
}

func infoASCIIZ(f *os.File, field *Layout) string {

	buf := make([]byte, field.Length)
	_, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return string(buf)
}

func infoDOSDateTime(f *os.File, field *Layout) string {

	var b uint32
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	t := time.Date(1970, time.January, 1, 1, 0, int(b), 0, time.UTC)
	return fmt.Sprintf("%v", t)
}

/*
 31 30 29 28 27 26 25 24 23 22 21 20 19 18 17 16
|<---- year-1980 --->|<- month ->|<--- day ---->|

 15 14 13 12 11 10  9  8  7  6  5  4  3  2  1  0
|<--- hour --->|<---- minute --->|<- second/2 ->|
*/
func infoArjDateTime(f *os.File, field *Layout) string {

	var b uint32
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}

	// XXX not correctly decoded

	t := time.Date(1970, time.January, 1, 1, 0, int(b), 0, time.UTC)
	return fmt.Sprintf("%v", t)
}

func infoDOSOffsetSegment(f *os.File, field *Layout) string {

	var b [2]uint16
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	abs := b[1]*16 + b[0]
	return fmt.Sprintf("%04x:%04x = %04x", b[1], b[0], abs)
}

func infoMajorMinor32le(f *os.File, field *Layout) string {

	var b [2]uint16
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%d.%d", b[0], b[1])
}

func infoMinorMajor16le(f *os.File, field *Layout) string {

	var b [2]uint8
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%d.%d", b[1], b[0])
}

func infoMajorMinor16be(f *os.File, field *Layout) string {

	var b [2]uint8
	if err := binary.Read(f, binary.BigEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%d.%d", b[0], b[1])
}

func infoMajorMinor16le(f *os.File, field *Layout) string {

	var b [2]uint8
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%d.%d", b[0], b[1])
}

func infoUint32be(f *os.File, field *Layout) string {

	var i uint32
	if err := binary.Read(f, binary.BigEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint64be(f *os.File, field *Layout) string {

	var i uint64
	if err := binary.Read(f, binary.BigEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint16be(f *os.File, field *Layout) string {

	var i uint16
	if err := binary.Read(f, binary.BigEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint64le(f *os.File, field *Layout) string {

	var i uint64
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint32le(f *os.File, field *Layout) string {

	var i uint32
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoInt32le(f *os.File, field *Layout) string {

	var i int32
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint16le(f *os.File, field *Layout) string {

	var i uint16
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoInt16le(f *os.File, field *Layout) string {

	var i int16
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoUint8(f *os.File, field *Layout) string {

	var i uint8
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoInt8(f *os.File, field *Layout) string {

	var i int8
	if err := binary.Read(f, binary.LittleEndian, &i); err != nil && err != io.EOF {
		return fmt.Sprintf("%v", err)
	}
	return field.prettyDecimalAndHex(int64(i))
}

func infoBytes(f *os.File, field *Layout) string {

	return "chunk of bytes"
}
