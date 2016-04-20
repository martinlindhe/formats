package exe

// MacOS Mach-O executable
// NOTE: on MacOS, there is C headers in /usr/include/mach-o
// TODO: handle the CIGAM byte ordered files (ppc, need samples)
// https://github.com/thetlk/Mach-O/tree/master/pymacho

// STATUS: 2%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

const (
	MH_MAGIC    = 0xfeedface
	MH_MAGIC_64 = 0xfeedfacf
	MH_CIGAM    = 0xcefaedfe
	MH_CIGAM_64 = 0xcffaedfe
)

var (
	machoMagicTypes = map[uint32]string{
		MH_MAGIC:    "MH_MAGIC",
		MH_MAGIC_64: "MH_MAGIC_64",
		MH_CIGAM:    "MH_CIGAM",
		MH_CIGAM_64: "MH_CIGAM_64",
	}
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
	machoFileTypes = map[uint32]string{
		1:  "object",       // relocatable object file
		2:  "execute",      // demand paged executable file
		3:  "fixed vm lib", // fixed VM shared library file
		4:  "core",         // core file
		5:  "preload",      // preloaded executable file
		6:  "dylib",        // dynamically bound shared library
		7:  "dylinker",     // dynamic link editor
		8:  "bundle",       // dynamically bound bundle file
		9:  "dylib stub",   // shared library stub for static linking only, no section contents
		10: "dsym",         // companion file with only debug sections
		11: "kext bundle",  // x86_64 kexts
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
	val := binary.LittleEndian.Uint32(b[:])
	if val == MH_MAGIC || val == MH_MAGIC_64 || val == MH_CIGAM || val == MH_CIGAM_64 {
		return true
	}
	return false
}

func parseMachO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	mhName, _ := parse.ReadToMap(file, parse.Uint32le, pos, machoMagicTypes)
	cpuTypeName, _ := parse.ReadToMap(file, parse.Uint32le, pos+4, machoCpuTypes)
	fileTypeName, _ := parse.ReadToMap(file, parse.Uint32le, pos+12, machoFileTypes)
	pl.FormatName = "mach-o " + cpuTypeName
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 28, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic = " + mhName, Type: parse.Uint32le},
			{Offset: pos + 4, Length: 4, Info: "cpu type = " + cpuTypeName, Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "cpu subtype", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 4, Info: "file type = " + fileTypeName, Type: parse.Uint32le},
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
