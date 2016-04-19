package exe

// https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func JavaClass(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isJavaClass(&hdr) {
		return nil, nil
	}
	return parseJavaClass(file, pl)
}

func isJavaClass(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] == 0xca && b[1] == 0xfe && b[2] == 0xba && b[3] == 0xbe {
		return true
	}
	return false
}

func parseJavaClass(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32be},
		}}}

	return &pl, nil
}
