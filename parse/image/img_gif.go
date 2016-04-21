package image

// STATUS: 90%

// XXX problems parsing samples/images/gif/gif_89a_005_with_application.gif

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/martinlindhe/formats/parse"
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
	gifExtensions = map[byte]string{
		1:    "text",
		0xf9: "graphic control",
		0xfe: "comment",
		0xff: "application",
	}
)

// section indicators
const (
	sExtension       = 0x21
	sImageDescriptor = 0x2c
	sTrailer         = 0x3b
)

// extensions
const (
	eText           = 0x01
	eGraphicControl = 0xf9
	eComment        = 0xfe
	eApplication    = 0xff
)

// misc
const (
	imgDescriptorLen = 10
)

func GIF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isGIF(&c.Header) {
		return nil, nil
	}
	return parseGIF(c)
}

func isGIF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'G' || b[1] != 'I' || b[2] != 'F' || b[3] != '8' {
		return false
	}
	if b[4] != '7' && b[4] != '9' {
		return false
	}
	return true
}

func parseGIF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl := c.ParsedLayout
	pl.FileKind = parse.Image

	header := gifHeader(c.File)
	pl.Layout = append(pl.Layout, header)
	pos += header.Length

	logicalDesc := gifLogicalDescriptor(c.File)
	pl.Layout = append(pl.Layout, logicalDesc)
	pos += logicalDesc.Length

	globalColorTableFlag := pl.DecodeBitfieldFromInfo(c.File, "global color table flag")

	if globalColorTableFlag == 1 {
		sizeOfGCT := pl.DecodeBitfieldFromInfo(c.File, "global color table size")

		if gctByteLen, ok := gctToLengthMap[byte(sizeOfGCT)]; ok {
			pl.Layout = append(pl.Layout, gifGlobalColorTable(c.File, gctByteLen))
			pos += gctByteLen
		}
	}

	for {
		_, err := c.File.Seek(pos, os.SEEK_SET)
		if err != nil {
			fmt.Println("seek err", err)
		}

		var b byte
		if err := binary.Read(c.File, binary.LittleEndian, &b); err != nil {
			if err == io.EOF {
				fmt.Println("warning: did not find gif trailer")
				return &pl, nil
			}
			return nil, err
		}

		fmt.Printf("ext %02x at %04x\n", b, pos)
		switch b {
		case sExtension:
			gfxExt, err := gifExtension(c.File, pos)
			if err != nil {
				return nil, err
			}
			pl.Layout = append(pl.Layout, *gfxExt)
			pos += gfxExt.Length

		case sImageDescriptor:
			imgDescriptor := gifImageDescriptor(c.File, pos)
			if imgDescriptor != nil {
				pl.Layout = append(pl.Layout, *imgDescriptor)
				pos += imgDescriptor.Length
			}
			if pl.DecodeBitfieldFromInfo(c.File, "local color table flag") > 0 {
				sizeOfLCT := pl.DecodeBitfieldFromInfo(c.File, "local color table size")
				if lctByteLen, ok := gctToLengthMap[byte(sizeOfLCT)]; ok {
					localTbl := gifLocalColorTable(c.File, pos, lctByteLen)
					pl.Layout = append(pl.Layout, localTbl)
					pos += localTbl.Length
				}
			}

			imgData, err := gifImageData(c.File, pos)
			if err != nil {
				return nil, err
			}
			pl.Layout = append(pl.Layout, *imgData)
			pos += imgData.Length

		case sTrailer:
			pl.Layout = append(pl.Layout, gifTrailer(c.File, pos))
			return &pl, nil

		default:
			fmt.Printf("warning: unknown gif section id %02x at %04x\n", b, pos)
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

func gifLogicalDescriptor(file *os.File) parse.Layout {

	pos := int64(6)
	return parse.Layout{
		Offset: pos,
		Length: 7,
		Info:   "logical screen descriptor",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "width", Type: parse.Uint16le},
			{Offset: pos + 2, Length: 2, Info: "height", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 1, Info: "packed", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 3, Info: "global color table size"},
				{Low: 3, Length: 1, Info: "sort flag"},
				{Low: 4, Length: 3, Info: "color resolution"},
				{Low: 7, Length: 1, Info: "global color table flag"},
			}},
			{Offset: pos + 5, Length: 1, Info: "background color", Type: parse.Uint8},
			{Offset: pos + 6, Length: 1, Info: "aspect ratio", Type: parse.Uint8},
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
				{Low: 0, Length: 2, Info: "local color table size"},
				{Low: 3, Length: 2, Info: "reserved"},
				{Low: 5, Length: 1, Info: "sort flag"},
				{Low: 6, Length: 1, Info: "interlace flag"},
				{Low: 7, Length: 1, Info: "local color table flag"},
			}}}}
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
		id := fmt.Sprintf("%d", cnt)
		childs = append(childs, parse.Layout{
			Offset: pos + i,
			Length: 3,
			Info:   "color " + id,
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

	extType, _ := parse.ReadUint8(file, pos+1)
	typeInfo, _ := parse.ReadToMap(file, parse.Uint8, pos+1, gifExtensions)
	typeSpecific := []parse.Layout{}
	size := int64(0)
	res := parse.Layout{
		Offset: pos,
		Length: size + 1,
		Info:   typeInfo + " extension",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "block id (extension)", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: "type = " + typeInfo, Type: parse.Uint8},
		}}
	pos += 2

	switch extType {
	case eText:
		size = 13

	case eGraphicControl:
		size = 7
		typeSpecific = []parse.Layout{
			{Offset: pos, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: "packed #2", Type: parse.Uint8},
			{Offset: pos + 2, Length: 2, Info: "delay time", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 1, Info: "transparent color index", Type: parse.Uint8},
			{Offset: pos + 5, Length: 1, Info: "block terminator", Type: parse.Uint8},
		}

	case eComment:
		// nothing to do but read the data.
		lenByte, _ := parse.ReadUint8(file, pos)

		size = 2 + int64(lenByte) + 1 // including terminating 0

		typeSpecific = []parse.Layout{
			{Offset: pos, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 1, Length: size - 2, Info: "data", Type: parse.ASCIIZ},
		}

		/*

		   			struct APPLICATIONEXTENTION {
		                   UBYTE ExtensionIntroducer; // 0x21
		                   UBYTE ApplicationLabel; // 0xFF

		                   struct APPLICATIONSUBBLOCK {
		                       UBYTE   BlockSize;
		                       char    ApplicationIdentifier[8];
		                       char    ApplicationAuthenticationCode[3];
		                   } ApplicationSubBlock;
		                   DATASUBBLOCKS ApplicationData;
		               } ApplicationExtension;
		*/
	case eApplication:
		size = 14
		typeSpecific = []parse.Layout{
			{Offset: pos, Length: 1, Info: "byte size", Type: parse.Uint8},
			{Offset: pos + 1, Length: 8, Info: "application id", Type: parse.ASCII},
			{Offset: pos + 9, Length: 3, Info: "application authentication code", Type: parse.ASCII},
		}
		extData, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, pos+3, 8)
		authCode, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, pos+11, 3)

		if extData == "NETSCAPE" && authCode == "2.0" {
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
		fmt.Printf("gif: unknown extension 0x%.2x\n", extType)
	}

	res.Length += size
	res.Childs = append(res.Childs, typeSpecific...)

	return &res, nil
}

