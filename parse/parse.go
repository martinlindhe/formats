package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

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

func knownLengthASCII(file *os.File, offset int64, length int) (string, error) {

	file.Seek(offset, os.SEEK_SET)

	var c byte
	s := ""

	len := 0
	for {
		if err := binary.Read(file, binary.LittleEndian, &c); err != nil {
			return s, err
		}
		if len == length {
			break
		}
		s += string(c)
		len++
	}
	return s, nil
}

// return string, bytes read, error
func countInitiatedASCII(r io.Reader) (string, int, error) {

	var count byte
	var c byte
	s := ""

	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return s, 0, err
	}
	readCnt := 0
	for i := 0; i < int(count); i++ {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return s, 0, err
		}
		s += string(c)
		readCnt++
	}
	return s, readCnt, nil
}

func ReadZeroTerminatedASCIIUntil(file *os.File, offset int64, maxLen int) (string, int, error) {

	file.Seek(offset, os.SEEK_SET)
	return zeroTerminatedASCIIUntil(file, maxLen)
}

// return string, bytes read, error
func zeroTerminatedASCIIUntil(r io.Reader, maxLen int) (string, int, error) {

	var c byte
	s := ""

	readCnt := 0
	for {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
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

// return string, bytes read, error
func zeroTerminatedASCII(r io.Reader) (string, int, error) {

	var c byte
	s := ""

	readCnt := 0
	for {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return s, 0, err
		}
		readCnt++
		if c == 0 {
			break
		}
		s += string(c)
	}
	return s, readCnt, nil
}

func readCString(file *os.File, offset int64) string {

	file.Seek(offset, os.SEEK_SET)

	s, _, _ := countInitiatedASCII(file)
	return s
}

func ReadBytesFrom(file *os.File, offset int64, size int64) []byte {

	file.Seek(offset, os.SEEK_SET)

	b := make([]byte, size)
	binary.Read(file, binary.LittleEndian, &b)
	return b
}

func ReadUint8(file *os.File, offset int64) (uint8, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint8
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func ReadUint16be(file *os.File, offset int64) (uint16, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func ReadUint16le(file *os.File, offset int64) (uint16, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func ReadUint32be(file *os.File, offset int64) (uint32, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func ReadUint32le(file *os.File, offset int64) (uint32, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func (pl *ParsedLayout) PercentMapped(totalSize int64) float64 {

	mapped := 0
	for _, l := range pl.Layout {
		mapped += int(l.Length)
	}

	//	fmt.Println("total =", totalSize, "mapped=", mapped, "in ", len(pl.Layout), " layouts")
	//	os.Exit(1)
	pct := (float64(mapped) / float64(totalSize)) * 100
	return pct
}

func (pl *ParsedLayout) PercentUnmapped(totalSize int64) float64 {
	return 100 - pl.PercentMapped(totalSize)
}
