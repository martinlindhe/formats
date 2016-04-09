package parse

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

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

func zeroTerminatedASCII(r io.Reader) (string, error) {

	var c byte
	s := ""
	for {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return s, err
		}
		if c == 0 {
			break
		}
		s += string(c)
	}
	return s, nil
}

func readBytesFrom(file *os.File, offset int64, size int64) []byte {

	file.Seek(offset, os.SEEK_SET)

	b := make([]byte, size)
	binary.Read(file, binary.LittleEndian, &b)
	return b
}

func readUint8(file *os.File, offset int64) (uint8, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint8
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func readUint16be(file *os.File, offset int64) (uint16, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint16
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func readUint32be(file *os.File, offset int64) (uint32, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.BigEndian, &b)
	return b, nil
}

func readUint32le(file *os.File, offset int64) (uint32, error) {

	file.Seek(offset, os.SEEK_SET)
	var b uint32
	binary.Read(file, binary.LittleEndian, &b)
	return b, nil
}

func fileSize(file *os.File) int64 {

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}
