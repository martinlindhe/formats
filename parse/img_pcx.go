package parse

// STATUS 1%

import (
	"encoding/binary"
	"os"
)

var (
	pcxPaletteType = map[uint16]string{
		1: "color",
		2: "grayscale",
	}
	pcxVersions = map[uint8]string{
		0: "2.5",
		2: "2.8 w/ palette",
		3: "2.8 w/out palette",
		5: "3.0 or better",
	}
)

func PCX(file *os.File) (*ParsedLayout, error) {

	if !isPCX(file) {
		return nil, nil
	}
	return parsePCX(file)
}

func isPCX(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 0xa {
		return false
	}
	if b[1] != 0 && b[1] != 2 && b[1] != 3 && b[1] != 5 {
		return false
	}
	return true
}

func parsePCX(file *os.File) (*ParsedLayout, error) {
	offset := int64(0)

	version, _ := readUint8(file, offset+1)
	versionName := "?"
	if val, ok := pcxVersions[version]; ok {
		versionName = val
	}

	paletteType, _ := readUint16le(file, offset+68)
	paletteTypeName := "?"
	if val, ok := pcxPaletteType[paletteType]; ok {
		paletteTypeName = val
	}

	fileLen := fileSize(file)

	res := ParsedLayout{
		FileKind: Image,
		Layout: []Layout{{
			Offset: offset,
			Length: 128, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 1, Info: "magic", Type: Uint8},
				{Offset: offset + 1, Length: 1, Info: "version = " + versionName, Type: Uint8},
				{Offset: offset + 2, Length: 1, Info: "encoding", Type: Uint8},
				{Offset: offset + 3, Length: 1, Info: "bits per plane", Type: Uint8},
				{Offset: offset + 4, Length: 2, Info: "x min", Type: Uint16le},
				{Offset: offset + 6, Length: 2, Info: "y min", Type: Uint16le},
				{Offset: offset + 8, Length: 2, Info: "x max", Type: Uint16le},
				{Offset: offset + 10, Length: 2, Info: "y max", Type: Uint16le},
				{Offset: offset + 12, Length: 2, Info: "vertical dpi", Type: Uint16le},
				{Offset: offset + 14, Length: 2, Info: "horizontal dpi", Type: Uint16le},
				{Offset: offset + 16, Length: 48, Info: "palette", Type: Bytes},
				{Offset: offset + 64, Length: 1, Info: "reserved", Type: Uint8},
				{Offset: offset + 65, Length: 1, Info: "color planes", Type: Uint8},
				{Offset: offset + 66, Length: 2, Info: "bytes per plane line", Type: Uint16le},
				{Offset: offset + 68, Length: 2, Info: "palette type = " + paletteTypeName, Type: Uint16le},
				{Offset: offset + 70, Length: 2, Info: "hScrSize", Type: Uint16le},
				{Offset: offset + 72, Length: 2, Info: "vScrSize", Type: Uint16le},
				{Offset: offset + 74, Length: 54, Info: "padding", Type: Bytes}, // XXX may be 56 byte if horiz dpi is absent
			}}, {
			Offset: offset + 128,
			Length: fileLen - 128,
			Info:   "image data",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset + 128, Length: fileLen - 128, Info: "image data", Type: Bytes},
			}}}}

	return &res, nil
}
