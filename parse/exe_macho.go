package parse

// NOTE: on OSX, there is C headers in /usr/include/mach-o
// OS X Mach-O executable
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

var (
	machoCpuTypes = map[uint32]string{
		1:         "VAX",
		2:         "ROMP",
		4:         "NS32032",
		5:         "NS32332",
		6:         "MC680x0",
		7:         "I386",
		8:         "MIPS",
		9:         "NS32532",
		11:        "HPPA",
		12:        "ARM",
		13:        "MC88000",
		14:        "SPARC",
		15:        "I860-be",
		16:        "I860-le",
		17:        "RS6000",
		18:        "POWERPC",
		255:       "VEO",
		0x1000000: "ABI64",
		0x1000007: "X86-64",
		0x1000018: "POWERPC64",
	}
)

func MachO(file *os.File) (*ParsedLayout, error) {

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

func parseMachO(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)

	cpuType, _ := ReadUint32le(file, pos+4)
	cpuTypeName := "?"
	if val, ok := machoCpuTypes[cpuType]; ok {
		cpuTypeName = val
	}

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: pos,
			Length: 28, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: Uint32le},
				{Offset: pos + 4, Length: 4, Info: "cpu type = " + cpuTypeName, Type: Uint32le},
				{Offset: pos + 8, Length: 4, Info: "cpu subtype", Type: Uint32le}, // XXX map ...
				{Offset: pos + 12, Length: 4, Info: "file type", Type: Uint32le},  // XXX ?
				{Offset: pos + 16, Length: 4, Info: "n cmds", Type: Uint32le},
				{Offset: pos + 20, Length: 4, Info: "size of cmds", Type: Uint32le},
				{Offset: pos + 24, Length: 4, Info: "flags", Type: Uint32le},
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
