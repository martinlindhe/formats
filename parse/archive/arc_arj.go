package archive

// STATUS borked

import (
	"encoding/binary"
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"io"
	"os"
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

func ARJ(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isARJ(file) {
		return nil, nil
	}

	mainHeader, err := parseARJMainHeader(file)

	// XXX rest of arj

	pl.FileKind = parse.Archive
	pl.Layout = mainHeader

	return &pl, err
}

func parseARJMainHeader(f *os.File) ([]parse.Layout, error) {

	pos, err := findARJHeader(f)
	if err != nil {
		return nil, err
	}

	mainHeaderLen := int64(35) // XXX

	pos = mainHeaderLen
	archiveName, _, _ := parse.ReadZeroTerminatedASCIIUntil(f, pos, 32)
	archiveNameLen := int64(len(archiveName)) + 1 // including terminating zero
	pos += archiveNameLen

	comment, _, _ := parse.ReadZeroTerminatedASCIIUntil(f, pos, 32)
	commentLen := int64(len(comment)) + 1

	chunk := parse.Layout{
		Offset: pos,
		Length: mainHeaderLen + archiveNameLen + commentLen + 8,
		Type:   parse.Group,
		Info:   "main header",
		Childs: []parse.Layout{
			// XXX convert arjMainHeader into []Layout and add to Childs in return
			{Offset: pos, Length: 2, Type: parse.Uint16le, Info: "magic"},
			{Offset: pos + 2, Length: 2, Type: parse.Uint16le, Info: "basic header size"}, // excl. Magic+HdrSize
			{Offset: pos + 4, Length: 1, Type: parse.Uint8, Info: "size up to and including 'extra data'"},
			{Offset: pos + 5, Length: 1, Type: parse.Uint8, Info: "archiver version number"},
			{Offset: pos + 6, Length: 1, Type: parse.Uint8, Info: "minimum archiver version to extract"},
			{Offset: pos + 7, Length: 1, Type: parse.Uint8, Info: "host OS"},   // XXX map hostOSes
			{Offset: pos + 8, Length: 1, Type: parse.Uint8, Info: "arj flags"}, // XXX show bitfield
			{Offset: pos + 9, Length: 1, Type: parse.Uint8, Info: "security version"},
			{Offset: pos + 10, Length: 1, Type: parse.Uint8, Info: "file type"},        // XXX map fileTypes
			{Offset: pos + 11, Length: 4, Type: parse.Uint32le, Info: "created time"},  // XXX time in "msdos-format"
			{Offset: pos + 15, Length: 4, Type: parse.Uint32le, Info: "modified time"}, // XXX time in "msdos-format"
			{Offset: pos + 19, Length: 4, Type: parse.Uint32le, Info: "archive size for secured archive"},
			{Offset: pos + 23, Length: 4, Type: parse.Uint32le, Info: "security envelope file position"},
			{Offset: pos + 27, Length: 4, Type: parse.Uint32le, Info: "filespec position in filename"},
			{Offset: pos + 31, Length: 2, Type: parse.Uint16le, Info: "length in bytes of security envelope data"},
			{Offset: pos + 33, Length: 1, Type: parse.Uint8, Info: "encryption version"},
			{Offset: pos + 34, Length: 1, Type: parse.Uint8, Info: "last chapter"}, // XXX
		},
	}
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
		{Offset: pos + 4, Length: 4, Type: parse.Uint32le, Info: "ext header size"},
	}...)
	pos += 8

	return []parse.Layout{chunk}, nil

	/*
	   XXX dont understand to parse 0x22, is 0 in both my samples
	   ?   extra data
	     1   arj protection factor
	     1   arj flags (second series)
	               (0x01 = ALTVOLNAME_FLAG) indicates special volume naming
	                                        option
	               (0x02 = reserved bit)
	     2   spare bytes
	*/
}

func isARJ(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	r := io.Reader(file)
	var b [2]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}
	return b[0] == 0x60 && b[1] == 0xea
}

/**
 * finds arj header and leaves file position at it
 */
func findARJHeader(file *os.File) (int64, error) {

	reader := io.Reader(file)

	pos, _ := file.Seek(0, os.SEEK_CUR)
	lastpos, _ := file.Seek(0, os.SEEK_END)
	lastpos -= 2

	if lastpos > arjMaxSFX {
		lastpos = arjMaxSFX
	}
	for ; pos < lastpos; pos++ {
		// fmt.Printf("setting pos to %04x\n", pos)
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
					// fmt.Println("yes 1")
					break
				}
			}
			pos++
		}
		if pos >= lastpos {
			// fmt.Println("yes 2")
			break
		}

		var headerSize uint16
		if err := binary.Read(reader, binary.LittleEndian, &headerSize); err != nil {
			return 0, err
		}

		// fmt.Printf("header size %02x\n", headerSize)

		if headerSize <= arjHeaderSizeMax {
			return pos, nil
		}
	}

	return 0, fmt.Errorf("could not find arj header")
}
