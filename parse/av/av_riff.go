package av

// RIFF container format by Microsoft and IBM. Used in WAV, AVI, WebP
// https://en.wikipedia.org/wiki/Resource_Interchange_File_Format

// https://en.wikipedia.org/wiki/WebP
// https://en.wikipedia.org/wiki/Audio_Video_Interleave
// https://en.wikipedia.org/wiki/WAV
// Extensions: .webp, .wav, .avi, .ani, .rmi

// STATUS: 2%

import (
	"fmt"

	"github.com/martinlindhe/formats/parse"
)

// RIFF parses the riff format
func RIFF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRIFF(c.Header) {
		return nil, nil
	}
	return parseRIFF(c)
}

func isRIFF(b []byte) bool {

	if b[0] != 'R' || b[1] != 'I' || b[2] != 'F' || b[3] != 'F' {
		return false
	}
	return true
}

var (
	riffFormatNames = map[string]string{
		"WEBP": "riff-webp",
		"WAVE": "riff-wav",
		"AVI ": "riff-avi",
	}
	riffMimeTypes = map[string]string{
		"WEBP": "image/webp",
		"WAVE": "audio/x-wav",
		"AVI ": "video/avi",
		"ACON": "riff-ani",
		"RMID": "riff-midi",
	}
)

func parseRIFF(c *parse.Checker) (*parse.ParsedLayout, error) {

	idTag := string(c.Header[8 : 8+4])

	if val, ok := riffMimeTypes[idTag]; ok {
		c.ParsedLayout.MimeType = val
	}

	if val, ok := riffFormatNames[idTag]; ok {
		c.ParsedLayout.FormatName = val
	} else {
		fmt.Println("error: unknown riff id tag:", idTag)
	}

	pos := int64(0)
	c.ParsedLayout.FileKind = parse.AudioVideo
	c.ParsedLayout.Layout = []parse.Layout{{
		Offset: pos,
		Length: 12, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 4, Info: "length", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "type tag", Type: parse.ASCII},
		}}}

	return &c.ParsedLayout, nil
}
