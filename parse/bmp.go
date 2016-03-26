package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func BMP(file *os.File) []Layout {

	res := []Layout{}
	if !isBmp(file) {
		return nil
	}

	fmt.Println("XXX TODO parse BMP")

	return res
}

func isBmp(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)

	r := io.Reader(file)

	var b [2]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}

	return b[0] == 'B' && b[1] == 'M'
}
