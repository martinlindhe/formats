package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

// LineEnding ...
type LineEnding int

// line endings
const (
	Crlf = 1 + iota // windows
	Cr              // old mac os
	Lf              // linux + modern mac os
)

// ReadToMap maps a bit value to a []string map
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
		log.Println("Unhandled type:", dataType.String())
	}
	return "?", nil
}

// ReadZeroTerminatedASCIIUntil returns ascii-string, bytes read, error
func ReadZeroTerminatedASCIIUntil(file *os.File, pos int64, maxLen int64) (string, int64, error) {
	c := byte(0)
	s := ""
	readCnt := int64(0)
	file.Seek(pos, os.SEEK_SET)

	for {
		if err := binary.Read(file, binary.LittleEndian, &c); err != nil {
			return s, 0, err
		}
		readCnt++
		if c == 0 {
			break
		}
		s += string(c)
		if readCnt >= maxLen {
			break
		}
	}
	return s, readCnt, nil
}

// ReadBytesFrom reads `size` bytes from `file`
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
	log.Fatal("ReadUnsignedInt: unhandled type " + field.Type.String())
	return 0
}

// ReadUint8 reads Uint8 from `file`
func ReadUint8(file *os.File, pos int64) (uint8, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint8
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

// ReadUint16be reads big endian Uint16 from `file`
func ReadUint16be(file *os.File, pos int64) (uint16, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

// ReadUint16le reads little endian Uint16 from `file`
func ReadUint16le(file *os.File, pos int64) (uint16, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

// ReadUint32be reads big endian Uint32 from `file`
func ReadUint32be(file *os.File, pos int64) (uint32, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

// ReadUint64be reads big endian Uint64 from `file`
func ReadUint64be(file *os.File, pos int64) (uint64, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint64
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

// ReadUint32le reads little endian Uint32 from `file`
func ReadUint32le(file *os.File, pos int64) (uint32, error) {
	file.Seek(pos, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

// ReadBytesUntilNewline is used to process text
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

func calcBitmask(mask *Mask, b uint32) uint32 {
	if bitmask, ok := bitmaskMap[mask.Length]; ok {
		tmp := bitmask << uint32(mask.Low)
		val := (b & tmp) >> uint32(mask.Low)
		return val
	}
	log.Fatal("mask missing for length " + fmt.Sprintf("%d", mask.Length))
	return 0
}
