package parse

// STATUS 20%

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

	res := ParsedLayout{}

	offset := int64(0)

	res.Layout = append(res.Layout, Layout{
		Offset: offset,
		Length: 36, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: offset, Length: 4, Info: "magic", Type: ASCII},
			Layout{Offset: offset + 4, Length: 4, Info: "reserved 1", Type: Uint32le},
			Layout{Offset: offset + 8, Length: 4, Info: "file size", Type: Uint32le},
			Layout{Offset: offset + 12, Length: 4, Info: "reserved 2", Type: Uint32le},
			Layout{Offset: offset + 16, Length: 4, Info: "offset to CFFILE", Type: Uint32le},
			Layout{Offset: offset + 20, Length: 4, Info: "reserved 3", Type: Uint32le},
			Layout{Offset: offset + 24, Length: 2, Info: "format version", Type: MinorMajor16},
			Layout{Offset: offset + 26, Length: 2, Info: "number of CFFOLDER entries", Type: Uint16le},
			Layout{Offset: offset + 28, Length: 2, Info: "number of CFFILE entries", Type: Uint16le},
			Layout{Offset: offset + 30, Length: 2, Info: "flags", Type: Uint16le},
			Layout{Offset: offset + 32, Length: 2, Info: "set id", Type: Uint16le},
			Layout{Offset: offset + 34, Length: 2, Info: "cabinet number", Type: Uint16le},
		},
	})

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

	offset += 36 // XXX

	dirEntries, _ := readUint16le(file, 26)

	for i := 0; i < int(dirEntries); i++ {
		chunk := Layout{
			Offset: offset,
			Length: 8,
			Info:   "CFFOLDER " + fmt.Sprintf("%d", i+1),
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: offset, Length: 4, Info: "offset of first CFDATA block", Type: Uint32le},
				Layout{Offset: offset + 4, Length: 2, Info: "number of CFDATA blocks", Type: Uint16le},
				Layout{Offset: offset + 6, Length: 2, Info: "compression type", Type: Uint16le},
				// XXX:
				// u1  abReserve[];   /* (optional) per-folder reserved area */
			}}

		offset += chunk.Length
		res.Layout = append(res.Layout, chunk)
	}

	fileEntries, _ := readUint16le(file, 28)

	cffOffset, _ := res.readUint32leFromInfo(file, "offset to CFFILE")
	if offset != int64(cffOffset) {
		fmt.Printf("cab: unexpected, offset = %x, cffOffset = %x\n", offset, cffOffset)
		offset = int64(cffOffset)
	}

	for i := 0; i < int(fileEntries); i++ {

		file.Seek(offset+16, os.SEEK_SET)
		_, nameLen, err := zeroTerminatedASCII(file)
		if err != nil {
			return nil, err
		}
		chunk := Layout{
			Offset: offset,
			Length: 16 + int64(nameLen),
			Info:   "CFFILE " + fmt.Sprintf("%d", i+1),
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: offset, Length: 4, Info: "uncompressed size", Type: Uint32le},
				Layout{Offset: offset + 4, Length: 4, Info: "uncompressed offset in folder", Type: Uint32le},
				Layout{Offset: offset + 8, Length: 2, Info: "index in CFFOLDER", Type: Uint16le},
				Layout{Offset: offset + 10, Length: 2, Info: "date stamp", Type: Uint16le},
				Layout{Offset: offset + 12, Length: 2, Info: "time stamp", Type: Uint16le},
				Layout{Offset: offset + 14, Length: 2, Info: "attributes", Type: Uint16le},
				Layout{Offset: offset + 16, Length: int64(nameLen), Info: "name", Type: ASCIIZ},
			}}

		offset += chunk.Length
		res.Layout = append(res.Layout, chunk)
	}

	// XXX actual file data remains

	return &res, nil
}
