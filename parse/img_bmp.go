package parse

// STATUS v1 and v2 dont map 100% of the files. v3,v4 and v5 seems mostly done

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var (
	bmpCompressions = map[uint32]string{
		0: "rgb",
		1: "rle8",
		2: "rle4",
		3: "bitfields",
		4: "jpeg",
		5: "png",
	}
	v2len = 64
	v3len = 40
	v4len = 68
)

func BMP(file *os.File) *ParsedLayout {

	if !isBMP(file) {
		return nil
	}
	return parseBMP(file)
}

func isBMP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	r := io.Reader(file)
	var b [2]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}
	return b[0] == 'B' && b[1] == 'M'
}

func parseBMP(file *os.File) *ParsedLayout {

	res := ParsedLayout{}

	fileHeader := Layout{
		Offset: 0,
		Length: 14,
		Info:   "bmp file header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: ASCII},
			Layout{Offset: 2, Length: 4, Info: "file size", Type: Uint32le},
			Layout{Offset: 6, Length: 4, Info: "reserved", Type: Uint32le},
			Layout{Offset: 10, Length: 4, Info: "offset to image data", Type: Uint32le},
		},
	}

	res.Layout = append(res.Layout, fileHeader)

	infoHeader, err := parseBMPInfoHeader(file)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	res.Layout = append(res.Layout, infoHeader)

	imageDataOffset := res.readUint32leFromInfo(file, "offset to image data")

	// body
	headerLen := fileHeader.Length + infoHeader.Length

	dataLayout := Layout{
		Offset: int64(imageDataOffset),
		Type:   Group,
		Info:   "image data",
		Length: fileSize(file) - headerLen,
		Childs: []Layout{
			Layout{Offset: int64(imageDataOffset), Length: fileSize(file) - headerLen, Type: Uint8, Info: "image data"},
		},
	}

	res.Layout = append(res.Layout, dataLayout)

	return &res
}

func readUint32le(reader io.Reader) (uint32, error) {

	var b uint32
	err := binary.Read(reader, binary.LittleEndian, &b)
	return b, err
}

func parseBMPInfoHeader(file *os.File) (Layout, error) {

	infoHeaderBase := int64(14)
	layout := Layout{
		Offset: infoHeaderBase,
		Type:   Group,
	}

	file.Seek(infoHeaderBase, os.SEEK_SET)

	infoHdrSize, err := readUint32le(file)
	if err != nil {
		return layout, err
	}

	layout.Length = int64(infoHdrSize)
	fmt.Println(layout.Length, "length")

	switch infoHdrSize {

	case 12: // OS/2 V1 - BITMAPCOREHEADER
		layout.Info = "bmp info header V1 OS/2"
		layout.Childs = parseBMPVersion1Header(file, infoHeaderBase)

	case 40: // Windows V3 - BITMAPINFOHEADER
		layout.Info = "bmp info header V3 Win"
		layout.Childs = parseBMPVersion3Header(file, infoHeaderBase)

	case 64: //OS/2 V2
		layout.Info = "bmp info header V2 OS/2"
		v3 := parseBMPVersion3Header(file, infoHeaderBase)
		v2 := parseBMPVersion2Header(file, infoHeaderBase+int64(v3len))
		layout.Childs = append(v3, v2...)

	case 108: //Windows V4 - BITMAPV4HEADER
		layout.Info = "bmp info header V4 Win"
		v3 := parseBMPVersion3Header(file, infoHeaderBase)
		v4 := parseBMPVersion4Header(file, infoHeaderBase+int64(v3len))
		layout.Childs = append(v3, v4...)

	case 124: //Windows V5 - BITMAPV5HEADER
		layout.Info = "bmp info header V5 Win"
		v3 := parseBMPVersion3Header(file, infoHeaderBase)
		v4 := parseBMPVersion4Header(file, infoHeaderBase+int64(v3len))
		v5 := parseBMPVersion5Header(file, infoHeaderBase+int64(v3len)+int64(v4len))
		layout.Childs = append(v3, v4...)
		layout.Childs = append(layout.Childs, v5...)

	default:
		return layout, fmt.Errorf("unrecognized header size %d", infoHdrSize)
	}

	return layout, nil
}

