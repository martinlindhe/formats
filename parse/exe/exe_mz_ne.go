package exe

// 16-bit NE exe (Win16, OS/2)
// STATUS: 70%

// XXX: http://www.program-transformation.org/Transform/NeFormat

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	neTargetOS = map[byte]string{
		1: "os/2",
		2: "windows",
		3: "european ms-dos 4.x",
		4: "windows 386",
		5: "boss", // Borland Operating System Services
	}

	neResourceType = map[uint16]string{
		0x8001: "cursor",
		0x8002: "bitmap",
		0x8003: "icon",
		0x8004: "menu",
		0x8005: "dialog",
		0x8006: "string",
		0x8007: "font directory",
		0x8008: "font",
		0x8009: "accelerator table",
		0x800a: "resource data",
		0x800c: "group cursor",
		0x800e: "group icon",
		0x8010: "version",
	}
)

// parses 16-bit Windows and OS/2 executables
func parseMZ_NEHeader(file *os.File, pos int64) ([]parse.Layout, error) {

	res := []parse.Layout{}
	targetOS, _ := parse.ReadToMap(file, parse.Uint8, pos+54, neTargetOS)
	res = append(res, parse.Layout{
		Offset: pos,
		Length: 64, // XXX
		Info:   "external header = NE",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
			{Offset: pos + 2, Length: 2, Info: "linker version", Type: parse.MajorMinor16le},
			{Offset: pos + 4, Length: 2, Info: "entry table offset", Type: parse.Uint16le},
			{Offset: pos + 6, Length: 2, Info: "entry table length", Type: parse.Uint16le},
			{Offset: pos + 8, Length: 4, Info: "file load crc", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 1, Info: "program flags", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 2, Info: "dgroup type"}, // XXX 0=none, 1=single shared, 2=multiple, 3=null
				{Low: 2, Length: 1, Info: "global initialization"},
				{Low: 3, Length: 1, Info: "protected mode only"},
				{Low: 4, Length: 1, Info: "8086 instructions"},
				{Low: 5, Length: 1, Info: "80286 instructions"},
				{Low: 6, Length: 1, Info: "80386 instructions"},
				{Low: 7, Length: 1, Info: "80x87 instructions"},
			}},
			{Offset: pos + 13, Length: 1, Info: "app flags", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 3, Info: "app type"},               // XXX 1=unaware of win api, 2=compatible with win api, 3=uses win api
				{Low: 3, Length: 1, Info: "OS/2 family app"},        // XXX
				{Low: 4, Length: 1, Info: "reserved"},               // XXX
				{Low: 5, Length: 1, Info: "errors in image"},        // XXX
				{Low: 6, Length: 1, Info: "non-conforming program"}, // XXX
				{Low: 7, Length: 1, Info: "dll or driver"},          // XXX
			}},
			{Offset: pos + 14, Length: 2, Info: "auto data segment index", Type: parse.Uint16le},
			{Offset: pos + 16, Length: 2, Info: "initial local heap size", Type: parse.Uint16le},
			{Offset: pos + 18, Length: 2, Info: "initial stack size", Type: parse.Uint16le},
			{Offset: pos + 20, Length: 4, Info: "entry point CS:IP", Type: parse.Uint32le},   // XXX type CS:IP,  XXX XFIXME READ PARSE DECODE
			{Offset: pos + 24, Length: 4, Info: "stack pointer SS:SP", Type: parse.Uint32le}, // XXX type
			{Offset: pos + 28, Length: 2, Info: "segment table entries", Type: parse.Uint16le},
			{Offset: pos + 30, Length: 2, Info: "module reference entires", Type: parse.Uint16le},
			{Offset: pos + 32, Length: 2, Info: "nonresident names table size", Type: parse.Uint16le},
			{Offset: pos + 34, Length: 2, Info: "segment table offset", Type: parse.Uint16le},
			{Offset: pos + 36, Length: 2, Info: "resource table offset", Type: parse.Uint16le},
			{Offset: pos + 38, Length: 2, Info: "resident names table offset", Type: parse.Uint16le},
			{Offset: pos + 40, Length: 2, Info: "module reference table offset", Type: parse.Uint16le},
			{Offset: pos + 42, Length: 2, Info: "imported names table offset", Type: parse.Uint16le},
			{Offset: pos + 44, Length: 4, Info: "nonresident names table offset", Type: parse.Uint32le},
			{Offset: pos + 48, Length: 2, Info: "movable entry points in entry table", Type: parse.Uint16le},
			{Offset: pos + 50, Length: 2, Info: "file alignment size shift", Type: parse.Uint16le}, //  File alignment size shift count, 0 is equivalent to 9 (default 512-byte pages)
			{Offset: pos + 52, Length: 2, Info: "resource table entries", Type: parse.Uint16le},
			{Offset: pos + 54, Length: 1, Info: "target os = " + targetOS, Type: parse.Uint8},
			{Offset: pos + 55, Length: 1, Info: "extra flags", Type: parse.Uint8, Masks: []parse.Mask{
				{Low: 0, Length: 1, Info: "long filename support"},
				{Low: 1, Length: 1, Info: "win2 protected mode"},
				{Low: 2, Length: 1, Info: "win2 proportional fonts"},
				{Low: 3, Length: 1, Info: "fastload area"},
				{Low: 4, Length: 4, Info: "reserved"},
			}},
			{Offset: pos + 56, Length: 2, Info: "offset to fastload", Type: parse.Uint16le}, // XXX only used by windows
			{Offset: pos + 58, Length: 2, Info: "length of fastload", Type: parse.Uint16le}, // XXX only used by windows, offset to segment reference thunks or length of gangload area.
			{Offset: pos + 60, Length: 2, Info: "reserved", Type: parse.Uint16le},
			{Offset: pos + 62, Length: 2, Info: "expected windows version", Type: parse.MinorMajor16le}, // XXX only used by windows
		}})

	moduleReferenceEntries, _ := parse.ReadUint16le(file, pos+30)
	moduleReferenceOffset, _ := parse.ReadUint16le(file, pos+40)
	if moduleReferenceEntries > 0 {
		res = append(res, *parseNEModuleReferenceTable(pos+int64(moduleReferenceOffset), moduleReferenceEntries))
	}

	entryTableOffset, _ := parse.ReadUint16le(file, pos+4)
	entryTableLength, _ := parse.ReadUint16le(file, pos+6)
	res = append(res, *parseNEEntryTable(file, pos+int64(entryTableOffset), entryTableLength))

	segmentTableOffset, _ := parse.ReadUint16le(file, pos+34)
	segmentTableEntries, _ := parse.ReadUint16le(file, pos+28)
	if segmentTableEntries > 0 {
		res = append(res, *parseNESegmentTable(file, pos+int64(segmentTableOffset), segmentTableEntries))
	}

	importedNamesTableOffset, _ := parse.ReadUint16le(file, pos+42)
	res = append(res, *parseNEImportedTable(file, pos+int64(importedNamesTableOffset)))

	residentNamesTableOffset, _ := parse.ReadUint16le(file, pos+38)
	res = append(res, *parseNEResidentTable(file, pos+int64(residentNamesTableOffset)))

	nonResidentNamesTableOffset, _ := parse.ReadUint32le(file, pos+44)
	nonresidentNamesTableSize, _ := parse.ReadUint16le(file, pos+32)
	res = append(res, *parseNENonResidentTable(file, int64(nonResidentNamesTableOffset), nonresidentNamesTableSize))

	resourceTableOffset, _ := parse.ReadUint16le(file, pos+36)
	resourceTableEntries, _ := parse.ReadUint16le(file, pos+52)
	res = append(res, *parseNEResourceTable(file, pos+int64(resourceTableOffset), resourceTableEntries))

	fastloadAreaOffset, _ := parse.ReadUint16le(file, pos+56)
	fastloadAreaLength, _ := parse.ReadUint16le(file, pos+58)

	if fastloadAreaLength > 0 {
		// XXX offset seems wrong
		res = append(res, parse.Layout{
			Offset: int64(fastloadAreaOffset) * 16,
			Length: int64(fastloadAreaLength),
			Info:   "fast-load area", // XXX
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: int64(fastloadAreaOffset) * 16, Length: int64(fastloadAreaLength), Info: "fast-load data", Type: parse.Bytes},
			}})
	}

	return res, nil
}

