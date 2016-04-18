package exe

// NOTE: on OSX, there is C headers in /usr/include/mach-o
// OS X Mach-O executable
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	machoCpuTypes = map[uint32]string{
		1:         "vax",
		2:         "romp",
		4:         "ns32032",
		5:         "ns32332",
		6:         "mc680x0",
		7:         "i386",
		8:         "mips",
		9:         "ns32532",
		11:        "hppa",
		12:        "arm",
		13:        "mc88000",
		14:        "sparc",
		15:        "i860-be",
		16:        "i860-le",
		17:        "rs6000",
		18:        "powerpc",
		255:       "veo",
		0x1000000: "abi64",
		0x1000007: "x86-64",
		0x1000018: "powerpc64",
	}
)

func MachO(file *os.File) (*parse.ParsedLayout, error) {

	if !isMachO(file) {
		return nil, nil
	}
	return parseMachO(file)
}

func isMachO(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b == 0xfeedfacf {
		return true
	}

	return false
}

func parseMachO(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)

	cpuType, _ := parse.ReadUint32le(file, pos+4)
	cpuTypeName := "?"
	if val, ok := machoCpuTypes[cpuType]; ok {
		cpuTypeName = val
	}

	res := parse.ParsedLayout{
		FileKind: parse.Executable,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 28, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
				{Offset: pos + 4, Length: 4, Info: "cpu type = " + cpuTypeName, Type: parse.Uint32le},
				{Offset: pos + 8, Length: 4, Info: "cpu subtype", Type: parse.Uint32le}, // XXX map ...
				{Offset: pos + 12, Length: 4, Info: "file type", Type: parse.Uint32le},  // XXX ?
				{Offset: pos + 16, Length: 4, Info: "n cmds", Type: parse.Uint32le},
				{Offset: pos + 20, Length: 4, Info: "size of cmds", Type: parse.Uint32le},
				{Offset: pos + 24, Length: 4, Info: "flags", Type: parse.Uint32le},
			}}}}

	/* XXX

	   struct segment_command {
	     uint32_t  cmd;
	     uint32_t  cmdsize;
	     char      segname[16];
	     uint32_t  vmaddr;
	     uint32_t  vmsize;
	     uint32_t  fileoff;
	     uint32_t  filesize;
	     vm_prot_t maxprot;
	     vm_prot_t initprot;
	     uint32_t  nsects;
	     uint32_t  flags;
	   };
	*/

	return &res, nil
}