func parseBMPVersion1Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "info header size"},
		Layout{Offset: baseOffset + 4, Length: 2, Type: Uint16le, Info: "width"},
		Layout{Offset: baseOffset + 6, Length: 2, Type: Uint16le, Info: "height"},
		Layout{Offset: baseOffset + 8, Length: 2, Type: Uint16le, Info: "planes"},
		Layout{Offset: baseOffset + 10, Length: 2, Type: Uint16le, Info: "bpp"},
	}
}

func parseBMPVersion2Header(file *os.File, baseOffset int64) []Layout {
	fmt.Printf(" base offset is %x\n", baseOffset)

	return []Layout{
		Layout{Offset: baseOffset, Length: 2, Type: Uint16le, Info: "units"},
		Layout{Offset: baseOffset + 2, Length: 2, Type: Uint16le, Info: "reserved"},
		Layout{Offset: baseOffset + 4, Length: 2, Type: Uint16le, Info: "recording"},
		Layout{Offset: baseOffset + 6, Length: 2, Type: Uint16le, Info: "rendering"},
		Layout{Offset: baseOffset + 8, Length: 4, Type: Uint32le, Info: "size1"},
		Layout{Offset: baseOffset + 12, Length: 4, Type: Uint32le, Info: "size2"},
		Layout{Offset: baseOffset + 16, Length: 4, Type: Uint32le, Info: "color encoding"},
		Layout{Offset: baseOffset + 20, Length: 4, Type: Uint32le, Info: "identifier"},
	}
}

func parseBMPVersion3Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "info header size"},
		Layout{Offset: baseOffset + 4, Length: 4, Type: Uint32le, Info: "width"},
		Layout{Offset: baseOffset + 8, Length: 4, Type: Uint32le, Info: "height"},
		Layout{Offset: baseOffset + 12, Length: 2, Type: Uint16le, Info: "planes"},
		Layout{Offset: baseOffset + 14, Length: 2, Type: Uint16le, Info: "bpp"},
		Layout{Offset: baseOffset + 16, Length: 4, Type: Uint32le, Info: "compression"}, // XXX decode value
		Layout{Offset: baseOffset + 20, Length: 4, Type: Uint32le, Info: "size of picture"},
		Layout{Offset: baseOffset + 24, Length: 4, Type: Uint32le, Info: "horizontal resolution"},
		Layout{Offset: baseOffset + 28, Length: 4, Type: Uint32le, Info: "vertical resolution"},
		Layout{Offset: baseOffset + 32, Length: 4, Type: Uint32le, Info: "number of used colors"},
		Layout{Offset: baseOffset + 36, Length: 4, Type: Uint32le, Info: "number of important colors"},
	}
}

func parseBMPVersion4Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "red mask"},
		Layout{Offset: baseOffset + 4, Length: 4, Type: Uint32le, Info: "green mask"},
		Layout{Offset: baseOffset + 8, Length: 4, Type: Uint32le, Info: "blue mask"},
		Layout{Offset: baseOffset + 12, Length: 4, Type: Uint32le, Info: "alpha mask"},
		Layout{Offset: baseOffset + 16, Length: 4, Type: Uint32le, Info: "cs type"}, // XXX "BGRs" ???

		//TODO: parse & display CIEXYZTRIPLE endpoint data: FXPT2DOT30  X, Y, Z
		Layout{Offset: baseOffset + 20, Length: 3 * 4, Type: Uint8, Info: "ciexyz red"},
		Layout{Offset: baseOffset + 20 + (3 * 4), Length: 3 * 4, Type: Uint8, Info: "ciexyz green"},
		Layout{Offset: baseOffset + 20 + (3 * 4) + (3 * 4), Length: 3 * 4, Type: Uint8, Info: "ciexyz blue"},
		Layout{Offset: baseOffset + 20 + (3 * 4) + (3 * 4) + (3 * 4), Length: 4, Type: Uint32le, Info: "gamma red"},
		Layout{Offset: baseOffset + 20 + (3 * 4) + (3 * 4) + (3 * 4) + 4, Length: 4, Type: Uint32le, Info: "gamma green"},
		Layout{Offset: baseOffset + 20 + (3 * 4) + (3 * 4) + (3 * 4) + 8, Length: 4, Type: Uint32le, Info: "gamma blue"},
	}
}

func parseBMPVersion5Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "intent"},
		Layout{Offset: baseOffset + 4, Length: 4, Type: Uint32le, Info: "profile data"},
		Layout{Offset: baseOffset + 8, Length: 4, Type: Uint32le, Info: "profile size"},
		Layout{Offset: baseOffset + 12, Length: 4, Type: Uint32le, Info: "reserved"},
	}
}
