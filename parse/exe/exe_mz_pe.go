package exe

// PE exe (Win32, Win64)
// STATUS: 50%

// http://wiki.osdev.org/PE

import (
	"fmt"
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	peHeaderLen        = int64(24)
	peSectionHeaderLen = int64(40)
	peOptHeaderMainLen = int64(96)

	peTypes = map[uint16]string{
		0x10b: "PE32",
		0x20b: "PE32+ (64-bit)",
	}
	peMachines = map[uint16]string{
		0x14c:  "intel 386",
		0x8664: "AMD64",
	}
	peSubsystems = map[uint16]string{
		1: "native",
		2: "GUI",
		3: "console",
		5: "OS/2",
		7: "POSIX",
	}
	peRvaChunks = map[int64]string{
		0:  "exports",
		1:  "imports",
		2:  "resources",
		5:  "base reolcations",
		9:  "thread local storage",
		12: "import address table",
		14: "CLR header",
	}
)

// parses 32/64-bit Windows executables
func parseMZ_PEHeader(file *os.File, pos int64) ([]parse.Layout, error) {

	optHeaderSize, _ := parse.ReadUint16le(file, pos+20)
	numberOfSections, _ := parse.ReadUint16le(file, pos+6)
	machineName, _ := parse.ReadToMap(file, parse.Uint16le, pos+4, peMachines)
	res := []parse.Layout{{
		Offset: pos,
		Length: peHeaderLen,
		Info:   "PE header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "identifier", Type: parse.ASCIIZ},
			{Offset: pos + 4, Length: 2, Info: "machine = " + machineName, Type: parse.Uint16le},
			{Offset: pos + 6, Length: 2, Info: "number of sections", Type: parse.Uint16le},
			{Offset: pos + 8, Length: 4, Info: "timestamp", Type: parse.Uint32le}, // XXX format, convert, etc: var TimeDateStamp = NumberOfSections.RelativeToLittleEndianDateStamp("TimeDateStamp");
			{Offset: pos + 12, Length: 4, Info: "symbol table offset", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "symbol table entries", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 2, Info: "optional header size", Type: parse.Uint16le},
			{Offset: pos + 22, Length: 2, Info: "characteristics", Type: parse.Uint16le},
		}}}
	pos += peHeaderLen

	if optHeaderSize > 0 {
		optHeader := parsePEOptHeader(file, pos, optHeaderSize)
		res = append(res, optHeader)
		pos += optHeader.Length
	}

	sectionHeader := parse.Layout{
		Offset: pos,
		Length: int64(numberOfSections) * peSectionHeaderLen,
		Info:   "section header",
		Type:   parse.Group}

	for i := 0; i < int(numberOfSections); i++ {

		sectionName, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, pos, 256)
		rawDataSize, _ := parse.ReadUint32le(file, pos+16)
		rawDataOffset, _ := parse.ReadUint32le(file, pos+20)

		res = append(res, parse.Layout{
			Offset: int64(rawDataOffset),
			Length: int64(rawDataSize),
			Info:   "section " + sectionName,
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: int64(rawDataOffset), Length: int64(rawDataSize), Info: "data", Type: parse.Bytes},
			}})

		chunk := []parse.Layout{
			{Offset: pos, Length: 8, Info: "name", Type: parse.ASCIIZ},
			{Offset: pos + 8, Length: 4, Info: "virtual size", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 4, Info: "virtual address", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "raw data size", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "raw data offset", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 4, Info: "reallocations offset", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "linenumbers offset", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 2, Info: "reallocations count", Type: parse.Uint16le},
			{Offset: pos + 34, Length: 2, Info: "linenumbers count", Type: parse.Uint16le},
			{Offset: pos + 36, Length: 4, Info: "flags", Type: parse.Uint32le, Masks: []parse.Mask{
				// XXX fix bit map
				{Low: 0, Length: 1, Info: "0x00000020 = Code"},
				{Low: 0, Length: 1, Info: "0x00000040 = Initialized data"},
				{Low: 0, Length: 1, Info: "0x00000080 = Uninitialized data"},
				{Low: 0, Length: 1, Info: "0x00000200 = Info"},
				{Low: 0, Length: 1, Info: "0x02000000 = Discardable"},
				{Low: 0, Length: 1, Info: "0x10000000 = Shared"},
				{Low: 0, Length: 1, Info: "0x20000000 = Executable"},
				{Low: 0, Length: 1, Info: "0x40000000 = Readable"},
				{Low: 0, Length: 1, Info: "0x80000000 = Writeable"},
			}}}

		sectionHeader.Childs = append(sectionHeader.Childs, chunk...)
		pos += peSectionHeaderLen
	}

	res = append(res, sectionHeader)

	return res, nil
}

