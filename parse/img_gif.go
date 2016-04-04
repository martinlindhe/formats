package parse

import (
	"encoding/binary"
	"io"
	"os"
)

// STATUS xxx, incompelete. dont map the image data

func GIF(file *os.File) *ParsedLayout {

	if !isGIF(file) {
		return nil
	}
	return parseGIF(file)
}

func isGIF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	r := io.Reader(file)
	var b [5]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'G' || b[1] != 'I' || b[2] != 'F' || b[3] != '8' {
		return false
	}

	if b[4] == '7' || b[4] == '9' {
		return true
	}

	return false
}

func parseGIF(file *os.File) *ParsedLayout {

	res := ParsedLayout{}

	fileHeader := Layout{
		Offset: 0,
		Length: 13,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 3, Info: "magic", Type: ASCII},

			Layout{Offset: 3, Length: 3, Info: "version", Type: ASCII},

			Layout{Offset: 6, Length: 2, Info: "width", Type: Uint16le},

			Layout{Offset: 8, Length: 2, Info: "height", Type: Uint16le},

			// Packed contains the following four subfields of data (bit 0 is the least significant bit, or LSB):
			//    Bits 0-2    Size of the Global Color Table
			//    Bit 3   Color Table Sort Flag
			//    Bits 4-6    Color Resolution
			/*
				var Packed = ScreenHeight.RelativeToByte("Packed");
				byte PackedValue = ReadByte(Packed.offset);
				int bpp = (PackedValue & 0x7) + 1;
				Log("bpp = " + bpp);
				int colorRes = (PackedValue & 0x70) >> 4;
				Log("color Res = " + colorRes);
			*/
			Layout{Offset: 10, Length: 1, Info: "packed", Type: Uint8}, // XXX bitmask decode bpp and color resolution

			Layout{Offset: 11, Length: 1, Info: "background color", Type: Uint8},

			Layout{Offset: 12, Length: 1, Info: "aspect ratio", Type: Uint8},
		},
	}

	res.Layout = append(res.Layout, fileHeader)
	/*

	   # typedef struct _GifColorTable {       global and local colortables!
	   #  BYTE Red;          // Red Color Element
	   #  BYTE Green;        // Green Color Element
	   #  BYTE Blue;         // Blue Color Element
	   # } GIFCOLORTABLE;

	   # next,

	   # typedef struct _GifImageDescriptor {
	   #  BYTE Separator;    // Image Descriptor identifier
	   #  WORD Left;         // X position of image on the display
	   #  WORD Top;          // Y position of image on the display
	   #  WORD Width;        // Width of the image in pixels
	   #  WORD Height;       // Height of the image in pixels
	   #  BYTE Packed;       // Image and Color Table Data Information
	   # } GIFIMGDESC;

	   # finally

	   # image data XXX
	*/
	return &res
}
