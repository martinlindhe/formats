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

	res.Layout = append(res.Layout, gifHeader(file))
	res.Layout = append(res.Layout, gifGlobalColorTable(file))

	gfxExt := gifGraphicsControlExtension(file)
	if gfxExt != nil {
		res.Layout = append(res.Layout, *gfxExt)
	}

	imgDescriptor := gifImageDescriptor(file)
	if imgDescriptor != nil {
		res.Layout = append(res.Layout, *imgDescriptor)
	}

	// 0x2b för nästa

	localColorTbl := gifLocalColorTable(file)
	if localColorTbl != nil {
		res.Layout = append(res.Layout, *localColorTbl)
	}

	res.Layout = append(res.Layout, gifImageData(file))

	res.Layout = append(res.Layout, gifTrailer(file))

	return &res
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

func gifLocalColorTable(file *os.File) *Layout {
	// XXX The local color table would always immediately follow an
	// image descriptor but will only be there if the local color table flag is set to 1

	// XXX not present in sample
	return nil
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

func gifGlobalColorTable(file *os.File) Layout {

	baseOffset := int64(0x0d)

	// XXX TODO The size of the global color table is determined by the value in the packed byte of the logical screen descriptor.

	return Layout{
		Offset: baseOffset,
		Length: 4 * 3,
		Info:   "global color table",
		Type:   Group,
		Childs: []Layout{
			// XXX dont hardcode
			Layout{Offset: baseOffset, Length: 3, Info: "color 0", Type: RGB},
			Layout{Offset: baseOffset + 3, Length: 3, Info: "color 1", Type: RGB},
			Layout{Offset: baseOffset + 6, Length: 3, Info: "color 2", Type: RGB},
			Layout{Offset: baseOffset + 9, Length: 3, Info: "color 3", Type: RGB},
		},
	}

}

func gifHeader(file *os.File) Layout {

	return Layout{
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
}
