package image

// STATUS ~80%: gif89 most files ok, gif87 broken!

// XXX samples/gif/gif_89a_002_anim.gif  lzw block decode seems broken, start offset wrong?
// XXX samples/gif/gif_87a_001.gif is broken!

import (
	"encoding/binary"
	"fmt"
	"github.com/martinlindhe/formats/parse"
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

func GIF(file *os.File) (*parse.ParsedLayout, error) {

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

func parseGIF(file *os.File) (*parse.ParsedLayout, error) {

	res := parse.ParsedLayout{
		FileKind: parse.Image,
	}

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

func gifHeader(file *os.File) parse.Layout {

	pos := int64(0)
	return parse.Layout{
		Offset: pos,
		Length: 6,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "signature", Type: parse.ASCII},
			{Offset: pos + 3, Length: 3, Info: "version", Type: parse.ASCII},
		},
	}
}

func gifImageDescriptor(file *os.File, pos int64) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Length: imgDescriptorLen,
		Info:   "image descriptor",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "image separator", Type: parse.Uint8},
			{Offset: pos + 1, Length: 2, Info: "image left", Type: parse.Uint16le},
			{Offset: pos + 3, Length: 2, Info: "image top", Type: parse.Uint16le},
			{Offset: pos + 5, Length: 2, Info: "image width", Type: parse.Uint16le},
			{Offset: pos + 7, Length: 2, Info: "image height", Type: parse.Uint16le},
			{Offset: pos + 9, Length: 1, Info: "packed #3", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 2, Info: "size of local color table"},
				{Low: 3, Length: 2, Info: "reserved"},
				{Low: 5, Length: 1, Info: "sort flag"},
				{Low: 6, Length: 1, Info: "interlace flag"},
				{Low: 7, Length: 1, Info: "local color table flag"},
			}},
		},
	}
	return &res
}

func gifGlobalColorTable(file *os.File, byteLen int64) parse.Layout {

	pos := int64(0x0d)

	childs := []parse.Layout{}

	cnt := 0
	for i := int64(0); i < byteLen; i += 3 {
		cnt++
		childs = append(childs, parse.Layout{
			Offset: pos + i,
			Length: 3,
			Info:   fmt.Sprintf("color %d", cnt),
			Type:   parse.RGB})
	}

	return parse.Layout{
		Offset: pos,
		Length: byteLen,
		Info:   "global color table",
		Type:   parse.Group,
		Childs: childs}
}

func gifLocalColorTable(file *os.File, pos int64, byteLen int64) parse.Layout {

	childs := []parse.Layout{}

	cnt := 0
	for i := int64(0); i < byteLen; i += 3 {
		cnt++
		childs = append(childs, parse.Layout{
			Offset: pos + i,
			Length: 3,
			Info:   fmt.Sprintf("color %d", cnt),
			Type:   parse.RGB})
	}

	return parse.Layout{
		Offset: pos,
		Length: byteLen,
		Info:   "local color table",
		Type:   parse.Group,
		Childs: childs}
}

