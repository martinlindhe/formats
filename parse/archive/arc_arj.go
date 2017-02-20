package archive

// STATUS: 80%

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/martinlindhe/formats/parse"
)

const (
	arjMaxSFX        = 500000 // size of self-extracting prefix
	arjFirstHdrSize  = 0x1e
	arjCommentMax    = 2048
	arjFileNameMax   = 512
	arjHeaderSizeMax = (arjFirstHdrSize + 10 + arjFileNameMax + arjCommentMax)
)

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
	arjEncryptionVersion = map[byte]string{
		0: "old",
		1: "old",
		2: "new",
		3: "reserved",
		4: "40 bit key GOST",
	}
	arjMethod = map[byte]string{
		0: "stored",
		1: "compressed 1a (best)",
		2: "compressed 1b",
		3: "compressed 1c",
		4: "compressed 2 (fastest)",
		8: "no data, no CRC",
		9: "no data",
	}
	arjFileType = map[byte]string{
		0: "binary",
		1: "7-bit text",
		3: "directory",
		4: "volume label",
		5: "chapter label",
	}
)

// ARJ parses the arj format
func ARJ(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isARJ(c.Header) {
		return nil, nil
	}

	arj, err := parseARJ(c.File)

	c.ParsedLayout.FileKind = parse.Archive
	c.ParsedLayout.MimeType = "application/x-arj"
	c.ParsedLayout.Layout = arj

	return &c.ParsedLayout, err
}

func isARJ(b []byte) bool {

	return b[0] == 0x60 && b[1] == 0xea
}

// finds arj header and leaves file position at it
func findARJHeader(file *os.File) (int64, error) {
	reader := io.Reader(file)
	pos := int64(0)
	lastpos, _ := file.Seek(0, os.SEEK_END)
	lastpos -= 2
	if lastpos > arjMaxSFX {
		lastpos = arjMaxSFX
	}

	for ; pos < lastpos; pos++ {
		pos2, _ := file.Seek(pos, os.SEEK_SET)
		if pos != pos2 {
			fmt.Printf("warning: expected %d, got %d\n", pos, pos2)
		}

		var c byte
		if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
			return 0, err
		}

		for pos < lastpos {
			if c != 0x60 {
				// low order first
				if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
					return 0, err
				}
			} else {
				if err := binary.Read(reader, binary.LittleEndian, &c); err != nil {
					return 0, err
				}
				if c == 0xEA {
					break
				}
			}
			pos++
		}
		if pos >= lastpos {
			break
		}
		var headerSize uint16
		if err := binary.Read(reader, binary.LittleEndian, &headerSize); err != nil {
			return 0, err
		}
		if headerSize <= arjHeaderSizeMax {
			return pos, nil
		}
	}
	return 0, fmt.Errorf("could not find arj header in " + file.Name())
}

