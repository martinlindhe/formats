package windows

// STATUS: 1%
// Extensions: .mdl
// ??? found on Windows 10, Windows/WinSxS/amd64_windows-media-faceanalysis_31bf3856ad364e35_10.0.10240.16384_none_8cb86b56f21902dd/FaceAnalysisColor.mdl

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// LDMF parses the ldmf format
func LDMF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLDMF(c.Header) {
		return nil, nil
	}
	return parseLDMF(c.File, c.ParsedLayout)
}

func isLDMF(b []byte) bool {

	if b[0] != 'L' || b[1] != 'D' || b[2] != 'M' || b[3] != 'F' {
		return false
	}
	return true
}

func parseLDMF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
