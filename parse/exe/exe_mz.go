package exe

// https://en.wikipedia.org/wiki/.exe
// STATUS: 60%

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	mzHeaderLen  = int64(28) // XXX
	subHeaderLen = int64(36) // XXX
)

// MZ parses the mz executable format
func MZ(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMZ(c.Header) {
		return nil, nil
	}
	return parseMZ(c)
}

func isMZ(b []byte) bool {

	if b[0] != 'M' || b[1] != 'Z' {
		return false
	}
	return true
}

func parseMZ(c *parse.Checker) (*parse.ParsedLayout, error) {

	c.ParsedLayout.FileKind = parse.Executable
	c.ParsedLayout.MimeType = "application/x-dosexec"
	pos := int64(0)
	mz := parse.Layout{
		Offset: pos,
		Length: mzHeaderLen,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 2, Length: 2, Info: "extra bytes", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 2, Info: "pages", Type: parse.Uint16le},
			{Offset: pos + 6, Length: 2, Info: "relocation items", Type: parse.Uint16le},
			{Offset: pos + 8, Length: 2, Info: "header size in paragraphs", Type: parse.Uint16le}, // 1 paragraph = group of 16 bytes
			{Offset: pos + 10, Length: 2, Info: "min allocation", Type: parse.Uint16le},
			{Offset: pos + 12, Length: 2, Info: "max allocation", Type: parse.Uint16le},
			{Offset: pos + 14, Length: 2, Info: "initial ss", Type: parse.Uint16le},
			{Offset: pos + 16, Length: 2, Info: "initial sp", Type: parse.Uint16le},
			{Offset: pos + 18, Length: 2, Info: "checksum", Type: parse.Uint16le},
			{Offset: pos + 20, Length: 2, Info: "initial ip", Type: parse.Uint16le},
			{Offset: pos + 22, Length: 2, Info: "initial cs", Type: parse.Uint16le},
			{Offset: pos + 24, Length: 2, Info: "relocation offset", Type: parse.Uint16le},
			{Offset: pos + 26, Length: 2, Info: "overlay", Type: parse.Uint16le},
		}}

	c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, mz)

	custom := findCustomDOSHeaders(c.File, c.Header)
	if custom != nil {
		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, custom...)
	}

	hdrSizeInParagraphs, _ := parse.ReadUint16le(c.File, pos+8)
	ip, _ := parse.ReadUint16le(c.File, pos+20)
	cs, _ := parse.ReadUint16le(c.File, pos+22)
	relocOffset, _ := parse.ReadUint16le(c.File, pos+24)

	if relocOffset == 0x40 {
		// 0x40 for new-(NE,LE,LX,W3,PE etc.) executable
		pos += mzHeaderLen

		newHeaderPos, _ := parse.ReadUint32le(c.File, pos+32)
		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, parse.Layout{
			Offset: pos,
			Length: subHeaderLen,
			Info:   "sub header", // XXX name
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 8, Info: "reserved", Type: parse.Bytes},
				{Offset: pos + 8, Length: 2, Info: "oem id", Type: parse.Uint16le},
				{Offset: pos + 10, Length: 2, Info: "oem info", Type: parse.Uint16le},
				{Offset: pos + 12, Length: 20, Info: "reserved 2", Type: parse.Uint16le},
				{Offset: pos + 32, Length: 4, Info: "start of ext header", Type: parse.Uint32le},
			}})

		pos = int64(newHeaderPos)
		newHeaderID, _, _ := parse.ReadZeroTerminatedASCIIUntil(c.File, pos, 2)

		switch newHeaderID {
		case "LX":
			// OS/2 (32-bit)
			c.ParsedLayout.FormatName = "mz-lx"
			header, _ := parseMzLxHeader(c.File, pos)
			c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, header...)

		case "LE":
			// Win, OS/2 (mixed 16/32-bit)
			c.ParsedLayout.FormatName = "mz-le"
			header, _ := parseMzLeHeader(c.File, pos)
			c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, header...)

		case "NE":
			// Win16, OS/2
			c.ParsedLayout.FormatName = "mz-ne"
			header, _ := parseMzNeHeader(c.File, pos)
			c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, header...)

		case "PE":
			// Win32, Win64
			c.ParsedLayout.FormatName = "mz-pe"
			header, _ := parseMzPeHeader(c.File, pos)
			c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, header...)

		default:
			log.Println("error: unknown newHeaderID: " + newHeaderID)
		}

		exeStart := int64(((hdrSizeInParagraphs + cs) * 16) + ip)

		dosStubLen := int64(newHeaderPos) - exeStart
		pos = exeStart
		dosStub := parse.Layout{
			Offset: pos,
			Length: dosStubLen,
			Info:   "dos stub",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: dosStubLen, Info: "dos stub", Type: parse.Bytes},
			}}

		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, dosStub)

	} else {
		cs, _ := parse.ReadUint16le(c.File, pos+22)
		relocItems, _ := parse.ReadUint16le(c.File, pos+6)
		if relocItems > 0 {
			pos = int64(relocOffset)
			reloc := parse.Layout{
				Offset: pos,
				Length: int64(relocItems) * 4,
				Info:   "relocation table",
				Type:   parse.Group}

			for i := 1; i <= int(relocItems); i++ {
				id := fmt.Sprintf("%d", i)

				c.File.Seek(pos, os.SEEK_SET)

				var b [2]uint16
				if err := binary.Read(c.File, binary.LittleEndian, &b); err != nil && err != io.EOF {
					return nil, err
				}
				// XXX abs offset seems wrong
				// log.Println("x = ", hdrSizeInParagraphs+cs-b[1])
				abs := (hdrSizeInParagraphs+cs-b[1])*16 + b[0]

				reloc.Childs = append(reloc.Childs, []parse.Layout{
					{Offset: pos, Length: 4, Info: "offset:segment " + id, Type: parse.DOSOffsetSegment},
				}...)

				relocLen := int64(4) // XXX section length ?
				c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, parse.Layout{
					Offset: int64(abs),
					Length: relocLen,
					Info:   "relocation " + id,
					Type:   parse.Group,
					Childs: []parse.Layout{
						{Offset: int64(abs), Length: relocLen, Info: "relocation " + id, Type: parse.Bytes},
					}})
				pos += 4
			}
			c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, reloc)
		}

		// XXX point to dos entry point
		exeStart := int64(((hdrSizeInParagraphs + cs) * 16) + ip)

		pos = exeStart
		dosEntry := parse.Layout{
			Offset: pos,
			Length: 4, // XXX
			Info:   "program",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "entry point", Type: parse.Bytes},
			}}

		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, dosEntry)
	}

	c.ParsedLayout.Sort()

	return &c.ParsedLayout, nil
}