func parseNEModuleReferenceTable(pos int64, count uint16) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Length: int64(count) * 2,
		Info:   "NE module reference table",
		Type:   parse.Group}

	for i := uint16(1); i <= count; i++ {
		res.Childs = append(res.Childs, parse.Layout{
			Offset: pos,
			Length: 2,
			Info:   "module reference " + fmt.Sprintf("%d", i),
			Type:   parse.Uint16le})
		pos += 2
	}
	return &res
}

func parseNEEntryTable(file *os.File, pos int64, length uint16) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		// Length: int64(length),
		Info: "NE entry table",
		Type: parse.Group}

	// The entry-table data is organized by bundle, each of which begins with
	// a 2-byte header. The first byte of the header specifies the number of
	// entries in the bundle (a value of 00h designates the end of the table).
	// The second byte specifies whether the corresponding segment is movable
	// or fixed. If the value in this byte is 0FFh, the segment is movable.
	// If the value in this byte is 0FEh, the entry does not refer to a segment
	// but refers, instead, to a constant defined within the module. If the
	// value in this byte is neither 0FFh nor 0FEh, it is a segment index.

	for {

		items, _ := parse.ReadUint8(file, pos)
		segNumber, _ := parse.ReadUint8(file, pos+1)

		if items == 0 {
			// NOTE: tagging the empty "items" block as end marker
			res.Childs = append(res.Childs, parse.Layout{
				Offset: pos,
				Length: 1,
				Info:   "end marker",
				Type:   parse.Uint16le})
			pos += 1
			break
		}

		res.Childs = append(res.Childs, []parse.Layout{
			{Offset: pos, Length: 1, Info: "items", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: "segment", Type: parse.Uint8},
		}...)
		pos += 2

		for i := 1; i <= int(items); i++ {

			switch segNumber {
			case 0xff:

				id := fmt.Sprintf("%d", i)
				res.Childs = append(res.Childs, []parse.Layout{
					{Offset: pos, Length: 1, Info: "movable " + id + " flags", Type: parse.Uint8, Masks: []parse.Mask{
						{Low: 0, Length: 1, Info: "exported"},
						{Low: 1, Length: 1, Info: "global data segment"},
						{Low: 2, Length: 1, Info: "reserved"},
						{Low: 3, Length: 5, Info: "ring transition words"},
					}},
					{Offset: pos + 1, Length: 2, Info: "movable " + id + " int3f", Type: parse.Uint16le},
					{Offset: pos + 3, Length: 1, Info: "movable " + id + " segment", Type: parse.Uint8},
					{Offset: pos + 4, Length: 2, Info: "movable " + id + " offset", Type: parse.Uint16le},
				}...)
				pos += 6

				/*			case 0xfe:
							// panic("  TODO   refer to constant defined within module")
							// struct entry_tab_fixed_s
							// unsigned char flags;
							// unsigned short offset;
				*/
			default:
				fmt.Println("  TODO segment index ", segNumber, ", entries", items)
				//NOTE: only sample i seen was empty here
				// panic("xxx")
			}
		}
	}

	res.Length = int64(length)
	return &res
}

