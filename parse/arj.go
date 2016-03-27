package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type arjMainHeader struct {
	Magic                  uint16 `note:"header id"`
	HdrSize                uint16 `note:"basic header size"` // excl. Magic+HdrSize
	SizeWithExtraData      byte   `note:"size up to and including 'extra data'"`
	Version                byte   `note:"archiver version number"`
	ExtractVersion         byte   `note:"minimum archiver version to extract"`
	HostOS                 byte   `note:"host OS",map:"hostOSes"` // XXX map
	Flags                  byte   `note:"arj flags"`
	SecurityVersion        byte   `note:"security version"`
	FileType               byte   `note:"file type",map:"fileTypes"` // XXX map
	Reserved               byte   `note:"reserved"`
	FileCTime              uint32 `note:"created time"`  // XXX time in "msdos-format"
	FileMTime              uint32 `note:"modified time"` // XXX time in "msdos-format"
	ArchiveSize            uint32 `note:"archive size for secured archive"`
	SecurityEnvelopePos    uint32 `note:"security envelope file position"`
	FilespecOffset         uint16 `note:"filespec position in filename"`
	SecurityEnvelopeLength uint16 `note:"length in bytes of security envelope data"`
	EncryptionVersion      byte   `note:"encryption version"`
	LastChapter            byte   `note:"last chapter"` // XXX
}

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

func ARJ(file *os.File) *ParsedLayout {

	if !isARJ(file) {
		return nil
	}

	res := ParsedLayout{}

	var err error
	res.Layout, err = parseARJMainHeader(file)

	if err != nil {
		fmt.Println("error", err)
		return nil
	}

	// XXX rest of arj
	return &res
}

func parseARJMainHeader(f *os.File) ([]Layout, error) {

	var err error
	offset, err := findARJHeader(f)
	if err != nil {
		return nil, err
	}

	h := arjMainHeader{}

	reader := io.Reader(f)

	if err := binary.Read(reader, binary.LittleEndian, &h); err != nil {
		return nil, err
	}
	if h.Magic != 60000 {
		return nil, fmt.Errorf("wrong magic word %04x", h.Magic)
	}

	archiveName := ""
	comment := ""

	if archiveName, err = zeroTerminatedASCII(f); err != nil {
		return nil, err
	}
	if comment, err = zeroTerminatedASCII(f); err != nil {
		return nil, err
	}

	return []Layout{
		Layout{
			Offset: offset,
			Length: int64(h.SizeWithExtraData) + 4,
			Type:   Group,
			Info:   "arj main header",
			Childs: []Layout{
				// XXX convert arjMainHeader into []Layout and add to Childs in return

				Layout{Offset: offset, Length: int64(len(archiveName)), Type: ASCIIZ, Info: "archive name"},
				Layout{Offset: offset + int64(len(archiveName)), Length: int64(len(comment)), Type: ASCIIZ, Info: "comment"},
				Layout{Offset: offset + int64(len(archiveName)+len(comment)), Length: 4, Type: Uint32le, Info: "crc32"},
				Layout{Offset: offset + int64(len(archiveName)+len(comment)) + 4, Length: 4, Type: Uint32le, Info: "ext header size"},
			},
		},
	}, nil

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
