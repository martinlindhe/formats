package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ReadBitmask(file *os.File, layout *Layout, mask *Mask) uint32 {

	b := ReadUnsignedInt(file, layout)
	return CalcBitmask(mask, b)
}

func CalcBitmask(mask *Mask, b uint32) uint32 {

	if bitmask, ok := bitmaskMap[mask.Length]; ok {

		tmp := bitmask << uint32(mask.Low)
		val := (b & tmp) >> uint32(mask.Low)

		return val
	}

	panic("add mask for length " + fmt.Sprintf("%d", mask.Length))
}

func ReadToMap(file *os.File, dataType DataType, pos int64, i interface{}) (string, error) {

	switch dataType {
	case Uint8:
		idx, err := ReadUint8(file, pos)
		if err != nil {
			return "", err
		}
		a := i.(map[byte]string)
		if val, ok := a[idx]; ok {
			return val, nil
		}
	case Uint16le:
		idx, err := ReadUint16le(file, pos)
		if err != nil {
			return "", err
		}
		a := i.(map[uint16]string)
		if val, ok := a[idx]; ok {
			return val, nil
		}
	case Uint32le:
		idx, err := ReadUint32le(file, pos)
		if err != nil {
			return "", err
		}
		a := i.(map[uint32]string)
		if val, ok := a[idx]; ok {
			return val, nil
		}
	case Uint16be:
		idx, err := ReadUint16be(file, pos)
		if err != nil {
			return "", err
		}
		a := i.(map[uint16]string)
		if val, ok := a[idx]; ok {
			return val, nil
		}
	case Uint32be:
		idx, err := ReadUint32be(file, pos)
		if err != nil {
			return "", err
		}
		a := i.(map[uint32]string)
		if val, ok := a[idx]; ok {
			return val, nil
		}
	default:
		fmt.Println("UNHNANdlNED:", dataType.String())
	}
	return "?", nil
}

// ReadZeroTerminatedASCIIUntil returns string, bytes read, error
func ReadZeroTerminatedASCIIUntil(file *os.File, pos int64, maxLen int) (string, int, error) {

	file.Seek(pos, os.SEEK_SET)

	var c byte
	s := ""

	readCnt := 0
	for {
		if err := binary.Read(file, binary.LittleEndian, &c); err != nil {
			return s, 0, err
		}
		readCnt++
		if c == 0 {
			break
		}
		s += string(c)
		if readCnt == maxLen {
			break
		}
	}
	return s, readCnt, nil
}

func ReadBytesFrom(file *os.File, pos int64, size int64) []byte {

	file.Seek(pos, os.SEEK_SET)

	b := make([]byte, size)
	binary.Read(file, binary.LittleEndian, &b)
	return b
}

// ReadUnsignedInt reads field value, as an uint32
func ReadUnsignedInt(file *os.File, field *Layout) uint32 {

	switch field.Type {
	case Uint8:
		val, _ := ReadUint8(file, field.Offset)
		return uint32(val)

	case Uint16le:
		val, _ := ReadUint16le(file, field.Offset)
		return uint32(val)

	case Uint32le:
		val, _ := ReadUint32le(file, field.Offset)
		return val
	}

	panic("ReadUnsignedInt: unhandled type " + field.Type.String())
}

func ReadUint8(file *os.File, pos int64) (uint8, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint8
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func ReadUint16be(file *os.File, pos int64) (uint16, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func ReadUint16le(file *os.File, pos int64) (uint16, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func ReadUint32be(file *os.File, pos int64) (uint32, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func ReadUint64be(file *os.File, pos int64) (uint64, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint64
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func ReadUint32le(file *os.File, pos int64) (uint32, error) {

	file.Seek(pos, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

type LineEnding int

const (
	Crlf = 1 + iota // windows
	Cr              // old mac os
	Lf              // linux + modern mac os
)

func ReadBytesUntilNewline(file *os.File, pos int64) ([]byte, int64, error) {

	var c byte
	var b []byte
	readCnt := int64(0)

	lineEnding := Lf
	for {
		file.Seek(pos, os.SEEK_SET)

		if err := binary.Read(file, binary.LittleEndian, &c); err != nil {
			if err == io.EOF {
				break
			}
			return b, 0, err
		}
		pos++
		readCnt++
		b = appendByte(b, c)

		if c == '\r' {
			if lineEnding != Crlf {
				if err := binary.Read(file, binary.LittleEndian, &c); err != nil {
					if err == io.EOF {
						break
					}
					return b, 0, err
				}
				if c == '\n' {
					lineEnding = Crlf
				}
			}

			if lineEnding == Crlf {
				pos++
				readCnt++
			} else {
				lineEnding = Cr
			}
			break
		}

		if c == '\n' || c == '\r' {
			break
		}
	}

	return b, readCnt, nil
}

// from http://blog.golang.org/go-slices-usage-and-internals
func appendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}
