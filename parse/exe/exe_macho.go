package exe

// MacOS Mach-O executable
// NOTE: on MacOS, there is C headers in /usr/include/mach-o

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
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

func MachO(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMachO(&c.Header) {
		return nil, nil
	}
	return parseMachO(c.File, c.ParsedLayout)
}

func isMachO(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[3] == 0xfe && b[2] == 0xed && b[1] == 0xfa && b[0] == 0xcf {
		return true
	}

	return false
}

func parseMachO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	cpuTypeName, _ := parse.ReadToMap(file, parse.Uint32le, pos+4, machoCpuTypes)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
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
		}}}

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

	return &pl, nil
}