func parsePEOptHeader(file *os.File, pos int64, size uint16) parse.Layout {

	typeName, _ := parse.ReadToMap(file, parse.Uint16le, pos, peTypes)
	subsystemName, _ := parse.ReadToMap(file, parse.Uint16le, pos+68, peSubsystems)
	numberOfRva, _ := parse.ReadUint32le(file, pos+92)
	optHeader := parse.Layout{
		Offset: pos,
		Info:   "PE optional header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "type = " + typeName, Type: parse.Uint16le},
			{Offset: pos + 2, Length: 2, Info: "linker version", Type: parse.MajorMinor16le},
			{Offset: pos + 4, Length: 4, Info: "size of code", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "size of initialized data", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 4, Info: "size of uninitialized data", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "address of entry point", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "base of code", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 4, Info: "base of data", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "image base", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 4, Info: "section alignment", Type: parse.Uint32le},
			{Offset: pos + 36, Length: 4, Info: "file alignment", Type: parse.Uint32le},
			{Offset: pos + 40, Length: 4, Info: "os version", Type: parse.MajorMinor32le},
			{Offset: pos + 44, Length: 4, Info: "image version", Type: parse.MajorMinor32le},
			{Offset: pos + 48, Length: 4, Info: "subsystem version", Type: parse.MajorMinor32le},
			{Offset: pos + 52, Length: 4, Info: "win32 version value", Type: parse.Uint32le},
			{Offset: pos + 56, Length: 4, Info: "size of image", Type: parse.Uint32le},
			{Offset: pos + 60, Length: 4, Info: "size of headers", Type: parse.Uint32le},
			{Offset: pos + 64, Length: 4, Info: "checksum", Type: parse.Uint32le},
			{Offset: pos + 68, Length: 2, Info: "subsystem = " + subsystemName, Type: parse.Uint16le},
			{Offset: pos + 70, Length: 2, Info: "dll characteristics", Type: parse.Uint16le},
			{Offset: pos + 72, Length: 4, Info: "size of stack reserve", Type: parse.Uint32le},
			{Offset: pos + 76, Length: 4, Info: "size of stack commit", Type: parse.Uint32le},
			{Offset: pos + 80, Length: 4, Info: "size of heap reserve", Type: parse.Uint32le},
			{Offset: pos + 84, Length: 4, Info: "size of heap commit", Type: parse.Uint32le},
			{Offset: pos + 88, Length: 4, Info: "loader flags", Type: parse.Uint32le},
			{Offset: pos + 92, Length: 4, Info: "number of rva and sizes", Type: parse.Uint32le},
		}}
	pos += peOptHeaderMainLen

	if numberOfRva != 16 {
		fmt.Println("error: odd number of RVA:s = " + fmt.Sprintf("%d", numberOfRva))
	}

	ddLen := int64(8)

	optHeader.Length = peOptHeaderMainLen + (int64(numberOfRva) * ddLen)
	if optHeader.Length != int64(size) {
		fmt.Println("error: PE unexpected opt header len. expected ", size, " actual =", optHeader.Length)
	}

	for i := int64(0); i < int64(numberOfRva); i++ {

		info := "data directory " + fmt.Sprintf("%d", i)
		if val, ok := peRvaChunks[i]; ok {
			info = val
		}

		optHeader.Childs = append(optHeader.Childs, []parse.Layout{
			{Offset: pos, Length: 4, Info: info + " RVA", Type: parse.Uint32le},
			{Offset: pos + 4, Length: 4, Info: info + " size", Type: parse.Uint32le},
		}...)
		pos += 8
	}

	return optHeader
}