func gifExtension(file *os.File, pos int64) (*parse.Layout, error) {

	var extType byte
	if err := binary.Read(file, binary.LittleEndian, &extType); err != nil {
		return nil, err
	}

	typeSpecific := []parse.Layout{}
	typeInfo := ""

	size := int64(0)
	switch extType {
	case eText:
		size = 13
		typeInfo = "text"

	case eGraphicControl:
		size = 7
		typeInfo = "graphic control"
		typeSpecific = []parse.Layout{
			{Offset: pos + 2, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 3, Length: 1, Info: "packed #2", Type: parse.Uint8},
			{Offset: pos + 4, Length: 2, Info: "delay time", Type: parse.Uint16le},
			{Offset: pos + 6, Length: 1, Info: "transparent color index", Type: parse.Uint8},
			{Offset: pos + 7, Length: 1, Info: "block terminator", Type: parse.Uint8},
		}

	case eComment:
		// nothing to do but read the data.
		typeInfo = "comment"

		var lenByte byte
		if err := binary.Read(file, binary.LittleEndian, &lenByte); err != nil {
			return nil, err
		}

		size = 2 + int64(lenByte) + 1 // including terminating 0

		typeSpecific = []parse.Layout{
			{Offset: pos + 2, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 3, Length: size - 2, Info: "data", Type: parse.ASCIIZ},
		}

	case eApplication:
		typeInfo = "application"
		var lenByte byte
		if err := binary.Read(file, binary.LittleEndian, &lenByte); err != nil {
			return nil, err
		}

		size = 2 + int64(lenByte)

		typeSpecific = []parse.Layout{
			{Offset: pos + 2, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 3, Length: size - 2, Info: "data", Type: parse.Uint8},
		}

		extData := parse.ReadBytesFrom(file, pos+3, size-2)

		if string(extData) == "NETSCAPE2.0" {
			// animated gif extension
			subBlocks, err := gifSubBlocks(file, pos+3+11)
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
	file.Seek(pos+size+1, os.SEEK_SET)

	res := parse.Layout{
		Offset: pos,
		Length: size + 1,
		Info:   "extension",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "block id (extension)", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: typeInfo, Type: parse.Uint8},
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

func gifLogicalDescriptor(file *os.File) parse.Layout {

	pos := int64(0x06)
	return parse.Layout{
		Offset: pos,
		Length: 7,
		Info:   "logical screen descriptor",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "width", Type: parse.Uint16le},
			{Offset: pos + 2, Length: 2, Info: "height", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 1, Info: "packed", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 3, Info: "size of global color table"},
				{Low: 3, Length: 1, Info: "color table sort flag"},
				{Low: 4, Length: 3, Info: "color resolution"},
				{Low: 7, Length: 1, Info: "global color table flag"},
			}},
			{Offset: pos + 5, Length: 1, Info: "background color", Type: parse.Uint8},
			{Offset: pos + 6, Length: 1, Info: "aspect ratio", Type: parse.Uint8},
		},
	}
}

func gifImageData(file *os.File, pos int64) (*parse.Layout, error) {

	// XXX need to decode first bytes of lzw stream to decode stream length

	file.Seek(pos+1, os.SEEK_SET)

	length := int64(1)
	childs := []parse.Layout{{
		Offset: pos,
		Length: 1,
		Info:   "lzw code size",
		Type:   parse.Uint8}}

	lzwSubBlocks, err := gifSubBlocks(file, pos+1)
	if err != nil {
		return nil, err
	}
	childs = append(childs, lzwSubBlocks...)
	for _, b := range lzwSubBlocks {
		length += b.Length
	}

	res := parse.Layout{
		Offset: pos,
		Length: length,
		Info:   "image data",
		Type:   parse.Group,
		Childs: childs,
	}

	return &res, nil
}

func gifSubBlocks(file *os.File, pos int64) ([]parse.Layout, error) {

	length := int64(0)
	childs := []parse.Layout{}
	file.Seek(pos, os.SEEK_SET)

	for {
		var follows byte // number of bytes follows
		if err := binary.Read(file, binary.LittleEndian, &follows); err != nil {
			if err == io.EOF {
				fmt.Println("XXX sub blocks unexpected EOF")
				break
			}
			return nil, err
		}

		childs = append(childs, parse.Layout{
			Offset: pos + length,
			Length: 1,
			Info:   "block length",
			Type:   parse.Uint8})
		length += 1

		if follows == 0 {
			break
		}

		childs = append(childs, parse.Layout{
			Offset: pos + length,
			Length: int64(follows),
			Info:   "block",
			Type:   parse.Uint8})

		length += int64(follows)
		file.Seek(pos+length, os.SEEK_SET)
	}
	return childs, nil
}

func gifTrailer(file *os.File, pos int64) parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Length: 1,
		Info:   "trailer",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "trailer", Type: parse.Uint8},
		},
	}
	return res
}
