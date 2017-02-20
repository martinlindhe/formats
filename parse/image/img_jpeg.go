package image

// XXX recognize jpeg with exif ... samples/images/jpg/jpeg_003_exif_fujifilm-finepix40i.jpg is unrecognized

// STATUS: 80%

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/martinlindhe/formats/parse"
)

const (
	jpegSOI  = 0xd8
	jpegEOI  = 0xd9
	jpegSOS  = 0xda
	jpegDQT  = 0xdb
	jpegAPP0 = 0xe0
	jpegAPP1 = 0xe1
	jpegCOM  = 0xfe
)

var (
	jpegMarkers = map[byte]string{
		0xc0:     "baseline DCT (SOF0)",
		0xc1:     "extended sequential DCT (SOF1)",
		0xc2:     "progressive DCT (SOF2)",
		0xc3:     "lossless (SOF3)",
		0xc4:     "huffman table (DHT)",
		jpegSOI:  "start of image (SOI)",
		jpegEOI:  "end of image (EOI)",
		jpegSOS:  "start of scan (SOS)",
		jpegDQT:  "quantization table (DQT)",
		jpegAPP0: "APP0",
		jpegAPP1: "APP1",
		jpegCOM:  "comment (COM)",
	}
)

// JPEG parses the jpeg format
func JPEG(c *parse.Checker) (*parse.ParsedLayout, error) {
	if !isJPEG(c.Header) {
		return nil, nil
	}
	return parseJPEG(c.File, c.ParsedLayout)
}

func isJPEG(b []byte) bool {
	if b[0] != 0xff || b[1] != 0xd8 {
		return false
	}
	if b[6] != 'J' || b[7] != 'F' || b[8] != 'I' || b[9] != 'F' || b[10] != 0 {
		return false
	}
	return true
}

func parseJPEG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {
	pos := int64(0)
	pl.FileKind = parse.Image
	pl.MimeType = "image/jpeg"

	for {
		magic, _ := parse.ReadUint8(file, pos)
		marker, _ := parse.ReadUint8(file, pos+1)

		if magic != 0xff {
			fmt.Printf("jpeg parse error at %04x. expected ff, found %02x\n", pos, magic)
			break
		}

		if marker == jpegSOI {
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Length: 2,
				Info:   jpegMarkers[marker],
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16le},
				}})
			pos += 2
			continue
		}
		if marker == jpegEOI {
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Type:   parse.Group,
				Length: 2,
				Info:   jpegMarkers[marker],
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16le},
				}})
			break
		}

		if marker == jpegSOS {

			sos := parseJPEGSos(file, pos)
			pl.Layout = append(pl.Layout, sos)
			pos += sos.Length

			imgData := findJPEGImageData(file, pos)
			pl.Layout = append(pl.Layout, imgData)
			pos += imgData.Length
			continue
		}

		if marker == jpegAPP0 {
			app0 := parseJPEGApp0(file, pos)
			pl.Layout = append(pl.Layout, app0)
			pos += app0.Length
			continue
		}

		chunkLen, _ := parse.ReadUint16be(file, pos+2)

		dataType := parse.Bytes
		if marker == jpegCOM {
			dataType = parse.ASCII
		}

		pl.Layout = append(pl.Layout, parse.Layout{
			Offset: pos,
			Length: 2 + int64(chunkLen),
			Info:   jpegMarkers[marker],
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
				{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
				{Offset: pos + 4, Length: int64(chunkLen) - 2, Info: "data", Type: dataType},
			}})
		pos += 2 + int64(chunkLen)
	}

	return &pl, nil
}

func findJPEGImageData(file *os.File, pos int64) parse.Layout {

	dataStart := pos

	res := parse.Layout{
		Offset: dataStart,
		Info:   "image data",
		Type:   parse.Group}

	for {
		file.Seek(pos, os.SEEK_SET)

		var b uint16
		err := binary.Read(file, binary.BigEndian, &b)
		if err == io.EOF {
			fmt.Printf("error: jpeg EOF at pos %04x\n", pos)
			break
		}

		pos++
		marker := b & 0xff

		if b&0xff00 == 0xff00 && marker != 0 && (marker < 0xd0 || marker > 0xd8) {
			// eoi
			dataLen := pos - dataStart - 1
			res.Length = dataLen
			res.Childs = []parse.Layout{
				{Offset: dataStart, Length: dataLen, Info: "image data", Type: parse.Bytes},
			}
			break
		}
	}
	return res
}

func parseJPEGSos(file *os.File, pos int64) parse.Layout {

	components, _ := parse.ReadUint8(file, pos+4)
	chunk := parse.Layout{
		Offset: pos,
		Length: 5,
		Info:   jpegMarkers[jpegSOS],
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
			{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
			{Offset: pos + 4, Length: 1, Info: "color components", Type: parse.Uint8},
		}}
	pos += chunk.Length

	for i := 0; i < int(components); i++ {
		chunk.Childs = append(chunk.Childs, []parse.Layout{
			{Offset: pos, Length: 1, Info: "COMPSOS", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 4, Info: "dc table"},
				{Low: 4, Length: 4, Info: "ac table"},
			}},
		}...)
		chunk.Length++
		pos++
	}

	chunk.Childs = append(chunk.Childs, []parse.Layout{
		{Offset: pos, Length: 1, Info: "ss", Type: parse.Uint8},
		{Offset: pos + 1, Length: 1, Info: "se", Type: parse.Uint8},
		{Offset: pos + 2, Length: 1, Info: "a", Type: parse.Uint8, Masks: []parse.Mask{
			{Low: 0, Length: 4, Info: "al"},
			{Low: 4, Length: 4, Info: "ah"},
		}},
	}...)
	chunk.Length += 3

	return chunk
}

func parseJPEGApp0(file *os.File, pos int64) parse.Layout {

	return parse.Layout{
		Offset: pos,
		Length: 18,
		Info:   jpegMarkers[jpegAPP0],
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
			{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
			{Offset: pos + 4, Length: 5, Info: "identifier", Type: parse.ASCII},
			{Offset: pos + 9, Length: 2, Info: "revision", Type: parse.MajorMinor16be},
			{Offset: pos + 11, Length: 1, Info: "units used", Type: parse.Uint8},
			{Offset: pos + 12, Length: 2, Info: "width", Type: parse.Uint16be},
			{Offset: pos + 14, Length: 2, Info: "height", Type: parse.Uint16be},
			{Offset: pos + 16, Length: 1, Info: "horizontal pixels", Type: parse.Uint8},
			{Offset: pos + 17, Length: 1, Info: "vertical pixels", Type: parse.Uint8},
		}}
}