func parseNESegmentTable(file *os.File, pos int64, count uint16) *parse.Layout {

	segmentLen := int64(8)

	res := parse.Layout{
		Offset: pos,
		Length: segmentLen * int64(count),
		Info:   "NE segment table",
		Type:   parse.Group}

	for i := 1; i <= int(count); i++ {
		id := fmt.Sprintf("%d", i)

		res.Childs = append(res.Childs, []parse.Layout{
			{Offset: pos, Length: 2, Info: "segment " + id + " offset", Type: parse.Uint16le}, // in segments. 0 = no data exists
			{Offset: pos + 2, Length: 2, Info: "segment " + id + " length", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 2, Info: "segment " + id + " flags", Type: parse.Uint16le, Masks: []parse.Mask{
				{Low: 0, Length: 1, Info: "segment " + id + " type"}, // 0=code, 1=data
				{Low: 1, Length: 1, Info: "allocated"},
				{Low: 2, Length: 1, Info: "loaded"},
				{Low: 3, Length: 1, Info: "iterated"},
				{Low: 4, Length: 1, Info: "1=moveable, 0=fixed"},
				{Low: 5, Length: 1, Info: "shareable"},
				{Low: 6, Length: 1, Info: "1=preload, 0=loadoncall"},
				{Low: 7, Length: 1, Info: "execute only/read only"},
				{Low: 8, Length: 1, Info: "reloc data"},
				{Low: 9, Length: 3, Info: "reserved"},
				{Low: 12, Length: 1, Info: "discardable"},
				{Low: 13, Length: 3, Info: "reserved"},
			}},
			{Offset: pos + 6, Length: 2, Info: "segment " + id + " min alloc size", Type: parse.Uint16le}, // 0 = 64k
		}...)
		pos += segmentLen
	}

	return &res
}

