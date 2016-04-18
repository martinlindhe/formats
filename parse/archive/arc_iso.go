package archive

// STATUS 1% , XXX

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func ISO(file *os.File) (*parse.ParsedLayout, error) {

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

func parseISO(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0x8000)

	res := parse.ParsedLayout{
		FileKind: parse.Archive,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 3,
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 3, Info: "magic", Type: parse.Bytes},
			}}}}

	return &res, nil
}
