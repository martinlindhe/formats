package image

// Adobe Photoshop Document
// https://en.wikipedia.org/wiki/Adobe_Photoshop#File_format
// Extensions: .psd

// NOTE about the mime-type:
// vnd.adobe.photoshop was registered by Adobe: https://www.iana.org/assignments/media-types/image/vnd.adobe.photoshop
// however, Adobe seems to be using a varity of mime types:
// image/photoshop, image/x-photoshop, image/psd, application/photoshop, application/psd, application/x-photoshop

// STATUS: 2%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	psdColorModes = map[uint16]string{
		0: "bitmap",
		1: "grayscale",
		2: "indexed",
		3: "RGB",
		4: "CMYK",
		7: "multichannel",
		8: "duotone",
		9: "lab",
	}
)

func PSD(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPSD(c.Header) {
		return nil, nil
	}
	return parsePSD(c.File, c.ParsedLayout)
}

func isPSD(b []byte) bool {

	if b[0] != '8' || b[1] != 'B' || b[2] != 'P' || b[3] != 'S' {
		return false
	}

	// version: uint16be
	if b[4] != 0 || b[5] != 1 {
		if b[5] == 2 {
			fmt.Println("TODO: psd version 2 file = 'big' file")
		}
		return false
	}

	return true
}

func parsePSD(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	colorMode, _ := parse.ReadToMap(file, parse.Uint16be, pos+24, psdColorModes)
	pl.FileKind = parse.Image
	pl.MimeType = "vnd.adobe.photoshop"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 26, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 2, Info: "version", Type: parse.Uint16be},
			{Offset: pos + 6, Length: 6, Info: "reserved", Type: parse.Bytes},
			{Offset: pos + 12, Length: 2, Info: "channels", Type: parse.Uint16be},
			{Offset: pos + 14, Length: 4, Info: "height", Type: parse.Uint32be},
			{Offset: pos + 18, Length: 4, Info: "width", Type: parse.Uint32be},
			{Offset: pos + 22, Length: 2, Info: "bits per channel", Type: parse.Uint16be},
			{Offset: pos + 24, Length: 2, Info: "color mode = " + colorMode, Type: parse.Uint16be},
		}}}

	return &pl, nil
}
