package parse

// STATUS 2%

import (
	"encoding/binary"
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

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 36, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 4, Info: "magic", Type: ASCII},
			Layout{Offset: 4, Length: 4, Info: "reserved1", Type: Uint32le},
			Layout{Offset: 8, Length: 4, Info: "file size", Type: Uint32le},
			Layout{Offset: 12, Length: 4, Info: "reserved2", Type: Uint32le},
			Layout{Offset: 16, Length: 4, Info: "offset to CFFILE", Type: Uint32le},
			Layout{Offset: 20, Length: 4, Info: "reserved3", Type: Uint32le},

			Layout{Offset: 24, Length: 2, Info: "format version", Type: MinorMajor16},
			Layout{Offset: 26, Length: 2, Info: "number of CFFOLDER entries", Type: Uint16le},

			Layout{Offset: 28, Length: 2, Info: "number of CFFILE entries", Type: Uint16le},
			Layout{Offset: 30, Length: 2, Info: "flags", Type: Uint16le},
			Layout{Offset: 32, Length: 2, Info: "set id", Type: Uint16le},

			Layout{Offset: 34, Length: 2, Info: "cabinet number", Type: Uint16le},

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
		},
	})

	return &res, nil
}
