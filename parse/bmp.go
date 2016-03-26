package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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
	/*
		case 12: // OS/2 V1 - BITMAPCOREHEADER
			layout.Info = "bitmap info header, OS/2 V1"
			header.Nodes.Add(ParseOS2Version1Header(headerLen.offset + headerLen.length))
	*/
	case 40: // Windows V3 BITMAPINFOHEADER
		layout.Info = "bmp info header V3 Win"
		layout.Childs = parseBMPVersion3Header(file, 40)
		/*
			case 64: //OS/2 V2
				layout.Info = "bitmap info header, OS/2 V2"
				var v3 = ParseVersion3Header(headerLen.offset + headerLen.length)
				header.Nodes.Add(v3)

				var os2_v2 = ParseOS2Version2Header(headerLen.offset + headerLen.length + v3.length)
				header.Nodes.Add(os2_v2)

			case 108: //Windows V4 - BITMAPV4HEADER
				layout.Info = "bitmap info header, Windows V4"
				var v3 = ParseVersion3Header(headerLen.offset + headerLen.length)
				header.Nodes.Add(v3)
				var v4 = ParseVersion4Header(headerLen.offset + headerLen.length + v3.length)
				header.Nodes.Add(v4)

			case 124: //Windows V5 - BITMAPV5HEADER
				layout.Info = "bitmap info header, Windows V5"
				var v3 = ParseVersion3Header(headerLen.offset + headerLen.length)
				header.Nodes.Add(v3)
				var v4 = ParseVersion4Header(headerLen.offset + headerLen.length + v3.length)
				header.Nodes.Add(v4)
				var v5 = ParseVersion5Header(headerLen.offset + headerLen.length + v3.length + v4.length)
				header.Nodes.Add(v5)
		*/
	default:
		return layout, fmt.Errorf("unrecognized header size %d", infoHdrSize)
	}

	return layout, nil
}

func parseBMPVersion3Header(file *os.File, baseOffset int64) []Layout {

	return []Layout{
		Layout{Offset: baseOffset, Length: 4, Type: Uint32le, Info: "info header size V3"},
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

	/*
		switch (compressionValue) {
		case 0:
			compression.Text += "RGB";
			break;
		case 1:
			compression.Text += "RLE8";
			break;
		case 2:
			compression.Text += "RLE4";
			break;
		case 3:
			compression.Text += "BITFIELDS";
			break;
		case 4:
			compression.Text += "JPEG";
			break;
		case 5:
			compression.Text += "PNG";
			break;
		default:
			throw new Exception("unknown " + compressionValue);
		}
	*/
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
