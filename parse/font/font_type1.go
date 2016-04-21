package font

// PostScript Type 1 font
// https://en.wikipedia.org/wiki/PostScript_fonts

/*
0	string		%!PS-AdobeFont-1.	PostScript Type 1 font text
>20	string		>\0			(%s)
6	string		%!PS-AdobeFont-1.	PostScript Type 1 font program data
0	string		%!FontType1	PostScript Type 1 font program data
6	string		%!FontType1	PostScript Type 1 font program data
0	string		%!PS-Adobe-3.0\ Resource-Font	PostScript Type 1 font text
*/

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func PSType1(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPSType1(c.File) {
		return nil, nil
	}
	return parsePSType1(c.File, c.ParsedLayout)
}

func isPSType1(file *os.File) bool {

	chk := getPS1Type(file)
	if chk == 0 {
		return false
	}
	return true
}

func parsePSType1(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	layout := parse.Layout{}
	fileType := getPS1Type(file)
	pl.FileKind = parse.Font
	pl.Layout = []parse.Layout{layout}

	if fileType == FontText {
		pl.FormatName = "type1 text"
		layout = parse.Layout{
			Offset: pos,
			Length: 17, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 17, Info: "magic", Type: parse.ASCII},
			}}
	} else {
		pl.FormatName = "type1 program data"
		layout = parse.Layout{
			Offset: pos,
			Length: 17, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 6, Info: "unknown", Type: parse.Bytes},
				{Offset: pos + 6, Length: 17, Info: "magic", Type: parse.ASCII},
			}}
	}

	return &pl, nil
}

func getPS1Type(file *os.File) ps1Type {

	expected := "%!PS-AdobeFont-1."
	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, 0, len(expected))
	if s == expected {
		return FontText
	}
	s, _, _ = parse.ReadZeroTerminatedASCIIUntil(file, 6, len(expected))
	if s == expected {
		return FontProgramData
	}
	return 0
}

type ps1Type int

const (
	FontText = 1 + iota
	FontProgramData
)
