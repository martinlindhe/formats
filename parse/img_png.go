package parse

// handles PNG and MNG images
// STATUS: 80% PNG/APNG
// STATUS: 20% MNG. parsing gives up after first IEND, should continue...

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func PNG(file *os.File) (*ParsedLayout, error) {

	if !isPNG(file) {
		return nil, nil
	}
	return parsePNG(file)
}

func getPNGHeader(file *os.File) ([8]byte, error) {
	file.Seek(0, os.SEEK_SET)

	var b [8]uint8
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
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

func parsePNG(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)
	res := ParsedLayout{
		FileKind: Image,
	}

	b, err := getPNGHeader(file)
	if err != nil {
		return nil, err
	}
	fileType := "PNG"
	if b[1] == 'M' {
		fileType = "MNG"
	}

	fileHeader := Layout{
		Offset: offset,
		Length: 8,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			{Offset: 0, Length: 8, Info: "magic = " + fileType, Type: Bytes},
		},
	}

	res.Layout = append(res.Layout, fileHeader)

	offset = 8

	chunks := []Layout{}
	for {
		l := Layout{
			Offset: offset,
			Length: 8,
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 4, Info: "length", Type: Uint32be},
				{Offset: offset + 4, Length: 4, Info: "type", Type: ASCII},
			},
		}
		chunkLength, err := readUint32be(file, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		typeCode, err := knownLengthASCII(file, offset+4, 4)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		l.Info = "chunk " + typeCode
		offset += l.Length

		if typeCode == "IHDR" {
			if chunkLength != 13 {
				fmt.Println("warning: IHDR size must be 13")
			}
			l.Childs = append(l.Childs, []Layout{
				{Offset: offset, Length: 4, Info: "width", Type: Uint32be},
				{Offset: offset + 4, Length: 4, Info: "height", Type: Uint32be},
				{Offset: offset + 8, Length: 1, Info: "bit depth", Type: Uint8},
				{Offset: offset + 9, Length: 1, Info: "color type", Type: Uint8},
				{Offset: offset + 10, Length: 1, Info: "compression method", Type: Uint8}, // XXX show meaning of value
				{Offset: offset + 11, Length: 1, Info: "filter method", Type: Uint8},
				{Offset: offset + 12, Length: 1, Info: "interlace method", Type: Uint8},
			}...)
		} else {
			l.Childs = append(l.Childs, []Layout{
				{Offset: offset, Length: int64(chunkLength), Info: typeCode + " data", Type: Bytes},
			}...)
		}

		offset += int64(chunkLength)
		l.Length += int64(chunkLength)

		l.Childs = append(l.Childs, []Layout{
			{Offset: offset, Length: 4, Info: "crc", Type: Uint32be},
		}...)
		l.Length += 4
		offset += 4

		chunks = append(chunks, l)

		if typeCode == "IEND" {
			break
		}
	}

	res.Layout = append(res.Layout, chunks...)

	return &res, nil
}
