package archive

// STATUS 80% some polishing remains

import (
	"encoding/binary"
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func CAB(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isCAB(file) {
		return nil, nil
	}
	return parseCAB(file, pl)
}

func isCAB(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'M' || b[1] != 'S' || b[2] != 'C' || b[3] != 'F' {
		return false
	}

	return true
}

func parseCAB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 36, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 4, Info: "reserved 1", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "file size", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 4, Info: "reserved 2", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "offset to CFFILE", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "reserved 3", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 2, Info: "format version", Type: parse.MinorMajor16le},
			{Offset: pos + 26, Length: 2, Info: "CFFOLDER entries", Type: parse.Uint16le},
			{Offset: pos + 28, Length: 2, Info: "CFFILE entries", Type: parse.Uint16le},
			{Offset: pos + 30, Length: 2, Info: "flags", Type: parse.Uint16le},
			{Offset: pos + 32, Length: 2, Info: "set id", Type: parse.Uint16le},
			{Offset: pos + 34, Length: 2, Info: "cabinet number", Type: parse.Uint16le},
		}}}

	/* XXX
	   u2  cbCFHeader;       // (optional) size of per-cabinet reserved area
	   u1  cbCFFolder;       // (optional) size of per-folder reserved area
	   u1  cbCFData;         // (optional) size of per-datablock reserved area
	   u1  abReserve[];      // (optional) per-cabinet reserved area
	   u1  szCabinetPrev[];  // (optional) name of previous cabinet file
	   u1  szDiskPrev[];     // (optional) name of previous disk
	   u1  szCabinetNext[];  // (optional) name of next cabinet file
	   u1  szDiskNext[];     // (optional) name of next disk
	*/

	pos += 36 // XXX

	dirEntries, _ := parse.ReadUint16le(file, 26)

	cfDataBlocks := map[uint32]uint16{}

	for i := 0; i < int(dirEntries); i++ {
		chunk := parse.Layout{
			Offset: pos,
			Length: 8,
			Info:   "CFFOLDER " + fmt.Sprintf("%d", i+1),
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "offset of first CFDATA block", Type: parse.Uint32le},
				{Offset: pos + 4, Length: 2, Info: "CFDATA blocks", Type: parse.Uint16le},
				{Offset: pos + 6, Length: 2, Info: "compression type", Type: parse.Uint16le},
				// XXX:
				// u1  abReserve[];   /* (optional) per-folder reserved area */
			}}

		cfdataPos, _ := parse.ReadUint32le(file, pos)
		cfdataBlocks, _ := parse.ReadUint16le(file, pos+4)
		cfDataBlocks[cfdataPos] = cfdataBlocks
		pl.Layout = append(pl.Layout, chunk)
		pos += chunk.Length
	}

	fileEntries, _ := parse.ReadUint16le(file, 28)

	cffOffset, _ := pl.ReadUint32leFromInfo(file, "offset to CFFILE")
	if pos != int64(cffOffset) {
		fmt.Printf("cab: unexpected, offset = %x, cffOffset = %x\n", pos, cffOffset)
		pos = int64(cffOffset)
	}

	for i := 0; i < int(fileEntries); i++ {

		_, nameLen, err := parse.ReadZeroTerminatedASCIIUntil(file, pos+16, 256)
		if err != nil {
			return nil, err
		}
		chunk := parse.Layout{
			Offset: pos,
			Length: 16 + int64(nameLen),
			Info:   "CFFILE " + fmt.Sprintf("%d", i+1),
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "uncompressed size", Type: parse.Uint32le},
				{Offset: pos + 4, Length: 4, Info: "uncompressed offset in folder", Type: parse.Uint32le},
				{Offset: pos + 8, Length: 2, Info: "index in CFFOLDER", Type: parse.Uint16le},
				{Offset: pos + 10, Length: 2, Info: "date stamp", Type: parse.Uint16le},
				{Offset: pos + 12, Length: 2, Info: "time stamp", Type: parse.Uint16le},
				{Offset: pos + 14, Length: 2, Info: "attributes", Type: parse.Uint16le},
				{Offset: pos + 16, Length: int64(nameLen), Info: "name", Type: parse.ASCIIZ},
			}}
		pos += chunk.Length
		pl.Layout = append(pl.Layout, chunk)
	}

	// map the compressed data
	for dataOffset, cnt := range cfDataBlocks {
		pos = int64(dataOffset)
		for i := 1; i < int(cnt); i++ {
			cbLen, _ := parse.ReadUint16le(file, int64(dataOffset)+4)
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Length: 8 + int64(cbLen),
				Info:   "CFDATA " + fmt.Sprintf("%d", i),
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 4, Info: "checksum", Type: parse.Uint32le},
					{Offset: pos + 4, Length: 2, Info: "compressed len", Type: parse.Uint16le},
					{Offset: pos + 6, Length: 2, Info: "uncompressed len", Type: parse.Uint16le},
					{Offset: pos + 8, Length: int64(cbLen), Info: "compressed data", Type: parse.Bytes},
				}})
			pos += 8 + int64(cbLen)
		}
	}

	return &pl, nil
}
