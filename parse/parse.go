package parse

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

func zeroTerminatedASCII(r io.Reader) (string, error) {

	var c byte
	s := ""
	for {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return "", err
		}
		if c == 0 {
			break
		}
		s += string(c)
	}
	return s, nil
}

func fileSize(file *os.File) int64 {

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}
