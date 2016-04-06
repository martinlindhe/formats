package parse

import (
	"encoding/binary"
	"fmt"
	"os"
)

// STATUS xxx, incompelete. need to redo the  parseGIF logic to detect optional chunks, and have a loop

var (
	gctToLengthMap = map[byte]int64{
		0: 2 * 3,
		1: 4 * 3,
		2: 8 * 3,
		3: 16 * 3,
		4: 32 * 3,
		5: 64 * 3,
		6: 128 * 3,
		7: 256 * 3,
	}
)

// Section indicators.
const (
	sExtension       = 0x21
	sImageDescriptor = 0x2C
	sTrailer         = 0x3B
)

func GIF(file *os.File) *ParsedLayout {

	if !isGIF(file) {
		return nil
	}
	return parseGIF(file)
}

func isGIF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)

	var b [5]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'G' || b[1] != 'I' || b[2] != 'F' || b[3] != '8' {
		return false
	}
	if b[4] != '7' && b[4] != '9' {
		return false
	}
	return true
}

func parseGIF(file *os.File) *ParsedLayout {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, gifHeader(file))
	res.Layout = append(res.Layout, gifLogicalDescriptor(file))

	sizeOfGCT := res.decodeBitfieldFromInfo(file, "size of the global color table")
	if gctByteLen, ok := gctToLengthMap[byte(sizeOfGCT)]; ok {
		res.Layout = append(res.Layout, gifGlobalColorTable(file, gctByteLen))
	}

	for {

		var b byte
		if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
			fmt.Println("error", err)
			return nil
		}

		switch b {
		case sExtension:
			// XXX All extension blocks begin with 21
			gfxExt := gifGraphicsControlExtension(file)
			if gfxExt != nil {
				res.Layout = append(res.Layout, *gfxExt)
			}

		case sImageDescriptor:
			imgDescriptor := gifImageDescriptor(file)
			if imgDescriptor != nil {
				res.Layout = append(res.Layout, *imgDescriptor)
			}

			// XXX directly after image descriptor, depending on some flag ???
			localColorTbl := gifLocalColorTable(file)
			if localColorTbl != nil {
				res.Layout = append(res.Layout, *localColorTbl)
			}
			res.Layout = append(res.Layout, gifImageData(file)) // XXX ?

		case sTrailer:
			res.Layout = append(res.Layout, gifTrailer(file)) // XXX value 0x3b
			return &res
		}
	}
}

func gifHeader(file *os.File) Layout {

	return Layout{
		Offset: 0,
		Length: 6,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 3, Info: "signature", Type: ASCII},
			Layout{Offset: 3, Length: 3, Info: "version", Type: ASCII},
		},
	}
}

func gifTrailer(file *os.File) Layout {
	baseOffset := int64(0x44) // XXX base is unknown ...

	res := Layout{
		Offset: baseOffset,
		Length: 1,
		Info:   "trailer",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: baseOffset, Length: 1, Info: "trailer", Type: Uint8},
		},
	}
	return res
}

func gifImageData(file *os.File) Layout {

	baseOffset := int64(0x2b) // XXX base is unknown ...
	length := int64(25)       // XXX length=

	res := Layout{
		Offset: baseOffset,
		Length: length,
		Info:   "image data",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: baseOffset, Length: length, Info: "image data", Type: Uint8},
		},
	}
	return res
}

func gifImageDescriptor(file *os.File) *Layout {

	baseOffset := int64(0x21) // XXX base is unknown ...

	res := Layout{
		Offset: baseOffset,
		Length: 10,
		Info:   "image descriptor",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: baseOffset, Length: 1, Info: "image separator", Type: Uint8},
			Layout{Offset: baseOffset + 1, Length: 2, Info: "image left", Type: Uint16le},
			Layout{Offset: baseOffset + 3, Length: 2, Info: "image top", Type: Uint16le},
			Layout{Offset: baseOffset + 5, Length: 2, Info: "image width", Type: Uint16le},
			Layout{Offset: baseOffset + 7, Length: 2, Info: "image height", Type: Uint16le},
			Layout{Offset: baseOffset + 9, Length: 1, Info: "packed #3", Type: Uint8},
		},
	}
	return &res
}

func gifGlobalColorTable(file *os.File, byteLen int64) Layout {

	baseOffset := int64(0x0d)

	childs := []Layout{}

	cnt := 0
	for i := int64(0); i < byteLen; i += 3 {
		cnt++
		childs = append(childs, Layout{Offset: baseOffset + i, Length: 3, Info: fmt.Sprintf("color %d", cnt), Type: RGB})
	}

	return Layout{
		Offset: baseOffset,
		Length: byteLen,
		Info:   "global color table",
		Type:   Group,
		Childs: childs,
	}
}

func gifLocalColorTable(file *os.File) *Layout {
	// XXX The local color table would always immediately follow an
	// image descriptor but will only be there if the local color table flag is set to 1

	// XXX not present in sample
	return nil
}

func gifGraphicsControlExtension(file *os.File) *Layout {

	// XXX this is optional, how do we detect if it is there?
	// return nil

	// XXX "first byte is extension introducer. All extension blocks begin with 21"
	// XXX "Finally we have the block terminator which is always 00. "

	baseOffset := int64(0x19) // XXX base is unknown ...

	res := Layout{
		Offset: baseOffset,
		Length: 8,
		Info:   "gfx control extension",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: baseOffset, Length: 1, Info: "extension introducer", Type: Uint8},
			Layout{Offset: baseOffset + 1, Length: 1, Info: "graphic control label", Type: Uint8},
			Layout{Offset: baseOffset + 2, Length: 1, Info: "byte size", Type: Uint8},
			Layout{Offset: baseOffset + 3, Length: 1, Info: "packed #2", Type: Uint8},
			Layout{Offset: baseOffset + 4, Length: 2, Info: "delay time", Type: Uint16le},
			Layout{Offset: baseOffset + 6, Length: 1, Info: "transparent color index", Type: Uint8},
			Layout{Offset: baseOffset + 7, Length: 1, Info: "block terminator", Type: Uint8},
		},
	}
	return &res
}

func gifLogicalDescriptor(file *os.File) Layout {
	base := int64(0x06)
	return Layout{
		Offset: base,
		Length: 7,
		Info:   "logical screen descriptor",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: base, Length: 2, Info: "width", Type: Uint16le},

			Layout{Offset: base + 2, Length: 2, Info: "height", Type: Uint16le},

			// Packed contains the following four subfields of data (bit 0 is the least significant bit, or LSB):
			//    Bits 0-2    Size of the Global Color Table
			//    Bit 3   Color Table Sort Flag
			//    Bits 4-6    Color Resolution

			// XXX bitmask decode bpp and color resolution:
			Layout{Offset: base + 4, Length: 1, Info: "packed", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 3, Info: "size of the global color table"}, // XXX extract to a var
				Mask{Low: 3, Length: 1, Info: "color table sort flag"},
				Mask{Low: 4, Length: 3, Info: "color resolution"},
				Mask{Low: 7, Length: 1, Info: "reserved xxx"},
			}},

			Layout{Offset: base + 5, Length: 1, Info: "background color", Type: Uint8},

			Layout{Offset: base + 6, Length: 1, Info: "aspect ratio", Type: Uint8},
		},
	}
}
