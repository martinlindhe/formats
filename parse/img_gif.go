package parse

// STATUS ~80%: gif89 most files ok, gif87 broken!

// XXX samples/gif/gif_89a_002_anim.gif  lzw block decode seems broken, start offset wrong?
// XXX samples/gif/gif_87a_001.gif is broken!

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

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

	// XXX 1. make test using a specific file, with known PACKED value, and use that to test the decode stuff!

	// XXX hack... decodeBitfieldFromInfo should return 1 but returns 2 now for soem reason?!
	if res.DecodeBitfieldFromInfo(file, "global color table flag") != 0 {
		if res.DecodeBitfieldFromInfo(file, "global color table flag") != 1 {
			fmt.Println(res.DecodeBitfieldFromInfo(file, "global color table flag"))
			panic("res is odd!")
		}
		sizeOfGCT := res.DecodeBitfieldFromInfo(file, "size of global color table")
		if gctByteLen, ok := gctToLengthMap[byte(sizeOfGCT)]; ok {
			res.Layout = append(res.Layout, gifGlobalColorTable(file, gctByteLen))
		}
	}

	for {

		offset, _ := file.Seek(0, os.SEEK_CUR)

		var b byte
		if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
			if err == io.EOF {
				fmt.Println("XXX did not find gif trailer!")
				return &res, nil
			}
			return nil, err
		}

		switch b {
		case sExtension:
			gfxExt, err := gifExtension(file, offset)
			if err != nil {
				return nil, err
			}
			res.Layout = append(res.Layout, *gfxExt)

		case sImageDescriptor:
			imgDescriptor := gifImageDescriptor(file, offset)
			if imgDescriptor != nil {
				res.Layout = append(res.Layout, *imgDescriptor)
			}
			if res.DecodeBitfieldFromInfo(file, "local color table flag") == 1 {
				// XXX this is untested due to lack of sample with a local color table
				sizeOfLCT := res.DecodeBitfieldFromInfo(file, "size of local color table")
				if lctByteLen, ok := gctToLengthMap[byte(sizeOfLCT)]; ok {
					localTbl := gifLocalColorTable(file, offset+imgDescriptorLen, lctByteLen)
					res.Layout = append(res.Layout, localTbl)
				}
			}

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

			Layout{Offset: baseOffset + 9, Length: 1, Info: "packed #3", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 2, Info: "size of local color table"},
				Mask{Low: 3, Length: 2, Info: "reserved"},
				Mask{Low: 5, Length: 1, Info: "sort flag"},
				Mask{Low: 6, Length: 1, Info: "interlace flag"},
				Mask{Low: 7, Length: 1, Info: "local color table flag"},
			}},
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

func gifLocalColorTable(file *os.File, baseOffset int64, byteLen int64) Layout {

	childs := []Layout{}

	cnt := 0
	for i := int64(0); i < byteLen; i += 3 {
		cnt++
		childs = append(childs, Layout{Offset: baseOffset + i, Length: 3, Info: fmt.Sprintf("color %d", cnt), Type: RGB})
	}

	return Layout{
		Offset: baseOffset,
		Length: byteLen,
		Info:   "local color table",
		Type:   Group,
		Childs: childs,
	}
}

func gifExtension(file *os.File, baseOffset int64) (*Layout, error) {

	var extType byte
	if err := binary.Read(file, binary.LittleEndian, &extType); err != nil {
		return nil, err
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

		var lenByte byte
		if err := binary.Read(file, binary.LittleEndian, &lenByte); err != nil {
			return nil, err
		}

		size = 2 + int64(lenByte) + 1 // including terminating 0

		typeSpecific = []Layout{
			Layout{Offset: baseOffset + 2, Length: 1, Info: "byte size", Type: Uint8},
			Layout{Offset: baseOffset + 3, Length: size - 2, Info: "data", Type: ASCIIZ},
		}

	case eApplication:
		typeInfo = "application"
		var lenByte byte
		if err := binary.Read(file, binary.LittleEndian, &lenByte); err != nil {
			return nil, err
		}

		size = 2 + int64(lenByte)

		typeSpecific = []Layout{
			Layout{Offset: baseOffset + 2, Length: 1, Info: "byte size", Type: Uint8},
			Layout{Offset: baseOffset + 3, Length: size - 2, Info: "data", Type: Uint8},
		}

		extData := readBytesFrom(file, baseOffset+3, size-2)

		if string(extData) == "NETSCAPE2.0" {
			// animated gif extension
			subBlocks, err := gifSubBlocks(file, baseOffset+3+11)
			if err != nil {
				return nil, err
			}
			typeSpecific = append(typeSpecific, subBlocks...)
			for _, b := range subBlocks {
				size += b.Length
			}
		}

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

	return &res, nil
}

func gifReadBlock(file *os.File) (int, error) {

	var b byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return 0, err
	}

	// return io.ReadFull(file, d.tmp[:n])
	return 0, nil
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
			Layout{Offset: base + 4, Length: 1, Info: "packed", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 3, Info: "size of global color table"},
				Mask{Low: 3, Length: 1, Info: "color table sort flag"},
				Mask{Low: 4, Length: 3, Info: "color resolution"},
				Mask{Low: 7, Length: 1, Info: "global color table flag"},
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

	lzwSubBlocks, err := gifSubBlocks(file, baseOffset+1)
	if err != nil {
		return nil, err
	}
	childs = append(childs, lzwSubBlocks...)
	for _, b := range lzwSubBlocks {
		length += b.Length
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

func gifSubBlocks(file *os.File, baseOffset int64) ([]Layout, error) {

	length := int64(0)
	childs := []Layout{}
	file.Seek(baseOffset, os.SEEK_SET)

	for {
		var follows byte // number of bytes follows
		if err := binary.Read(file, binary.LittleEndian, &follows); err != nil {
			if err == io.EOF {
				fmt.Println("XXX sub blocks unexpected EOF")
				break
			}
			return nil, err
		}

		childs = append(childs, Layout{Offset: baseOffset + length, Length: 1, Info: "block length", Type: Uint8})
		length += 1

		if follows == 0 {
			break
		}

		childs = append(childs, Layout{Offset: baseOffset + length, Length: int64(follows), Info: "block", Type: Uint8})

		length += int64(follows)
		file.Seek(baseOffset+length, os.SEEK_SET)
	}
	return childs, nil
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