func parseARJ(f *os.File) ([]parse.Layout, error) {
	pos, err := findARJHeader(f)
	if err != nil {
		return nil, err
	}

	hostOSName, _ := parse.ReadToMap(f, parse.Uint8, pos+7, arjHostOS)
	encryptionName, _ := parse.ReadToMap(f, parse.Uint8, pos+32, arjEncryptionVersion)
	mainHeaderLen := int64(34)

	chunk := parse.Layout{
		Offset: pos,
		Type:   parse.Group,
		Info:   "main header",
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Type: parse.Uint16le, Info: "magic"},
			{Offset: pos + 2, Length: 2, Type: parse.Uint16le, Info: "header size"},
			{Offset: pos + 4, Length: 1, Type: parse.Uint8, Info: "size up to and including 'extra data'"},
			{Offset: pos + 5, Length: 1, Type: parse.MinorMajor8Five, Info: "archiver version"},
			{Offset: pos + 6, Length: 1, Type: parse.MinorMajor8Five, Info: "minimum archiver version to extract"},
			{Offset: pos + 7, Length: 1, Type: parse.Uint8, Info: "host OS = " + hostOSName},
			{Offset: pos + 8, Length: 1, Type: parse.Uint8, Info: "flags", Masks: []parse.Mask{
				{Low: 0, Length: 1, Info: "garbled"},
				{Low: 1, Length: 1, Info: "ansipage"}, // aka "OLD_SECURED_FLAG"
				{Low: 2, Length: 1, Info: "volume"},
				{Low: 3, Length: 1, Info: "arjprot"},
				{Low: 4, Length: 1, Info: "pathsym"},
				{Low: 5, Length: 1, Info: "backup"},
				{Low: 6, Length: 1, Info: "secured"},
				{Low: 7, Length: 1, Info: "altname"},
			}},
			{Offset: pos + 9, Length: 1, Type: parse.Uint8, Info: "security version"},
			{Offset: pos + 10, Length: 1, Type: parse.Uint8, Info: "file type"}, // must equal 2
			{Offset: pos + 11, Length: 1, Type: parse.Uint8, Info: "reserved"},
			{Offset: pos + 12, Length: 4, Type: parse.ArjDateTime, Info: "created time"},
			{Offset: pos + 16, Length: 4, Type: parse.ArjDateTime, Info: "modified time"},
			{Offset: pos + 20, Length: 4, Type: parse.Uint32le, Info: "archive size for secured archive"},
			{Offset: pos + 24, Length: 4, Type: parse.Uint32le, Info: "security envelope file position"},
			{Offset: pos + 28, Length: 2, Type: parse.Uint16le, Info: "filespec position in filename"},
			{Offset: pos + 30, Length: 2, Type: parse.Uint16le, Info: "length in bytes of security envelope data"},
			{Offset: pos + 32, Length: 1, Type: parse.Uint8, Info: "encryption version = " + encryptionName},
			{Offset: pos + 33, Length: 1, Type: parse.Uint8, Info: "last chapter"},
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

	_, archiveNameLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos+mainHeaderLen, arjFileNameMax)
	if err != nil {
		return nil, err
	}

	_, commentLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos+mainHeaderLen+archiveNameLen, arjCommentMax)
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
		{Offset: pos, Length: commentLen, Type: parse.ASCIIZ, Info: "archive comment"},
	}...)
	pos += commentLen

	chunk.Childs = append(chunk.Childs, []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "crc32"},
		{Offset: pos + 4, Length: 2, Type: parse.Uint16le, Info: "ext header size"},
	}...)

	f.Seek(pos+6, os.SEEK_SET)
	// NOTE: if ext header size > 0, it should follow here. currently unused in file format

	res := []parse.Layout{chunk}

	files, err := parseARJLocalFiles(f)
	if err != nil {
		return nil, err
	}
	res = append(res, files...)

	return res, nil
}

