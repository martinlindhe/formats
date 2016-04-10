package parse

// STATUS 1% , XXX

import (
	"encoding/binary"
	"os"
)

func ISO(file *os.File) (*ParsedLayout, error) {

	if !isISO(file) {
		return nil, nil
	}
	return parseISO(file)
}

func isISO(file *os.File) bool {

	/* XXX
	   if (BaseStream.Length < 0x8000 + 100)
	       return false;
	*/
	file.Seek(0x8000, os.SEEK_SET)

	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 1 || b[1] != 'C' || b[2] != 'D' {
		return false
	}

	return true
}

func parseISO(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0x8000,
		Length: 3,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0x8000, Length: 3, Info: "magic", Type: Bytes},
		},
	})

	return &res, nil
}
