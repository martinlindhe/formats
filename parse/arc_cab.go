package parse

// STATUS 80% some polishing remains

import (
	"encoding/binary"
	"fmt"
	"os"
)

func CAB(file *os.File) (*ParsedLayout, error) {

	if !isCAB(file) {
		return nil, nil
	}
	return parseCAB(file)
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

func parseCAB(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)

	res := ParsedLayout{
		FileKind: Archive,
		Layout: []Layout{{
			Offset: pos,
			Length: 36, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: ASCII},
				{Offset: pos + 4, Length: 4, Info: "reserved 1", Type: Uint32le},
				{Offset: pos + 8, Length: 4, Info: "file size", Type: Uint32le},
				{Offset: pos + 12, Length: 4, Info: "reserved 2", Type: Uint32le},
				{Offset: pos + 16, Length: 4, Info: "offset to CFFILE", Type: Uint32le},
				{Offset: pos + 20, Length: 4, Info: "reserved 3", Type: Uint32le},
				{Offset: pos + 24, Length: 2, Info: "format version", Type: MinorMajor16le},
				{Offset: pos + 26, Length: 2, Info: "CFFOLDER entries", Type: Uint16le},
				{Offset: pos + 28, Length: 2, Info: "CFFILE entries", Type: Uint16le},
				{Offset: pos + 30, Length: 2, Info: "flags", Type: Uint16le},
				{Offset: pos + 32, Length: 2, Info: "set id", Type: Uint16le},
				{Offset: pos + 34, Length: 2, Info: "cabinet number", Type: Uint16le},
			}}}}

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

	dirEntries, _ := readUint16le(file, 26)

	cfDataBlocks := map[uint32]uint16{}

	for i := 0; i < int(dirEntries); i++ {
		chunk := Layout{
			Offset: pos,
			Length: 8,
			Info:   "CFFOLDER " + fmt.Sprintf("%d", i+1),
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "offset of first CFDATA block", Type: Uint32le},
				{Offset: pos + 4, Length: 2, Info: "CFDATA blocks", Type: Uint16le},
				{Offset: pos + 6, Length: 2, Info: "compression type", Type: Uint16le},
				// XXX:
				// u1  abReserve[];   /* (optional) per-folder reserved area */
			}}

		cfdataPos, _ := readUint32le(file, pos)
		cfdataBlocks, _ := readUint16le(file, pos+4)
		cfDataBlocks[cfdataPos] = cfdataBlocks

		pos += chunk.Length
		res.Layout = append(res.Layout, chunk)
	}

	fileEntries, _ := readUint16le(file, 28)

	cffOffset, _ := res.readUint32leFromInfo(file, "offset to CFFILE")
	if pos != int64(cffOffset) {
		fmt.Printf("cab: unexpected, offset = %x, cffOffset = %x\n", pos, cffOffset)
		pos = int64(cffOffset)
	}

	for i := 0; i < int(fileEntries); i++ {

		file.Seek(pos+16, os.SEEK_SET)
		_, nameLen, err := zeroTerminatedASCII(file)
		if err != nil {
			return nil, err
		}
		chunk := Layout{
			Offset: pos,
			Length: 16 + int64(nameLen),
			Info:   "CFFILE " + fmt.Sprintf("%d", i+1),
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "uncompressed size", Type: Uint32le},
				{Offset: pos + 4, Length: 4, Info: "uncompressed offset in folder", Type: Uint32le},
				{Offset: pos + 8, Length: 2, Info: "index in CFFOLDER", Type: Uint16le},
				{Offset: pos + 10, Length: 2, Info: "date stamp", Type: Uint16le},
				{Offset: pos + 12, Length: 2, Info: "time stamp", Type: Uint16le},
				{Offset: pos + 14, Length: 2, Info: "attributes", Type: Uint16le},
				{Offset: pos + 16, Length: int64(nameLen), Info: "name", Type: ASCIIZ},
			}}

		pos += chunk.Length
		res.Layout = append(res.Layout, chunk)
	}

	// map the compressed data
	for dataOffset, cnt := range cfDataBlocks {
		pos = int64(dataOffset)
		for i := 1; i < int(cnt); i++ {
			cbLen, _ := readUint16le(file, int64(dataOffset)+4)
			res.Layout = append(res.Layout, Layout{
				Offset: pos,
				Length: 8 + int64(cbLen),
				Info:   "CFDATA " + fmt.Sprintf("%d", i),
				Type:   Group,
				Childs: []Layout{
					{Offset: pos, Length: 4, Info: "checksum", Type: Uint32le},
					{Offset: pos + 4, Length: 2, Info: "compressed len", Type: Uint16le},
					{Offset: pos + 6, Length: 2, Info: "uncompressed len", Type: Uint16le},
					{Offset: pos + 8, Length: int64(cbLen), Info: "compressed data", Type: Bytes},
				}})
			pos += 8 + int64(cbLen)
		}
	}

	return &res, nil
}