func parseNEImportedTable(file *os.File, pos int64) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Info:   "NE imported names table",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "reserved", Type: parse.Uint8}, // XXX ?
		}}

	pos++

	var len byte

	totLen := int64(1)
	for {

		len, _ = parse.ReadUint8(file, pos)

		b := parse.ReadBytesFrom(file, pos+1, int64(len))

		subLen := int64(len) + 1
		info := string(b)
		subType := parse.ASCIIC
		brk := false
		if len <= 1 || b[0] == 0 || b[0] == 0xff {
			info = "end marker"
			subLen = 2
			subType = parse.Uint16le
			brk = true
		}

		res.Childs = append(res.Childs, parse.Layout{
			Offset: pos,
			Length: subLen,
			Info:   info,
			Type:   subType,
		})

		pos += subLen
		totLen += subLen

		if brk {
			break
		}
	}

	res.Length = totLen

	return &res
}

func parseNEResidentTable(file *os.File, pos int64) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Info:   "NE resident names table",
		Type:   parse.Group}

	residentLen := int64(0)
	var len byte
	for {

		len, _ = parse.ReadUint8(file, pos)
		chunkLen := 1 + int64(len)

		if len == 0 {
			res.Childs = append(res.Childs,
				parse.Layout{Offset: pos, Length: 1, Info: "end marker", Type: parse.Uint8})
		} else {
			res.Childs = append(res.Childs, []parse.Layout{
				{Offset: pos, Length: 1 + int64(len), Info: "data", Type: parse.ASCIIC},
				{Offset: pos + 1 + int64(len), Length: 2, Info: "ord", Type: parse.Uint16le}, // XXX ordinal value
			}...)
			chunkLen += 2
		}

		pos += chunkLen
		residentLen += chunkLen

		if len == 0 {
			break
		}
	}

	res.Length = residentLen
	return &res
}

func parseNENonResidentTable(file *os.File, pos int64, size uint16) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Info:   "NE nonresident names table",
		Type:   parse.Group}

	nonresidentLen := int64(0)

	var len byte
	for {
		len, _ = parse.ReadUint8(file, pos)
		if len == 0 {
			res.Childs = append(res.Childs,
				parse.Layout{Offset: pos, Length: 1, Info: "end marker", Type: parse.Uint8})
			nonresidentLen += 1
		} else {
			res.Childs = append(res.Childs, []parse.Layout{
				{Offset: pos, Length: 1 + int64(len), Info: "name", Type: parse.ASCIIC},
				{Offset: pos + 1 + int64(len), Length: 2, Info: "ord", Type: parse.Uint16le},
			}...)
		}
		if len == 0 {
			break
		}
		nonresidentLen += 1 + int64(len) + 2
		pos += 1 + int64(len) + 2
	}
	if int64(size) != nonresidentLen {
		fmt.Println("warning: NE nonresident table length expected ", int64(size), ", found", nonresidentLen)
	}
	res.Length = nonresidentLen
	return &res
}

func parseNEResourceTable(file *os.File, pos int64, count uint16) *parse.Layout {

	res := parse.Layout{
		Offset: pos,
		Info:   "NE resource table",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "shift", Type: parse.Uint16le},
		}}

	len := int64(2)
	pos += 2
	tnameInfoLen := int64(12)

	for {

		resourceType, _ := parse.ReadUint16le(file, pos)
		if resourceType == 0 {
			res.Childs = append(res.Childs, parse.Layout{
				Offset: pos,
				Length: 2,
				Info:   "end marker",
				Type:   parse.Uint16le})
			len += 2
			break
		}

		resourceCount, _ := parse.ReadUint16le(file, pos+2)

		info := "type"
		if val, ok := neResourceType[resourceType]; ok {
			info += " = " + val
		}

		res.Childs = append(res.Childs, []parse.Layout{ // TTYPEINFO
			{Offset: pos, Length: 2, Info: info, Type: parse.Uint16le},
			{Offset: pos + 2, Length: 2, Info: "resource count", Type: parse.Uint16le},
			{Offset: pos + 4, Length: 4, Info: "reserved", Type: parse.Uint32le},
		}...)

		pos += 8
		len += 8

		for i := 0; i < int(resourceCount); i++ {
			res.Childs = append(res.Childs, []parse.Layout{ // TNAMEINFO
				{Offset: pos, Length: 2, Info: "offset", Type: parse.Uint16le},
				{Offset: pos + 2, Length: 2, Info: "size", Type: parse.Uint16le},
				{Offset: pos + 4, Length: 2, Info: "flags", Type: parse.Uint16le},
				{Offset: pos + 6, Length: 2, Info: "id", Type: parse.Uint16le},
				{Offset: pos + 8, Length: 2, Info: "reserved 1", Type: parse.Uint16le},
				{Offset: pos + 10, Length: 2, Info: "reserved 2", Type: parse.Uint16le},
			}...)

			pos += tnameInfoLen
			len += tnameInfoLen
		}
	}

	res.Length = len
	return &res
}
