package parse

// STATUS borked

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	arjHeaderSize    = 0x22 // XXX what is size?1
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

func ARJ(file *os.File) (*ParsedLayout, error) {

	if !isARJ(file) {
		return nil, nil
	}

	mainHeader, err := parseARJMainHeader(file)

	// XXX rest of arj

	return &ParsedLayout{
		FileKind: Archive,
		Layout:   mainHeader}, err
}

func parseARJMainHeader(f *os.File) ([]Layout, error) {

	var err error
	offset, err := findARJHeader(f)
	if err != nil {
		return nil, err
	}

	mainHeaderLen := int64(35) // XXX hdr len?!

	f.Seek(mainHeaderLen, os.SEEK_SET)

	archiveName := ""
	comment := ""

	if archiveName, _, err = zeroTerminatedASCII(f); err != nil {
		return nil, err
	}
	archiveNameLen := int64(len(archiveName)) + 1 // including terminating zero

	if comment, _, err = zeroTerminatedASCII(f); err != nil {
		return nil, err
	}
	commentLen := int64(len(comment)) + 1

	chunk := Layout{
		Offset: offset,
		Length: mainHeaderLen + int64(len(archiveName)+len(comment)) + 8,
		Type:   Group,
		Info:   "main header",
		Childs: []Layout{
			// XXX convert arjMainHeader into []Layout and add to Childs in return
			{Offset: offset, Length: 2, Type: Uint16le, Info: "magic"},
			{Offset: offset + 2, Length: 2, Type: Uint16le, Info: "basic header size"}, // excl. Magic+HdrSize
			{Offset: offset + 4, Length: 1, Type: Uint8, Info: "size up to and including 'extra data'"},
			{Offset: offset + 5, Length: 1, Type: Uint8, Info: "archiver version number"},
			{Offset: offset + 6, Length: 1, Type: Uint8, Info: "minimum archiver version to extract"},
			{Offset: offset + 7, Length: 1, Type: Uint8, Info: "host OS"},   // XXX map hostOSes
			{Offset: offset + 8, Length: 1, Type: Uint8, Info: "arj flags"}, // XXX show bitfield
			{Offset: offset + 9, Length: 1, Type: Uint8, Info: "security version"},
			{Offset: offset + 10, Length: 1, Type: Uint8, Info: "file type"},        // XXX map fileTypes
			{Offset: offset + 11, Length: 4, Type: Uint32le, Info: "created time"},  // XXX time in "msdos-format"
			{Offset: offset + 15, Length: 4, Type: Uint32le, Info: "modified time"}, // XXX time in "msdos-format"
			{Offset: offset + 19, Length: 4, Type: Uint32le, Info: "archive size for secured archive"},
			{Offset: offset + 23, Length: 4, Type: Uint32le, Info: "security envelope file position"},
			{Offset: offset + 27, Length: 4, Type: Uint32le, Info: "filespec position in filename"},
			{Offset: offset + 31, Length: 2, Type: Uint16le, Info: "length in bytes of security envelope data"},
			{Offset: offset + 33, Length: 1, Type: Uint8, Info: "encryption version"},
			{Offset: offset + 34, Length: 1, Type: Uint8, Info: "last chapter"}, // XXX
		},
	}
	offset += mainHeaderLen

	chunk.Childs = append(chunk.Childs, []Layout{
		{Offset: offset, Length: archiveNameLen, Type: ASCIIZ, Info: "archive name"},
	}...)
	offset += archiveNameLen

	chunk.Childs = append(chunk.Childs, []Layout{
		{Offset: offset, Length: commentLen, Type: ASCIIZ, Info: "comment"},
	}...)
	offset += commentLen

	chunk.Childs = append(chunk.Childs, []Layout{
		{Offset: offset, Length: 4, Type: Uint32le, Info: "crc32"},
		{Offset: offset + 4, Length: 4, Type: Uint32le, Info: "ext header size"},
	}...)
	offset += 8

	return []Layout{chunk}, nil

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
		fmt.Printf("setting pos to %04x\n", pos)
		pos2, _ := file.Seek(pos, os.SEEK_SET)
		if pos != pos2 {
			fmt.Printf("expected %d, got %d\n", pos, pos2)
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
			// fmt.Printf("arcpos %04x\n", arcpos)

			// XXX implement crc check?
			//crc = crcMask
			//fread_crc(header, headersize, fd)
			//if (crc ^ crcMask) == fget_crc(fd) {
			file.Seek(pos, os.SEEK_SET)
			return pos, nil
			//}
		}
	}

	return 0, fmt.Errorf("could not find arj header")
}