// parse local file headers until we reach eof
func parseARJLocalFiles(f *os.File) ([]parse.Layout, error) {
	res := []parse.Layout{}
	pos, _ := f.Seek(0, os.SEEK_CUR)

	for {
		magic, _ := parse.ReadUint16le(f, pos)
		if magic != 0xEA60 {
			log.Fatalf("Unexpected magic %04x at %04x", magic, pos)
		}
		length, _ := parse.ReadUint16le(f, pos+2)

		local := parse.Layout{
			Offset: pos,
			Length: 4,
			Type:   parse.Group,
			Info:   "local file header",
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Type: parse.Uint16le, Info: "magic"},
				{Offset: pos + 2, Length: 2, Type: parse.Uint16le, Info: "header size"},
			},
		}
		if length == 0 {
			local.Info = "eof marker"
		} else {
			localHostOSName, _ := parse.ReadToMap(f, parse.Uint8, pos+7, arjHostOS)
			methodName, _ := parse.ReadToMap(f, parse.Uint8, pos+9, arjMethod)
			fileTypeName, _ := parse.ReadToMap(f, parse.Uint8, pos+10, arjFileType)
			dataLength, _ := parse.ReadUint32le(f, pos+16)
			local.Childs = append(local.Childs, []parse.Layout{
				{Offset: pos + 4, Length: 1, Type: parse.Uint8, Info: "size up to and including 'extra data'"},
				{Offset: pos + 5, Length: 1, Type: parse.MinorMajor8Five, Info: "archiver version"},
				{Offset: pos + 6, Length: 1, Type: parse.MinorMajor8Five, Info: "minimum archiver version to extract"},
				{Offset: pos + 7, Length: 1, Type: parse.Uint8, Info: "host OS = " + localHostOSName},
				{Offset: pos + 8, Length: 1, Type: parse.Uint8, Info: "flags", Masks: []parse.Mask{
					{Low: 0, Length: 1, Info: "garbled"},
					{Low: 1, Length: 1, Info: "reserved"},
					{Low: 2, Length: 1, Info: "volume"},
					{Low: 3, Length: 1, Info: "extfile"},
					{Low: 4, Length: 1, Info: "pathsym"},
					{Low: 5, Length: 1, Info: "backup"},
					{Low: 6, Length: 2, Info: "reserved2"},
				}},
				{Offset: pos + 9, Length: 1, Type: parse.Uint8, Info: "method = " + methodName},
				{Offset: pos + 10, Length: 1, Type: parse.Uint8, Info: "file type = " + fileTypeName},
				{Offset: pos + 11, Length: 1, Type: parse.Uint8, Info: "reserved"},
				{Offset: pos + 12, Length: 4, Type: parse.ArjDateTime, Info: "modified time"},
				{Offset: pos + 16, Length: 4, Type: parse.Uint32le, Info: "compressed size"},
				{Offset: pos + 20, Length: 4, Type: parse.Uint32le, Info: "original size"},
				{Offset: pos + 24, Length: 4, Type: parse.Uint32le, Info: "original file's CRC"},
				{Offset: pos + 28, Length: 2, Type: parse.Uint16le, Info: "filespec position in filename"},
				{Offset: pos + 30, Length: 2, Type: parse.Uint16le, Info: "file access mode"},
				{Offset: pos + 32, Length: 1, Type: parse.Uint8, Info: "first chapter of file's lifespan"},
				{Offset: pos + 33, Length: 1, Type: parse.Uint8, Info: "last chapter of file's lifespan"},
			}...)
			local.Length += 30

			withExtData, _ := parse.ReadUint8(f, pos+4)
			if withExtData == 0x2E {
				local.Childs = append(local.Childs, []parse.Layout{
					{Offset: pos + 34, Length: 4, Type: parse.Uint32le, Info: "extended file position"},
					// XXX the following twelve bytes may be present in ARJ 2.62 and above:
					{Offset: pos + 38, Length: 4, Type: parse.ArjDateTime, Info: "accessed time"},
					{Offset: pos + 42, Length: 4, Type: parse.ArjDateTime, Info: "created time"},
					{Offset: pos + 46, Length: 4, Type: parse.Uint32le, Info: "original file size"},
				}...)
				local.Length += 16
			} else if withExtData == 0x1E {
				// no ext data
			} else {
				log.Fatalf("sample please. local file ext data = %02x", withExtData)
			}

			pos += local.Length

			_, fileNameLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos, arjFileNameMax)
			if err != nil {
				return nil, err
			}
			local.Childs = append(local.Childs, []parse.Layout{
				{Offset: pos, Length: fileNameLen, Type: parse.ASCIIZ, Info: "file name"},
			}...)
			pos += fileNameLen
			local.Length += fileNameLen

			_, commentLen, err := parse.ReadZeroTerminatedASCIIUntil(f, pos, arjCommentMax)
			if err != nil {
				return nil, err
			}

			local.Childs = append(local.Childs, []parse.Layout{
				{Offset: pos, Length: commentLen, Type: parse.ASCIIZ, Info: "file comment"},
			}...)
			pos += commentLen
			local.Length += commentLen

			local.Childs = append(local.Childs, []parse.Layout{
				{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "basic header crc32"},
				{Offset: pos + 4, Length: 2, Type: parse.Uint16le, Info: "ext header size"},
			}...)
			pos += 6
			local.Length += 6
			// NOTE: if ext header size > 0, it should follow here. currently unused in file format

			local.Childs = append(local.Childs, []parse.Layout{
				{Offset: pos, Length: int64(dataLength), Type: parse.Bytes, Info: "compressed data"},
			}...)
			pos += int64(dataLength)
			local.Length += int64(dataLength)
		}
		res = append(res, local)
		if length == 0 {
			break
		}
	}
	return res, nil
}
