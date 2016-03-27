package parse

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
	v3len = 40
)

func BMP(file *os.File) *ParsedLayout {

	if !isBMP(file) {
		return nil
	}
	return parseBMP(file)
}

func parseBMP(file *os.File) *ParsedLayout {

	res := ParsedLayout{}

	fileHeader := Layout{
		Offset: 0,
		Length: 14,
		Info:   "bitmap file header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic (BMP image)", Type: ASCII},
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
		Type:   Uint8,
		Info:   "image data",
		Length: getFileSize(file) - headerLen,
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

	switch infoHdrSize {

	case 12: // OS/2 V1 - BITMAPCOREHEADER
		layout.Info = "bmp info header V1 OS/2"
		layout.Childs = parseBMPVersion1Header(file, int64(infoHdrSize))

	case 40: // Windows V3 - BITMAPINFOHEADER
		layout.Info = "bmp info header V3 Win"
		layout.Childs = parseBMPVersion3Header(file, int64(infoHdrSize))

	case 64: //OS/2 V2
		layout.Info = "bmp info header V2 OS/2"
		v3 := parseBMPVersion3Header(file, int64(infoHdrSize))
		v2 := parseBMPVersion2Header(file, int64(infoHdrSize)+int64(v3len))
		layout.Childs = append(v3, v2...)

	case 108: //Windows V4 - BITMAPV4HEADER
		layout.Info = "bmp info header V4 Win"
		v3 := parseBMPVersion3Header(file, int64(infoHdrSize))
		v4 := parseBMPVersion4Header(file, 60) // XXX what is base of v4 hdr
		layout.Childs = append(v3, v4...)

	case 124: //Windows V5 - BITMAPV5HEADER
		layout.Info = "bmp info header V5 Win"
		v3 := parseBMPVersion3Header(file, int64(infoHdrSize))
		v4 := parseBMPVersion4Header(file, 60) // XXX what is base of v4 hdr
		v5 := parseBMPVersion5Header(file, 80) // XXX what is base of v5 hdr
		layout.Childs = append(v3, v4...)
		layout.Childs = append(layout.Childs, v5...)

	default:
		return layout, fmt.Errorf("unrecognized header size %d", infoHdrSize)
	}

	return layout, nil
}

// v1 = 12 byte header
func parseBMPVersion1Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "info header size"},
		Layout{Offset: baseOffset + 4, Length: 2, Type: Uint16le, Info: "width"},
		Layout{Offset: baseOffset + 6, Length: 2, Type: Uint16le, Info: "height"},
		Layout{Offset: baseOffset + 8, Length: 2, Type: Uint16le, Info: "planes"},
		Layout{Offset: baseOffset + 10, Length: 2, Type: Uint16le, Info: "bpp"},
	}
}

// v2 = 24 byte header
func parseBMPVersion2Header(file *os.File, baseOffset int64) []Layout {

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

// v3 = 40 byte header
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

// v4 = 68 byte header
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

// v5 = 16 byte header
func parseBMPVersion5Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "intent"},
		Layout{Offset: baseOffset + 4, Length: 4, Type: Uint32le, Info: "profile data"},
		Layout{Offset: baseOffset + 8, Length: 4, Type: Uint32le, Info: "profile size"},
		Layout{Offset: baseOffset + 12, Length: 4, Type: Uint32le, Info: "reserved"},
	}
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
