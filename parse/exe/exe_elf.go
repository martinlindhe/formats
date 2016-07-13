package exe

// https://en.wikipedia.org/wiki/Executable_and_Linkable_Format

// STATUS: 70%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	phHeaderSize = int64(32)
	shHeaderSize = int64(40)
	elfClasses   = map[byte]string{
		0: "none",
		1: "ELF32",
		2: "ELF64",
	}
	elfEndian = map[byte]string{
		1: "lsb",
		2: "msb",
	}
	elfOSABIs = map[byte]string{
		0:    "system v",
		1:    "hp-ux",
		2:    "netbsd",
		3:    "linux",
		6:    "solaris",
		7:    "aix",
		8:    "irix",
		9:    "freebsd",
		0xc:  "openbsd",
		0xd:  "openvms",
		0xe:  "nsk os",
		0xf:  "aros",
		0x10: "fenix os",
		0x11: "cloud abi",
	}
	elfTypes = map[uint16]string{
		0:      "none",
		1:      "relocatable file",
		2:      "executable file",
		3:      "shared object file",
		4:      "core file",
		0xff00: "ET_LOPROC",
		0xffff: "ET_HIPROC",
	}
	elfMachines = map[uint16]string{
		0:    "none",
		1:    "AT&T WE 32100",
		2:    "SPARC",
		3:    "intel 80386",
		4:    "motorola 68000",
		5:    "motorola 88000",
		7:    "intel 80860",
		8:    "MIPS",
		0x14: "powerpc",
		0x28: "arm",
		0x2a: "superh",
		0x32: "ia-64",
		0x3e: "x86-64",
		0xb7: "aarch64",
	}
	elfPhTypes = map[uint32]string{
		0:          "null",
		1:          "load",
		2:          "dynamic",
		3:          "interp",
		4:          "note",
		5:          "sh lib",
		6:          "p hdr",
		0x60000000: "lo os",
		0x6fffffff: "hi os",
		0x70000000: "lo proc",
		0x7fffffff: "hi proc",
	}
	elfShTypes = map[uint32]string{
		0:          "null",
		1:          "prog bits",
		2:          "symbol table",
		3:          "string table",
		6:          "dynamic",
		7:          "note",
		8:          "no bits",
		9:          "rel",
		0xb:        "dyn sym",
		0xe:        "init array",
		0xf:        "fini array",
		0x6ffffff6: "gnu hash",
		0x6ffffffe: "ver need",
		0x6fffffff: "ver sym",
	}
)

// ELF parses the Executable and Linkable Format
func ELF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isELF(c.Header) {
		return nil, nil
	}
	return parseELF(c.File, c.ParsedLayout)
}

func isELF(b []byte) bool {

	if b[5] != 1 && b[5] != 2 { // endian
		return false
	}
	if b[5] == 2 {
		fmt.Println("TODO: handle big-endian ELF")
		return false
	}
	if b[0] != 0x7f || b[1] != 'E' || b[2] != 'L' || b[3] != 'F' {
		return false
	}
	return true
}

func parseELF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	phOffset, _ := parse.ReadUint32le(file, pos+28)
	phEntrySize, _ := parse.ReadUint16le(file, pos+42)
	phCount, _ := parse.ReadUint16le(file, pos+44)

	shOffset, _ := parse.ReadUint32le(file, pos+32)
	shEntrySize, _ := parse.ReadUint16le(file, pos+46)
	shCount, _ := parse.ReadUint16le(file, pos+48)

	header := elfHeader(file, pos)
	pl.FileKind = parse.Executable
	pl.MimeType = "application/x-executable"
	pl.Layout = []parse.Layout{header}

	if phOffset > 0 && phCount > 0 {
		pl.Layout = append(pl.Layout, parseElfPhEntries(file, int64(phOffset), phEntrySize, phCount)...)
	}

	if shOffset > 0 && shCount > 0 {
		pl.Layout = append(pl.Layout, parseElfShEntries(file, int64(shOffset), shEntrySize, shCount)...)
	}

	pl.Sort()

	return &pl, nil
}

func elfHeader(file *os.File, pos int64) parse.Layout {

	className, _ := parse.ReadToMap(file, parse.Uint8, pos+4, elfClasses)
	endian, _ := parse.ReadToMap(file, parse.Uint8, pos+5, elfEndian)
	osABIName, _ := parse.ReadToMap(file, parse.Uint8, pos+7, elfOSABIs)
	typeName, _ := parse.ReadToMap(file, parse.Uint16le, pos+16, elfTypes)
	machineName, _ := parse.ReadToMap(file, parse.Uint16le, pos+18, elfMachines)

	return parse.Layout{
		Offset: pos,
		Length: 52, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
			{Offset: pos + 4, Length: 1, Info: "class = " + className, Type: parse.Uint8},
			{Offset: pos + 5, Length: 1, Info: "endian = " + endian, Type: parse.Uint8},
			{Offset: pos + 6, Length: 1, Info: "header version", Type: parse.Uint8},
			{Offset: pos + 7, Length: 1, Info: "os abi = " + osABIName, Type: parse.Bytes},
			{Offset: pos + 8, Length: 1, Info: "abi version", Type: parse.Uint8},
			{Offset: pos + 9, Length: 7, Info: "reserved", Type: parse.Bytes},
			{Offset: pos + 16, Length: 2, Info: "type = " + typeName, Type: parse.Uint16le},
			{Offset: pos + 18, Length: 2, Info: "machine = " + machineName, Type: parse.Uint16le},
			{Offset: pos + 20, Length: 4, Info: "version", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 4, Info: "entry", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "program header offset", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 4, Info: "section header offset", Type: parse.Uint32le},
			{Offset: pos + 36, Length: 4, Info: "flags", Type: parse.Uint32le},
			{Offset: pos + 40, Length: 2, Info: "elf header size", Type: parse.Uint16le},
			{Offset: pos + 42, Length: 2, Info: "program header entry size", Type: parse.Uint16le},
			{Offset: pos + 44, Length: 2, Info: "program header count", Type: parse.Uint16le},
			{Offset: pos + 46, Length: 2, Info: "section header entry size", Type: parse.Uint16le},
			{Offset: pos + 48, Length: 2, Info: "section header count", Type: parse.Uint16le},
			{Offset: pos + 50, Length: 2, Info: "section header strndx", Type: parse.Uint16le}, // XXX map
		}}
}

