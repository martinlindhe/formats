package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

var (
	elfClasses = map[byte]string{
		0: "none",
		1: "32-bit",
		2: "64-bit",
	}
	elfDataEncodings = map[byte]string{
		1: "lsb",
		2: "msb",
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
		0: "none",
		1: "AT&T WE 32100",
		2: "SPARC",
		3: "intel 80386",
		4: "motorola 68000",
		5: "motorola 88000",
		7: "intel 80860",
		8: "MIPS RS3000",
	}
)

func ELF(file *os.File) (*ParsedLayout, error) {

	if !isELF(file) {
		return nil, nil
	}
	return parseELF(file)
}

func isELF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [16]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX 16 first bytes are id
	if b[0] == 0x7f && b[1] == 'E' && b[2] == 'L' && b[3] == 'F' {
		return true
	}

	return false
}

func parseELF(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)

	elfClass, _ := readUint8(file, pos+4)
	className := "?"
	if val, ok := elfClasses[elfClass]; ok {
		className = val
	}

	encoding, _ := readUint8(file, pos+5)
	encodingName := "?"
	if val, ok := elfDataEncodings[encoding]; ok {
		encodingName = val
	}

	elfType, _ := readUint16le(file, pos+16)
	typeName := "?"
	if val, ok := elfTypes[elfType]; ok {
		typeName = val
	}

	machine, _ := readUint16le(file, pos+18)
	machineName := "?"
	if val, ok := elfMachines[machine]; ok {
		machineName = val
	}

	phOffset, _ := readUint32le(file, pos+28)
	phEntrySize, _ := readUint16le(file, pos+42)
	phCount, _ := readUint16le(file, pos+44)

	shOffset, _ := readUint32le(file, pos+32)
	shEntrySize, _ := readUint16le(file, pos+48)
	shCount, _ := readUint16le(file, pos+50)

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: pos,
			Length: 52, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: Uint32le},
				{Offset: pos + 4, Length: 1, Info: "class = " + className, Type: Uint8},
				{Offset: pos + 5, Length: 1, Info: "data encoding = " + encodingName, Type: Uint8},
				{Offset: pos + 6, Length: 1, Info: "header version", Type: Uint8},
				{Offset: pos + 7, Length: 9, Info: "reserved", Type: Bytes},
				{Offset: pos + 16, Length: 2, Info: "type = " + typeName, Type: Uint16le},
				{Offset: pos + 18, Length: 2, Info: "machine = " + machineName, Type: Uint16le},
				{Offset: pos + 20, Length: 4, Info: "version", Type: Uint32le},
				{Offset: pos + 24, Length: 4, Info: "entry", Type: Uint32le},
				{Offset: pos + 28, Length: 4, Info: "program header offset", Type: Uint32le},
				{Offset: pos + 32, Length: 4, Info: "section header offset", Type: Uint32le},
				{Offset: pos + 36, Length: 4, Info: "flags", Type: Uint32le},
				{Offset: pos + 40, Length: 2, Info: "elf header size", Type: Uint16le},
				{Offset: pos + 42, Length: 2, Info: "program header entry size", Type: Uint16le},
				{Offset: pos + 44, Length: 2, Info: "program header count", Type: Uint16le},
				{Offset: pos + 46, Length: 2, Info: "section header entry size", Type: Uint16le},
				{Offset: pos + 48, Length: 2, Info: "section header count", Type: Uint16le},
				{Offset: pos + 50, Length: 2, Info: "section header strndx", Type: Uint16le}, // XXX map
			}}}}

	if phOffset > 0 {
		pos = int64(phOffset)
		phLen := int64(phEntrySize * phCount)
		res.Layout = append(res.Layout, Layout{
			Offset: pos,
			Length: phLen,
			Info:   "program header",
			Type:   Group,
			Childs: []Layout{}, // XXX childs
		})
	}

	if shOffset > 0 {
		pos = int64(shOffset)
		shLen := int64(shEntrySize * shCount)
		res.Layout = append(res.Layout, Layout{
			Offset: pos,
			Length: shLen,
			Info:   "section header",
			Type:   Group,
			Childs: []Layout{}, // XXX childs
		})
	}

	return &res, nil
}
