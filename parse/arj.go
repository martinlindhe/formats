package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ARJ(file *os.File) *ParsedLayout {

	if !isARJ(file) {
		return nil
	}

	res := ParsedLayout{}
	fmt.Println("XXX TODO parse ARJ")
	return &res
}

func isARJ(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	r := io.Reader(file)
	var b [2]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}
	return b[0] == 0x60 && b[1] == 0xea
}