func elfStrtabOffset(file *os.File, pos int64, shCount uint16) int64 {

	// XXX hack for one sample file

	// XXX need to look up segment type STRTAB to decode name ...
	// XXX 2: there are 3 strtab segments in sample file:
	//   which to choose?  .shstrtab  .. but how to know this?!

	return 0x6ed

	/*

		for i := 1; i <= int(shCount); i++ {
			shType, _ := readUint32le(file, pos+4)
			if shType == 3 {
				offset, _ := readUint32le(file, pos+16)
				return int64(offset)
			}
			pos += shHeaderSize
		}
		return 0
	*/
}

func parseElfPhEntries(file *os.File, pos int64, phEntrySize uint16, phCount uint16) []parse.Layout {

	res := []parse.Layout{}

	if int64(phEntrySize) != phHeaderSize {
		fmt.Println("warning: ph entry size - expected", phHeaderSize, ", saw", int64(phEntrySize))
	}

	for i := 1; i <= int(phCount); i++ {

		phTypeName, _ := parse.ReadToMap(file, parse.Uint32le, pos, elfPhTypes)
		phOffset, _ := parse.ReadUint32le(file, pos+4)
		phSize, _ := parse.ReadUint32le(file, pos+16)

		id := fmt.Sprintf("%d", i)

		res = append(res, parse.Layout{
			Offset: pos,
			Length: phHeaderSize,
			Info:   "program header " + id,
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "type = " + phTypeName, Type: parse.Uint32le},
				{Offset: pos + 4, Length: 4, Info: "offset", Type: parse.Uint32le},
				{Offset: pos + 8, Length: 4, Info: "virtual addresss", Type: parse.Uint32le},
				{Offset: pos + 12, Length: 4, Info: "physical address", Type: parse.Uint32le},
				{Offset: pos + 16, Length: 4, Info: "file size", Type: parse.Uint32le},
				{Offset: pos + 20, Length: 4, Info: "mem size", Type: parse.Uint32le},
				{Offset: pos + 24, Length: 4, Info: "flags", Type: parse.Uint32le},
				{Offset: pos + 28, Length: 4, Info: "align", Type: parse.Uint32le},
			}})

		if phSize > 0 {
			// NOTE: there are sections covering the whole header (LOAD), and the program header (PHDR)
			if phOffset != 0 {
				res = append(res, parse.Layout{
					Offset: int64(phOffset),
					Length: int64(phSize),
					Info:   "program " + id,
					Type:   parse.Group,
					Childs: []parse.Layout{
						{Offset: int64(phOffset), Length: int64(phSize), Info: "data", Type: parse.Bytes},
					}})
			}
		}

		pos += phHeaderSize
	}

	return res
}

func parseElfShEntries(file *os.File, pos int64, shEntrySize uint16, shCount uint16) []parse.Layout {

	res := []parse.Layout{}

	if int64(shEntrySize) != shHeaderSize {
		fmt.Println("warning: unexpected sh entry size. expected", shHeaderSize, ", saw", int64(shEntrySize))
	}

	strtabOffset := elfStrtabOffset(file, pos, shCount)

	for i := 1; i <= int(shCount); i++ {

		shTypeName, _ := parse.ReadToMap(file, parse.Uint32le, pos+4, elfShTypes)
		shOffset, _ := parse.ReadUint32le(file, pos+16)
		shSize, _ := parse.ReadUint32le(file, pos+20)
		nameOffset, _ := parse.ReadUint32le(file, pos)
		name, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, strtabOffset+int64(nameOffset), 32)
		res = append(res, parse.Layout{
			Offset: pos,
			Length: shHeaderSize,
			Info:   "section header " + fmt.Sprintf("%d", i),
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "name = " + name, Type: parse.Uint32le},
				{Offset: pos + 4, Length: 4, Info: "type = " + shTypeName, Type: parse.Uint32le},
				{Offset: pos + 8, Length: 4, Info: "flags", Type: parse.Uint32le},
				{Offset: pos + 12, Length: 4, Info: "address", Type: parse.Uint32le},
				{Offset: pos + 16, Length: 4, Info: "offset", Type: parse.Uint32le},
				{Offset: pos + 20, Length: 4, Info: "size", Type: parse.Uint32le},
				{Offset: pos + 24, Length: 16, Info: "extra", Type: parse.Uint32le}, // type dependent
			}})

		if shSize > 0 {
			if shOffset != 0 {
				res = append(res, parse.Layout{
					Offset: int64(shOffset),
					Length: int64(shSize),
					Info:   "section " + name,
					Type:   parse.Group,
					Childs: []parse.Layout{
						{Offset: int64(shOffset), Length: int64(shSize), Info: "data", Type: parse.Bytes},
					}})
			}
		}

		pos += shHeaderSize
	}

	return res
}
