package windows

// ???, not the RIFF WAVE format
// found on Windows 10 Windows/WinSxS/amd64_microsoft-windows-t..peech-en-gb-onecore_31bf3856ad364e35_10.0.10240.16384_none_e1ad0a33c01f1b40/M2057Sarah.voiceAssistant.WVE
// extensions: .wve

// STATUS: 1%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WAVE(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isWAVE(c.Header) {
		return nil, nil
	}
	return parseWAVE(c.File, c.ParsedLayout)
}

func isWAVE(b []byte) bool {

	// XXX just guessing
	if b[0] != 'W' || b[1] != 'A' || b[2] != 'V' || b[3] != 'E' {
		return false
	}
	return true
}

func parseWAVE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
