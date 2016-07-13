package font

// https://en.wikipedia.org/wiki/PostScript_fonts

// STATUS: 1%

import (
	"github.com/martinlindhe/formats/parse"
)

// PSType1 parses the PostScript Type 1 font format
func PSType1(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPSType1(c.Header) {
		return nil, nil
	}
	return parsePSType1(c)
}

func isPSType1(b []byte) bool {

	chk := getPS1Type(b)
	if chk == 0 {
		return false
	}
	return true
}

func parsePSType1(c *parse.Checker) (*parse.ParsedLayout, error) {

	pos := int64(0)
	layout := parse.Layout{}
	fileType := getPS1Type(c.Header)
	c.ParsedLayout.FileKind = parse.Font

	if fileType == fontText {
		c.ParsedLayout.FormatName = "type1 text"
		layout = parse.Layout{
			Offset: pos,
			Length: 17, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 17, Info: "magic", Type: parse.ASCII},
			}}
	} else {
		c.ParsedLayout.FormatName = "type1 program data"
		layout = parse.Layout{
			Offset: pos,
			Length: 23, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 6, Info: "unknown", Type: parse.Bytes},
				{Offset: pos + 6, Length: 17, Info: "magic", Type: parse.ASCII},
			}}
	}

	c.ParsedLayout.Layout = []parse.Layout{layout}
	return &c.ParsedLayout, nil
}

func getPS1Type(b []byte) ps1Type {

	/* XXX
	0	string		%!PS-AdobeFont-1.	PostScript Type 1 font text
	>20	string		>\0			(%s)
	6	string		%!PS-AdobeFont-1.	PostScript Type 1 font program data
	0	string		%!FontType1	PostScript Type 1 font program data
	6	string		%!FontType1	PostScript Type 1 font program data
	0	string		%!PS-Adobe-3.0\ Resource-Font	PostScript Type 1 font text
	*/

	expected := "%!PS-AdobeFont-1."
	s := string(b[0:len(expected)])
	if s == expected {
		return fontText
	}
	s = string(b[6 : 6+len(expected)])
	if s == expected {
		return fontProgramData
	}
	return 0
}

type ps1Type int

const (
	fontText = 1 + iota
	fontProgramData
)
