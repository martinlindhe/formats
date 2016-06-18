package exe

// Linear Executable (Win VxD:s, and OS/2)
// http://fileformats.archiveteam.org/wiki/Linear_Executable
// http://wiki.osdev.org/LE

// STATUS: 2%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	leCpuTypes = map[uint16]string{
		1:    "intel 80286+",
		2:    "intel 80386+",
		3:    "intel 80486+",
		4:    "intel 80586+",
		0x20: "intel i860 (N10) or compatible",
		0x21: `intel "N11" or compatible`,
		0x40: "MIPS Mark I (R2000, R3000) or compatible",
		0x41: "MIPS Mark II (R6000) or compatible",
		0x42: "MIPS Mark III (R4000) or compatible",
	}
	leTargetOSes = map[uint16]string{
		1: "os/2",
		2: "windows",
		3: "dos 4.x",
		4: "windows 386",
	}
)

func parseMzLeHeader(file *os.File, pos int64) ([]parse.Layout, error) {

	cpuType, _ := parse.ReadToMap(file, parse.Uint16le, pos+8, leCpuTypes)
	targetOS, _ := parse.ReadToMap(file, parse.Uint16le, pos+10, leTargetOSes)

	res := []parse.Layout{{
		Offset: pos,
		Length: 176, // XXX
		Info:   "external header = LE",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
			{Offset: pos + 2, Length: 1, Info: "byte order", Type: parse.Uint8},
			{Offset: pos + 3, Length: 1, Info: "word order", Type: parse.Uint8},
			{Offset: pos + 4, Length: 4, Info: "executable format level", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 2, Info: "cpu type = " + cpuType, Type: parse.Uint16le},
			{Offset: pos + 10, Length: 2, Info: "target os = " + targetOS, Type: parse.Uint16le},
			{Offset: pos + 12, Length: 4, Info: "module version", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "module type flags", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "memory page count", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 4, Info: "initial cs", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "initial eip", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 4, Info: "initial ss", Type: parse.Uint32le},
			{Offset: pos + 36, Length: 4, Info: "initial esp", Type: parse.Uint32le},
			{Offset: pos + 40, Length: 4, Info: "memory page size", Type: parse.Uint32le},
			{Offset: pos + 44, Length: 4, Info: "bytes on last page", Type: parse.Uint32le},
			{Offset: pos + 48, Length: 4, Info: "fix-up section size", Type: parse.Uint32le},
			{Offset: pos + 52, Length: 4, Info: "fix-up secrion checksum", Type: parse.Uint32le},
			{Offset: pos + 56, Length: 4, Info: "loader section size", Type: parse.Uint32le},
			{Offset: pos + 60, Length: 4, Info: "loader section checksum", Type: parse.Uint32le},
			{Offset: pos + 64, Length: 4, Info: "object table offset", Type: parse.Uint32le}, // XXX decode
			{Offset: pos + 68, Length: 4, Info: "object table count", Type: parse.Uint32le},
			{Offset: pos + 72, Length: 4, Info: "object page map offset", Type: parse.Uint32le},
			{Offset: pos + 76, Length: 4, Info: "object iterate data map offset", Type: parse.Uint32le},
			{Offset: pos + 80, Length: 4, Info: "resource table offset", Type: parse.Uint32le},
			{Offset: pos + 84, Length: 4, Info: "resource table entries", Type: parse.Uint32le},
			{Offset: pos + 88, Length: 4, Info: "resident names table offset", Type: parse.Uint32le}, // XXX decode
			{Offset: pos + 92, Length: 4, Info: "entry table offset", Type: parse.Uint32le},          // XXX decode
			{Offset: pos + 96, Length: 4, Info: "module directives table offset", Type: parse.Uint32le},
			{Offset: pos + 100, Length: 4, Info: "module directives entries", Type: parse.Uint32le},
			{Offset: pos + 104, Length: 4, Info: "fix-up page table offset", Type: parse.Uint32le},
			{Offset: pos + 108, Length: 4, Info: "fix-up record table offset", Type: parse.Uint32le},
			{Offset: pos + 112, Length: 4, Info: "imported modules name table offset", Type: parse.Uint32le},
			{Offset: pos + 116, Length: 4, Info: "imported modules count", Type: parse.Uint32le},
			{Offset: pos + 120, Length: 4, Info: "imported procedure name table offset", Type: parse.Uint32le},
			{Offset: pos + 124, Length: 4, Info: "per-page checksum table offset", Type: parse.Uint32le},
			{Offset: pos + 128, Length: 4, Info: "data pages offset from top of file", Type: parse.Uint32le},
			{Offset: pos + 132, Length: 4, Info: "preload page count", Type: parse.Uint32le},
			{Offset: pos + 136, Length: 4, Info: "non-resident names table offset from top of file", Type: parse.Uint32le},
			{Offset: pos + 140, Length: 4, Info: "non-resident names table length", Type: parse.Uint32le},
			{Offset: pos + 144, Length: 4, Info: "non-resident names table checksum", Type: parse.Uint32le},
			{Offset: pos + 148, Length: 4, Info: "automatic data object", Type: parse.Uint32le},
			{Offset: pos + 152, Length: 4, Info: "debug information offset", Type: parse.Uint32le},
			{Offset: pos + 156, Length: 4, Info: "debug information length", Type: parse.Uint32le},
			{Offset: pos + 160, Length: 4, Info: "preload instance pages number", Type: parse.Uint32le},
			{Offset: pos + 164, Length: 4, Info: "demand instance pages number", Type: parse.Uint32le},
			{Offset: pos + 168, Length: 4, Info: "extra heap allocation", Type: parse.Uint32le},
			{Offset: pos + 172, Length: 4, Info: "unknown", Type: parse.Uint32le},
		}}}

	return res, nil
}
