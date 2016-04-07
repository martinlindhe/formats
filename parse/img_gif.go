package parse

import (
	"encoding/binary"
	"fmt"
	"os"
)

// STATUS wip

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

// Extensions.
const (
	eText           = 0x01 // Plain Text
	eGraphicControl = 0xF9 // Graphic Control
	eComment        = 0xFE // Comment
	eApplication    = 0xFF // Application
)

// misc
const (
	imgDescriptorLen = 10
)

func GIF(file *os.File) (*ParsedLayout, error) {

	if !isGIF(file) {
		return nil, nil
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

func parseGIF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, gifHeader(file))
	res.Layout = append(res.Layout, gifLogicalDescriptor(file))

	sizeOfGCT := res.decodeBitfieldFromInfo(file, "size of the global color table")
	if gctByteLen, ok := gctToLengthMap[byte(sizeOfGCT)]; ok {
		res.Layout = append(res.Layout, gifGlobalColorTable(file, gctByteLen))
	}

	for {

		offset, _ := file.Seek(0, os.SEEK_CUR)

		var b byte
		if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
			return nil, err
		}

		switch b {
		case sExtension:
			gfxExt := gifExtension(file, offset)
			if gfxExt != nil {
				res.Layout = append(res.Layout, *gfxExt)
			}

		case sImageDescriptor:
			imgDescriptor := gifImageDescriptor(file, offset)
			if imgDescriptor != nil {
				res.Layout = append(res.Layout, *imgDescriptor)
			}
			/*
				// XXX directly after image descriptor, depending on some flag ???
				localColorTbl := gifLocalColorTable(file, offset+imgDescriptorLen)
				if localColorTbl != nil {
					res.Layout = append(res.Layout, *localColorTbl)
				}
			*/
			imgData, err := gifImageData(file, offset+imgDescriptorLen)
			if err != nil {
				return nil, err
			}
			res.Layout = append(res.Layout, *imgData)

		case sTrailer:
			res.Layout = append(res.Layout, gifTrailer(file, offset))
			return &res, nil
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

func gifImageDescriptor(file *os.File, baseOffset int64) *Layout {

	res := Layout{
		Offset: baseOffset,
		Length: imgDescriptorLen,
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

func gifExtension(file *os.File, baseOffset int64) *Layout {

	var extType byte
	if err := binary.Read(file, binary.LittleEndian, &extType); err != nil {
		fmt.Println("error", err)
		return nil
	}

	typeSpecific := []Layout{}
	typeInfo := ""

	size := int64(0)
	switch extType {
	case eText:
		size = 13
		typeInfo = "text"

	case eGraphicControl:
		size = 7
		typeInfo = "graphic control"
		typeSpecific = []Layout{
			Layout{Offset: baseOffset + 2, Length: 1, Info: "byte size", Type: Uint8},
			Layout{Offset: baseOffset + 3, Length: 1, Info: "packed #2", Type: Uint8},
			Layout{Offset: baseOffset + 4, Length: 2, Info: "delay time", Type: Uint16le},
			Layout{Offset: baseOffset + 6, Length: 1, Info: "transparent color index", Type: Uint8},
			Layout{Offset: baseOffset + 7, Length: 1, Info: "block terminator", Type: Uint8},
		}

	case eComment:
		// nothing to do but read the data.
		typeInfo = "comment"

	case eApplication:
		typeInfo = "application"
		var lenByte byte
		if err := binary.Read(file, binary.LittleEndian, &lenByte); err != nil {
			fmt.Println("error", err)
			return nil
		}

		size = int64(lenByte)

		typeSpecific = []Layout{
			Layout{Offset: baseOffset + 2, Length: 1, Info: "byte size", Type: Uint8},
			Layout{Offset: baseOffset + 3, Length: size, Info: "data", Type: Uint8},
		}

		/*
			// Application Extension with "NETSCAPE2.0" as string and 1 in data means
			// this extension defines a loop count.
			if extension == eApplication && string(d.tmp[:size]) == "NETSCAPE2.0" {
				n, err := d.readBlock()
				if n == 0 || err != nil {
					return err
				}
				if n == 3 && d.tmp[0] == 1 {
					d.loopCount = int(d.tmp[1]) | int(d.tmp[2])<<8
				}
			}
			for {
				n, err := d.readBlock()
				if n == 0 || err != nil {
					return err
				}
			}*/

	default:
		fmt.Printf("gif: unknown extension 0x%.2x", extType)
	}

	// skip past all data
	file.Seek(baseOffset+size+1, os.SEEK_SET)

	res := Layout{
		Offset: baseOffset,
		Length: size + 1,
		Info:   "extension",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: baseOffset, Length: 1, Info: "block id (extension)", Type: Uint8},
			Layout{Offset: baseOffset + 1, Length: 1, Info: typeInfo, Type: Uint8},
		},
	}

	res.Childs = append(res.Childs, typeSpecific...)

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

func gifImageData(file *os.File, baseOffset int64) (*Layout, error) {

	// XXX need to decode first bytes of lzw stream to decode stream length

	file.Seek(baseOffset+1, os.SEEK_SET)

	length := int64(1)

	childs := []Layout{}
	childs = append(childs, Layout{Offset: baseOffset, Length: 1, Info: "lzw code size", Type: Uint8})

	for {
		var follows byte // number of bytes follows
		if err := binary.Read(file, binary.LittleEndian, &follows); err != nil {
			return nil, err
		}

		childs = append(childs, Layout{Offset: baseOffset + length, Length: 1, Info: "lzw block length", Type: Uint8})
		length += 1

		if follows == 0 {
			break
		}

		childs = append(childs, Layout{Offset: baseOffset + length, Length: int64(follows), Info: "lzw block", Type: Uint8})

		length += int64(follows)
		file.Seek(baseOffset+length, os.SEEK_SET)
	}

	res := Layout{
		Offset: baseOffset,
		Length: length,
		Info:   "image data",
		Type:   Group,
		Childs: childs,
	}

	return &res, nil
}

func gifTrailer(file *os.File, baseOffset int64) Layout {

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
