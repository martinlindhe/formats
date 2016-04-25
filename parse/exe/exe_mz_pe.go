package exe

// PE/COFF exe (Win32, Win64)
// https://www.microsoft.com/whdc/system/platform/firmware/PECOFF.mspx
// http://wiki.osdev.org/PE

// STATUS: 50%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
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
			{Offset: pos + 8, Length: 4, Info: "timestamp", Type: parse.DOSDateTime},
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

	res = append(res, parsePESections(file, pos, numberOfSections)...)

	return res, nil
}

func parsePESections(file *os.File, pos int64, numberOfSections uint16) []parse.Layout {

	res := []parse.Layout{}
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
				{Low: 0, Length: 3, Info: "reserved"},
				{Low: 3, Length: 1, Info: "no padding", Spec: "IMAGE_SCN_TYPE_NO_PAD"},
				{Low: 4, Length: 1, Info: "reserved"},
				{Low: 5, Length: 1, Info: "code", Spec: "IMAGE_SCN_CNT_CODE"},
				{Low: 6, Length: 1, Info: "initialized data", Spec: "IMAGE_SCN_CNT_INITIALIZED_DATA"},
				{Low: 7, Length: 1, Info: "uninitialized data", Spec: "IMAGE_SCN_CNT_UNINITIALIZED_DATA"},
				{Low: 8, Length: 1, Info: "reserved", Spec: "IMAGE_SCN_LNK_OTHER"},
				{Low: 9, Length: 1, Info: "info", Spec: "IMAGE_SCN_LNK_INFO"},
				{Low: 10, Length: 1, Info: "reserved", Spec: ""},
				{Low: 11, Length: 1, Info: "remove", Spec: "IMAGE_SCN_LNK_REMOVE"},
				{Low: 12, Length: 1, Info: "COMDAT data", Spec: "IMAGE_SCN_LNK_COMDAT"},
				{Low: 13, Length: 4, Info: "reserved"},
				{Low: 17, Length: 1, Info: "data referenced through the global pointer", Spec: "IMAGE_SCN_GPREL"},
				{Low: 18, Length: 2, Info: "reserved"},
				{Low: 20, Length: 4, Info: "align x-bytes", Spec: "IMAGE_SCN_ALIGN_xBYTES"},
				{Low: 24, Length: 1, Info: "extended relocations", Spec: "IMAGE_SCN_LNK_NRELOC_OVFL"},
				{Low: 25, Length: 1, Info: "discardable", Spec: "IMAGE_SCN_MEM_DISCARDABLE"},
				{Low: 26, Length: 1, Info: "cacheable", Spec: "IMAGE_SCN_MEM_NOT_CACHED"},
				{Low: 27, Length: 1, Info: "pageable", Spec: "IMAGE_SCN_MEM_NOT_PAGED"},
				{Low: 28, Length: 1, Info: "shared", Spec: "IMAGE_SCN_MEM_SHARED"},
				{Low: 29, Length: 1, Info: "executable", Spec: "IMAGE_SCN_MEM_EXECUTE"},
				{Low: 30, Length: 1, Info: "readable", Spec: "IMAGE_SCN_MEM_READ"},
				{Low: 31, Length: 1, Info: "writeable", Spec: "IMAGE_SCN_MEM_WRITE"},
			}}}

		sectionHeader.Childs = append(sectionHeader.Childs, chunk...)
		pos += peSectionHeaderLen
	}

	res = append(res, sectionHeader)
	return res
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
		fmt.Println("error: expected 16 RVA:s, found " + fmt.Sprintf("%d", numberOfRva))
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

/* XXX
public long FileOffsetFromVirtualAddress(long va)
{
    if (this.sections.Count == 0) {
        Log("no sections - ERROR");
        return va;
        //throw new Exception("no sections");
    }

    // Log("translate VA " + va.ToString("x8")+ " to file offset");

    foreach (var section in this.sections) {
        if (va >= section.virtualOffset && (va < section.virtualOffset + section.length)) {
            long res = (va - section.virtualOffset) + section.realOffset;
            // Log("translated to " + res.ToString("x8"));
            return res;
        }
    }
    Log("FATAL ERROR not found for va " + va.ToString("x8"));
    return va;
    //throw new Exception("not found for va " + va.ToString("x8"));
}
*/
