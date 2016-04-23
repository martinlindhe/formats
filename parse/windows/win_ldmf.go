package windows

// ??? found on Windows 10, Windows/WinSxS/amd64_windows-media-faceanalysis_31bf3856ad364e35_10.0.10240.16384_none_8cb86b56f21902dd/FaceAnalysisColor.mdl

// Extensions: .mdl

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func LDMF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

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
