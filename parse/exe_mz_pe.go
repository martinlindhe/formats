package parse

// PE exe (Win32, Win64)
// STATUS: 50%

// http://wiki.osdev.org/PE

import (
	"fmt"
	"os"
)

var (
	peTypes = map[uint16]string{
		0x10b: "PE32",
		0x20b: "PE32+ (64-bit)",
	}
	peMachines = map[uint16]string{
		0x14c:  "Intel 386",
		0x8664: "AMD64",
	}
	peSubsystems = map[uint16]string{
		1: "Native",
		2: "GUI",
		3: "Console",
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
func parseMZ_PEHeader(file *os.File, offset int64) ([]Layout, error) {

	peHeaderLen := int64(24)
	sectionHeaderLen := int64(40)
	optHeaderSize, _ := readUint16le(file, offset+20)
	numberOfSections, _ := readUint16le(file, offset+6)

	machine, _ := readUint16le(file, offset+4)

	machineName := "?"
	if val, ok := peMachines[machine]; ok {
		machineName = val
	}

	res := []Layout{{
		Offset: offset,
		Length: peHeaderLen,
		Info:   "PE header",
		Type:   Group,
		Childs: []Layout{
			{Offset: offset, Length: 4, Info: "identifier", Type: ASCIIZ},
			{Offset: offset + 4, Length: 2, Info: "machine = " + machineName, Type: Uint16le},
			{Offset: offset + 6, Length: 2, Info: "number of sections", Type: Uint16le},
			{Offset: offset + 8, Length: 4, Info: "timestamp", Type: Uint32le}, // XXX format, convert, etc: var TimeDateStamp = NumberOfSections.RelativeToLittleEndianDateStamp("TimeDateStamp");
			{Offset: offset + 12, Length: 4, Info: "symbol table offset", Type: Uint32le},
			{Offset: offset + 16, Length: 4, Info: "symbol table entries", Type: Uint32le},
			{Offset: offset + 20, Length: 2, Info: "optional header size", Type: Uint16le},
			{Offset: offset + 22, Length: 2, Info: "characteristics", Type: Uint16le},
		}}}
	offset += peHeaderLen

	if optHeaderSize > 0 {
		optHeader := parsePEOptHeader(file, offset, optHeaderSize)
		res = append(res, optHeader)
		offset += optHeader.Length
	}

	sectionHeader := Layout{
		Offset: offset,
		Length: int64(numberOfSections) * sectionHeaderLen,
		Info:   "section header",
		Type:   Group,
	}

	for i := 0; i < int(numberOfSections); i++ {

		sectionName, _, _ := readZeroTerminatedASCII(file, offset)
		rawDataSize, _ := readUint32le(file, offset+16)
		rawDataOffset, _ := readUint32le(file, offset+20)

		res = append(res, Layout{
			Offset: int64(rawDataOffset),
			Length: int64(rawDataSize),
			Info:   "section " + sectionName,
			Type:   Group,
			Childs: []Layout{
				{Offset: int64(rawDataOffset), Length: int64(rawDataSize), Info: "data", Type: Bytes},
			}})

		chunk := []Layout{
			{Offset: offset, Length: 8, Info: "name", Type: ASCIIZ},
			{Offset: offset + 8, Length: 4, Info: "virtual size", Type: Uint32le},
			{Offset: offset + 12, Length: 4, Info: "virtual address", Type: Uint32le},
			{Offset: offset + 16, Length: 4, Info: "raw data size", Type: Uint32le},
			{Offset: offset + 20, Length: 4, Info: "raw data offset", Type: Uint32le},
			{Offset: offset + 24, Length: 4, Info: "reallocations offset", Type: Uint32le},
			{Offset: offset + 28, Length: 4, Info: "linenumbers offset", Type: Uint32le},
			{Offset: offset + 32, Length: 2, Info: "reallocations count", Type: Uint16le},
			{Offset: offset + 34, Length: 2, Info: "linenumbers count", Type: Uint16le},
			{Offset: offset + 36, Length: 4, Info: "flags", Type: Uint32le, Masks: []Mask{
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
			}},
		}
		sectionHeader.Childs = append(sectionHeader.Childs, chunk...)
		offset += sectionHeaderLen
	}

	res = append(res, sectionHeader)

	return res, nil
}

func parsePEOptHeader(file *os.File, offset int64, size uint16) Layout {

	typeId, _ := readUint16le(file, offset)
	typeName := "?"
	if val, ok := peTypes[typeId]; ok {
		typeName = val
	}

	subsystem, _ := readUint16le(file, offset+68)
	subsystemName := "?"
	if val, ok := peSubsystems[subsystem]; ok {
		subsystemName = val
	}

	numberOfRva, _ := readUint32le(file, offset+92)

	optHeaderMainLen := int64(96)

	optHeader := Layout{
		Offset: offset,
		Info:   "PE optional header",
		Type:   Group,
		Childs: []Layout{
			{Offset: offset, Length: 2, Info: "type = " + typeName, Type: Uint16le},
			{Offset: offset + 2, Length: 2, Info: "linker version", Type: MajorMinor16le},
			{Offset: offset + 4, Length: 4, Info: "size of code", Type: Uint32le},
			{Offset: offset + 8, Length: 4, Info: "size of initialized data", Type: Uint32le},
			{Offset: offset + 12, Length: 4, Info: "size of uninitialized data", Type: Uint32le},
			{Offset: offset + 16, Length: 4, Info: "address of entry point", Type: Uint32le},
			{Offset: offset + 20, Length: 4, Info: "base of code", Type: Uint32le},
			{Offset: offset + 24, Length: 4, Info: "base of data", Type: Uint32le},
			{Offset: offset + 28, Length: 4, Info: "image base", Type: Uint32le},
			{Offset: offset + 32, Length: 4, Info: "section alignment", Type: Uint32le},
			{Offset: offset + 36, Length: 4, Info: "file alignment", Type: Uint32le},
			{Offset: offset + 40, Length: 4, Info: "os version", Type: MajorMinor32le},
			{Offset: offset + 44, Length: 4, Info: "image version", Type: MajorMinor32le},
			{Offset: offset + 48, Length: 4, Info: "subsystem version", Type: MajorMinor32le},
			{Offset: offset + 52, Length: 4, Info: "win32 version value", Type: Uint32le},
			{Offset: offset + 56, Length: 4, Info: "size of image", Type: Uint32le},
			{Offset: offset + 60, Length: 4, Info: "size of headers", Type: Uint32le},
			{Offset: offset + 64, Length: 4, Info: "checksum", Type: Uint32le},
			{Offset: offset + 68, Length: 2, Info: "subsystem = " + subsystemName, Type: Uint16le},
			{Offset: offset + 70, Length: 2, Info: "dll characteristics", Type: Uint16le},
			{Offset: offset + 72, Length: 4, Info: "size of stack reserve", Type: Uint32le},
			{Offset: offset + 76, Length: 4, Info: "size of stack commit", Type: Uint32le},
			{Offset: offset + 80, Length: 4, Info: "size of heap reserve", Type: Uint32le},
			{Offset: offset + 84, Length: 4, Info: "size of heap commit", Type: Uint32le},
			{Offset: offset + 88, Length: 4, Info: "loader flags", Type: Uint32le},
			{Offset: offset + 92, Length: 4, Info: "number of rva and sizes", Type: Uint32le},
		}}
	offset += optHeaderMainLen

	if numberOfRva != 16 {
		panic("odd number of RVA:s = " + fmt.Sprintf("%d", numberOfRva))
	}

	ddLen := int64(8)

	optHeader.Length = optHeaderMainLen + (int64(numberOfRva) * ddLen)
	if optHeader.Length != int64(size) {
		fmt.Println("error: PE unexpected opt header len. expected ", size, " actual =", optHeader.Length)
	}

	for i := int64(0); i < int64(numberOfRva); i++ {

		info := "data directory " + fmt.Sprintf("%d", i)
		if val, ok := peRvaChunks[i]; ok {
			info = val
		}

		optHeader.Childs = append(optHeader.Childs, []Layout{
			{Offset: offset, Length: 4, Info: info + " RVA", Type: Uint32le},
			{Offset: offset + 4, Length: 4, Info: info + " size", Type: Uint32le},
		}...)
		offset += 8
	}

	return optHeader
}
