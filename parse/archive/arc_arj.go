package archive

// STATUS: 20%

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/martinlindhe/formats/parse"
)

const (
	arjBlockSizeMin  = 30
	arjBlockSizeMax  = 2600
	arjMaxSFX        = 500000 // size of self-extracting prefix
	arjHeaderIDHi    = 0xea
	arjHeaderIDLo    = 0x60
	arjFirstHdrSize  = 0x1e
	arjCommentMax    = 2048
	arjFileNameMax   = 512
	arjHeaderSizeMax = (arjFirstHdrSize + 10 + arjFileNameMax + arjCommentMax)
	arjCrcMask       = 0xffffffff
)

// ARJ parses the arj format
func ARJ(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isARJ(c.Header) {
		return nil, nil
	}

	mainHeader, err := parseARJMainHeader(c.File)

	// XXX rest of arj

	c.ParsedLayout.FileKind = parse.Archive
	c.ParsedLayout.MimeType = "application/x-arj"
	c.ParsedLayout.Layout = mainHeader

	return &c.ParsedLayout, err
}

func isARJ(b []byte) bool {

	return b[0] == 0x60 && b[1] == 0xea
}

// finds arj header and leaves file position at it
func findARJHeader(file *os.File) (int64, error) {

	reader := io.Reader(file)

	file.Seek(0, os.SEEK_SET)

	pos, _ := file.Seek(0, os.SEEK_CUR)
	lastpos, _ := file.Seek(0, os.SEEK_END)
	lastpos -= 2

	if lastpos > arjMaxSFX {
		lastpos = arjMaxSFX
	}
	// log.Println("starting", pos, lastpos)
	for ; pos < lastpos; pos++ {
		// log.Printf("setting pos to %04x\n", pos)
		pos2, _ := file.Seek(pos, os.SEEK_SET)
		if pos != pos2 {
			fmt.Printf("warning: expected %d, got %d\n", pos, pos2)
		}

		var c byte
		if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
			return 0, err
		}

		for pos < lastpos {
			if c != arjHeaderIDLo { // low order first
				if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
					return 0, err
				}
			} else {
				if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
					return 0, err
				}
				if c == arjHeaderIDHi {
					// log.Println("yes 1")
					break
				}
			}
			pos++
		}
		if pos >= lastpos {
			// log.Println("yes 2")
			break
		}

		var headerSize uint16
		if err := binary.Read(reader, binary.LittleEndian, &headerSize); err != nil {
			// log.Println("read err", err)
			return 0, err
		}

		log.Printf("header size %02x\n", headerSize)

		if headerSize <= arjHeaderSizeMax {
			log.Println("returning pos", pos)
			return pos, nil
		}
	}

	return 0, fmt.Errorf("could not find arj header in %s", file.Name())
}

var (
	arjHostOS = map[byte]string{
		0:  "MSDOS",
		1:  "PRIMOS",
		2:  "UNIX",
		3:  "AMIGA",
		4:  "MAC-OS",
		5:  "OS/2",
		6:  "APPLE GS",
		7:  "ATARI ST",
		8:  "NEXT",
		9:  "VAX VMS",
		10: "WIN95",
		11: "WIN32",
	}
)

func parseARJMainHeader(f *os.File) ([]parse.Layout, error) {

	pos, err := findARJHeader(f)
	if err != nil {
		return nil, err
	}

	hostOSName, _ := parse.ReadToMap(f, parse.Uint8, pos+7, arjHostOS)

	mainHeaderLen := int64(34)

	chunk := parse.Layout{
		Offset: pos,
		Type:   parse.Group,
		Info:   "main header",
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Type: parse.Uint16le, Info: "magic"},
			{Offset: pos + 2, Length: 2, Type: parse.Uint16le, Info: "basic header size"}, // excl. Magic+HdrSize
			{Offset: pos + 4, Length: 1, Type: parse.Uint8, Info: "size up to and including 'extra data'"},
			{Offset: pos + 5, Length: 1, Type: parse.Uint8, Info: "archiver version number"},
			{Offset: pos + 6, Length: 1, Type: parse.Uint8, Info: "minimum archiver version to extract"},
			{Offset: pos + 7, Length: 1, Type: parse.Uint8, Info: "host OS = " + hostOSName},
			{Offset: pos + 8, Length: 1, Type: parse.Uint8, Info: "arj flags"}, // XXX show bitfield
			{Offset: pos + 9, Length: 1, Type: parse.Uint8, Info: "security version"},
			{Offset: pos + 10, Length: 1, Type: parse.Uint8, Info: "file type"},
			{Offset: pos + 11, Length: 1, Type: parse.Uint8, Info: "reserved"},
			{Offset: pos + 12, Length: 4, Type: parse.ArjDateTime, Info: "created time"},
			{Offset: pos + 16, Length: 4, Type: parse.ArjDateTime, Info: "modified time"},
			{Offset: pos + 20, Length: 4, Type: parse.Uint32le, Info: "archive size for secured archive"},
			{Offset: pos + 24, Length: 4, Type: parse.Uint32le, Info: "security envelope file position"},
			{Offset: pos + 28, Length: 2, Type: parse.Uint16le, Info: "filespec position in filename"},
			{Offset: pos + 30, Length: 2, Type: parse.Uint16le, Info: "length in bytes of security envelope data"},
			{Offset: pos + 32, Length: 1, Type: parse.Uint8, Info: "encryption version"},
			{Offset: pos + 33, Length: 1, Type: parse.Uint8, Info: "last chapter"}, // XXX
		},
	}

	withExtData, _ := parse.ReadUint8(f, pos+4)
	if withExtData == 0x22 {
		chunk.Childs = append(chunk.Childs, []parse.Layout{
			{Offset: pos + 34, Length: 1, Type: parse.Uint8, Info: "protection factor"},
			{Offset: pos + 35, Length: 1, Type: parse.Uint8, Info: "flags (second series)"},
			{Offset: pos + 36, Length: 2, Type: parse.Uint8, Info: "spare bytes"},
		}...)
		mainHeaderLen += 4
	} else if withExtData == 0x1E {
		// no ext data
	} else {
		log.Fatalf("sample please. ext data = %02x", withExtData)
	}

	_, archiveNameLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos+mainHeaderLen, 255)
	if err != nil {
		return nil, err
	}

	_, commentLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos+mainHeaderLen+archiveNameLen, 4096)
	if err != nil {
		return nil, err
	}

	chunk.Length = mainHeaderLen + archiveNameLen + commentLen + 6

	pos += mainHeaderLen

	chunk.Childs = append(chunk.Childs, []parse.Layout{
		{Offset: pos, Length: archiveNameLen, Type: parse.ASCIIZ, Info: "archive name"},
	}...)
	pos += archiveNameLen

	chunk.Childs = append(chunk.Childs, []parse.Layout{
		{Offset: pos, Length: commentLen, Type: parse.ASCIIZ, Info: "comment"},
	}...)
	pos += commentLen

	chunk.Childs = append(chunk.Childs, []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "crc32"},
		{Offset: pos + 4, Length: 2, Type: parse.Uint16le, Info: "ext header size"},
	}...)
	pos += 6

	// XXX if ext header size > 0, it should follow here! need sample

	return []parse.Layout{chunk}, nil
}
