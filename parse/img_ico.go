package parse

// Windows Icon / Cursor image resources
// STATUS 80%

// TODO icon_embedded_png_001.ico has embedded PNG in image data
// TODO decode non-png as "standard BMP image"...
// TODO offer "save to file" for the "image data" (bytes type)

import (
	"encoding/binary"
	"fmt"
	"os"
)

func ICO(file *os.File) (*ParsedLayout, error) {

	if !isICO(file) {
		return nil, nil
	}
	return parseICO(file)
}

func readIconHeader(file *os.File) ([3]uint16, error) {

	file.Seek(0, os.SEEK_SET)
	var b [3]uint16
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
}

func isICO(file *os.File) bool {

	b, _ := readIconHeader(file)
	if b[0] != 0 {
		return false
	}

	// 1 = icon, 2 = cursor
	if b[1] != 1 && b[1] != 2 {
		return false
	}

	return true
}

func parseICO(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}
	typeName := ""

	hdr, _ := readIconHeader(file)
	switch hdr[1] {
	case 1:
		typeName = "icon"
	case 2:
		typeName = "cursor"
	default:
		typeName = "unknown"
	}

	offset := int64(0)

	fileHeader := Layout{
		Offset: offset,
		Length: 6,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: Uint16le},
			Layout{Offset: 2, Length: 2, Info: "type = " + typeName, Type: Uint16le},
			Layout{Offset: 4, Length: 2, Info: "number of resources", Type: Uint16le},
		},
	}

	offset += fileHeader.Length

	numIcons := hdr[2]
	resourceEntryLength := int64(16)

	fmt.Println("parsing ", numIcons, " resources")

	res.Layout = append(res.Layout, fileHeader)

	for i := 0; i < int(numIcons); i++ {
		resNum := fmt.Sprintf("%d", i+1)
		resource := Layout{
			Offset: offset,
			Length: resourceEntryLength,
			Info:   "resource " + resNum + " header",
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: offset, Length: 1, Info: "width", Type: Uint8},
				Layout{Offset: offset + 1, Length: 1, Info: "height", Type: Uint8},
				Layout{Offset: offset + 2, Length: 1, Info: "max number of colors", Type: Uint8},
				Layout{Offset: offset + 3, Length: 1, Info: "reserved", Type: Uint8},
				Layout{Offset: offset + 4, Length: 2, Info: "planes", Type: Uint16le},
				Layout{Offset: offset + 6, Length: 2, Info: "bit count", Type: Uint16le},
				Layout{Offset: offset + 8, Length: 4, Info: "data size of resource " + resNum, Type: Uint32le},
				Layout{Offset: offset + 12, Length: 4, Info: "offset to resource " + resNum, Type: Uint32le},
			}}

		res.Layout = append(res.Layout, resource)
		fileHeader.Length += resourceEntryLength
		offset += resourceEntryLength
	}

	for i := 0; i < int(numIcons); i++ {
		resNum := fmt.Sprintf("%d", i+1)

		dataOffset, err := res.readUint32leFromInfo(file, "offset to resource "+resNum)
		if err != nil {
			return nil, err
		}
		dataSize, err := res.readUint32leFromInfo(file, "data size of resource "+resNum)
		if err != nil {
			return nil, err
		}

		res.Layout = append(res.Layout, Layout{
			Offset: int64(dataOffset),
			Type:   Group,
			Info:   "resource " + resNum + " data",
			Length: int64(dataSize),
			Childs: []Layout{
				Layout{Offset: int64(dataOffset), Length: int64(dataSize), Info: "image data", Type: Bytes},
			}})
	}

	return &res, nil
}
