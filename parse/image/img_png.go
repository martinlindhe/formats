package image

// handles PNG and MNG images
// STATUS: 80% PNG/APNG
// STATUS: 20% MNG. parsing gives up after first IEND, should continue...

import (
	"encoding/binary"
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"io"
	"os"
)

func PNG(file *os.File) (*parse.ParsedLayout, error) {

	if !isPNG(file) {
		return nil, nil
	}
	return parsePNG(file)
}

func isPNG(file *os.File) bool {

	b, err := getPNGHeader(file)
	if err != nil {
		return false
	}

	if (b[0] == 0x89 && b[1] == 'P') || // png
		(b[0] == 0x8a && b[1] == 'M') { // mng
		if b[2] == 'N' && b[3] == 'G' && b[4] == 0xd &&
			b[5] == 0xa && b[6] == 0x1a && b[7] == 0xa {
			return true
		}
	}

	return false
}

func parsePNG(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Image}

	b, err := getPNGHeader(file)
	if err != nil {
		return nil, err
	}
	fileType := "PNG"
	if b[1] == 'M' {
		fileType = "MNG"
	}

	fileHeader := parse.Layout{
		Offset: pos,
		Length: 8,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: 0, Length: 8, Info: "magic = " + fileType, Type: parse.Bytes},
		}}

	res.Layout = append(res.Layout, fileHeader)

	pos = 8

	chunks := []parse.Layout{}
	for {
		l := parse.Layout{
			Offset: pos,
			Length: 8,
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "length", Type: parse.Uint32be},
				{Offset: pos + 4, Length: 4, Info: "type", Type: parse.ASCII},
			},
		}
		chunkLength, err := parse.ReadUint32be(file, pos)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		typeCode, _, err := parse.ReadZeroTerminatedASCIIUntil(file, pos+4, 4)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		l.Info = "chunk " + typeCode
		pos += l.Length

		if typeCode == "IHDR" {
			if chunkLength != 13 {
				fmt.Println("warning: IHDR size must be 13")
			}
			l.Childs = append(l.Childs, []parse.Layout{
				{Offset: pos, Length: 4, Info: "width", Type: parse.Uint32be},
				{Offset: pos + 4, Length: 4, Info: "height", Type: parse.Uint32be},
				{Offset: pos + 8, Length: 1, Info: "bit depth", Type: parse.Uint8},
				{Offset: pos + 9, Length: 1, Info: "color type", Type: parse.Uint8},
				{Offset: pos + 10, Length: 1, Info: "compression method", Type: parse.Uint8}, // XXX show meaning of value
				{Offset: pos + 11, Length: 1, Info: "filter method", Type: parse.Uint8},
				{Offset: pos + 12, Length: 1, Info: "interlace method", Type: parse.Uint8},
			}...)
		} else {
			l.Childs = append(l.Childs, []parse.Layout{
				{Offset: pos, Length: int64(chunkLength), Info: typeCode + " data", Type: parse.Bytes},
			}...)
		}

		pos += int64(chunkLength)
		l.Length += int64(chunkLength)

		l.Childs = append(l.Childs, []parse.Layout{
			{Offset: pos, Length: 4, Info: "crc", Type: parse.Uint32be},
		}...)
		l.Length += 4
		pos += 4

		chunks = append(chunks, l)

		if typeCode == "IEND" {
			break
		}
	}

	res.Layout = append(res.Layout, chunks...)

	return &res, nil
}

func getPNGHeader(file *os.File) ([8]byte, error) {
	file.Seek(0, os.SEEK_SET)

	var b [8]uint8
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
}