func gifImageData(file *os.File, pos int64) (*parse.Layout, error) {

	length := int64(1)
	childs := []parse.Layout{{
		Offset: pos,
		Length: 1,
		Info:   "lzw code size",
		Type:   parse.Uint8}}

	// decodes first bytes of lzw stream to calculate stream length

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
		Childs: childs}

	return &res, nil
}

func gifSubBlocks(file *os.File, pos int64) ([]parse.Layout, error) {

	length := int64(0)
	childs := []parse.Layout{}

	for {
		file.Seek(pos, os.SEEK_SET)

		var follows byte // number of bytes follows
		if err := binary.Read(file, binary.LittleEndian, &follows); err != nil {
			if err == io.EOF {
				fmt.Println("XXX sub blocks unexpected EOF")
				break
			}
			return nil, err
		}

		childs = append(childs, parse.Layout{
			Offset: pos,
			Length: 1,
			Info:   "block length",
			Type:   parse.Uint8})
		length += 1
		pos += 1

		// XXX special case 0xff ?

		if follows == 0 {
			break
		}

		childs = append(childs, parse.Layout{
			Offset: pos,
			Length: int64(follows),
			Info:   "block",
			Type:   parse.Bytes})
		length += int64(follows)
		pos += int64(follows)
	}
	return childs, nil
}

func gifTrailer(file *os.File, pos int64) parse.Layout {

	return parse.Layout{
		Offset: pos,
		Length: 1,
		Info:   "trailer",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "trailer", Type: parse.Uint8},
		}}
}
