package parse

// 16-bit NE exe (Win16, OS/2)
// STATUS: 10%

// XXX: http://www.program-transformation.org/Transform/NeFormat

import (
	"fmt"
	"os"
)

var (
	neTargetOS = map[byte]string{
		1: "OS/2",
		2: "Windows",
		3: "European MS-DOS 4.x",
		4: "Windows 386",
		5: "BOSS (Borland Operating System Services)",
	}
)

// parses 16-bit Windows and OS/2 executables
func parseMZ_NEHeader(file *os.File) ([]Layout, error) {

	offset := int64(0x400)

	targetOSId, _ := readUint8(file, offset+54)
	targetOS := "unknown"

	if val, ok := neTargetOS[targetOSId]; ok {
		targetOS = val
	}

	res := []Layout{}

	res = append(res, Layout{
		Offset: offset,
		Length: 64, // XXX
		Info:   "NE header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: offset, Length: 2, Info: "identifier", Type: ASCII},
			Layout{Offset: offset + 2, Length: 2, Info: "linker version", Type: MajorMinor16},
			Layout{Offset: offset + 4, Length: 2, Info: "entry table offset", Type: Uint16le},
			Layout{Offset: offset + 6, Length: 2, Info: "entry table length", Type: Uint16le},
			Layout{Offset: offset + 8, Length: 4, Info: "file load crc", Type: Uint32le},
			Layout{Offset: offset + 12, Length: 1, Info: "program flags", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 2, Info: "dgroup type"}, // XXX 0=none, 1=single shared, 2=multiple, 3=null
				Mask{Low: 2, Length: 1, Info: "global initialization"},
				Mask{Low: 3, Length: 1, Info: "protected mode only"},
				Mask{Low: 4, Length: 1, Info: "8086 instructions"},
				Mask{Low: 5, Length: 1, Info: "80286 instructions"},
				Mask{Low: 6, Length: 1, Info: "80386 instructions"},
				Mask{Low: 7, Length: 1, Info: "80x87 instructions"},
			}},
			Layout{Offset: offset + 13, Length: 1, Info: "app flags", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 3, Info: "app type"},               // XXX 1=unaware of win api, 2=compatible with win api, 3=uses win api
				Mask{Low: 3, Length: 1, Info: "OS/2 family app"},        // XXX
				Mask{Low: 4, Length: 1, Info: "reserved"},               // XXX
				Mask{Low: 5, Length: 1, Info: "errors in image"},        // XXX
				Mask{Low: 6, Length: 1, Info: "non-conforming program"}, // XXX
				Mask{Low: 7, Length: 1, Info: "dll or driver"},          // XXX
			}},
			Layout{Offset: offset + 14, Length: 2, Info: "auto data segment index", Type: Uint16le},
			Layout{Offset: offset + 16, Length: 2, Info: "initial local heap size", Type: Uint16le},
			Layout{Offset: offset + 18, Length: 2, Info: "initial stack size", Type: Uint16le},
			Layout{Offset: offset + 20, Length: 4, Info: "entry point CS:IP", Type: Uint32le},   // XXX type CS:IP,  XXX XFIXME READ PARSE DECODE
			Layout{Offset: offset + 24, Length: 4, Info: "stack pointer SS:SP", Type: Uint32le}, // XXX type
			Layout{Offset: offset + 28, Length: 2, Info: "segment table entries", Type: Uint16le},
			Layout{Offset: offset + 30, Length: 2, Info: "module reference entires", Type: Uint16le},
			Layout{Offset: offset + 32, Length: 2, Info: "nonresident names table size", Type: Uint16le},
			Layout{Offset: offset + 34, Length: 2, Info: "offset segment table", Type: Uint16le},
			Layout{Offset: offset + 36, Length: 2, Info: "offset resource table", Type: Uint16le},
			Layout{Offset: offset + 38, Length: 2, Info: "offset resident names table", Type: Uint16le},
			Layout{Offset: offset + 40, Length: 2, Info: "offset module reference table", Type: Uint16le},
			Layout{Offset: offset + 42, Length: 2, Info: "offset imported names table", Type: Uint16le},    // XXX (array of counted strings, terminated with a string of length 00h)
			Layout{Offset: offset + 44, Length: 4, Info: "offset nonresident names table", Type: Uint32le}, // Offset from start of file to nonresident names table
			Layout{Offset: offset + 48, Length: 2, Info: "movable entry points in entry table", Type: Uint16le},
			Layout{Offset: offset + 50, Length: 2, Info: "file alignment size shift", Type: Uint16le}, //  File alignment size shift count, 0 is equivalent to 9 (default 512-byte pages)
			Layout{Offset: offset + 52, Length: 2, Info: "resources", Type: Uint16le},                 // Number of resource table entries
			Layout{Offset: offset + 54, Length: 1, Info: "target os = " + targetOS, Type: Uint8},
			Layout{Offset: offset + 55, Length: 1, Info: "extra flags", Type: Uint8, Masks: []Mask{
				Mask{Low: 0, Length: 1, Info: "long filename support"},
				Mask{Low: 1, Length: 1, Info: "win2 protected mode"},
				Mask{Low: 2, Length: 1, Info: "win2 proportional fonts"},
				Mask{Low: 3, Length: 1, Info: "fastload area"},
				Mask{Low: 4, Length: 4, Info: "reserved"},
			}},
			Layout{Offset: offset + 56, Length: 2, Info: "offset to fastload", Type: Uint16le}, // XXX only used by windows
			Layout{Offset: offset + 58, Length: 2, Info: "length of fastload", Type: Uint16le}, // XXX only used by windows, offset to segment reference thunks or length of gangload area.
			Layout{Offset: offset + 60, Length: 2, Info: "reserved", Type: Uint16le},
			Layout{Offset: offset + 62, Length: 2, Info: "expected windows version", Type: MinorMajor16}, // XXX only used by windows
		}})

	moduleReferenceEntries, _ := readUint16le(file, offset+30)
	moduleReferenceOffset, _ := readUint16le(file, offset+40)
	res = append(res, *parseNEModuleReferenceTable(offset+int64(moduleReferenceOffset), moduleReferenceEntries))

	entryTableOffset, _ := readUint16le(file, offset+4)
	entryTableLength, _ := readUint16le(file, offset+6)
	res = append(res, *parseNEEntryTable(file, offset+int64(entryTableOffset), entryTableLength))

	segmentTableOffset, _ := readUint16le(file, offset+34)
	segmentTableEntries, _ := readUint16le(file, offset+28)
	res = append(res, *parseNESegmentTable(file, offset+int64(segmentTableOffset), segmentTableEntries))
	/*


	   neHead.Nodes.Add(ParseNEImportedTable(OffsetImportedNamesTableValue));

	   neHead.Nodes.Add(ParseNEResidentTable(OffsetResidentNamesTableValue));

	   neHead.Nodes.Add(ParseNENonResidentTable(OffsetNonresidentNamesTableValue));

	   neHead.Nodes.Add(ParseNEResourceTable(OffsetResourceTableValue));
	*/

	return res, nil
}

func parseNEModuleReferenceTable(offset int64, count uint16) *Layout {

	res := Layout{
		Offset: offset,
		Length: int64(count) * 2,
		Info:   "NE module reference table",
		Type:   Group}

	// The module-reference table contains offsets for
	// module names stored in the imported-name table.
	for i := uint16(1); i <= count; i++ {
		res.Childs = append(res.Childs, Layout{Offset: offset, Length: 2, Info: "module reference " + fmt.Sprintf("%d", i), Type: Uint16le})
		offset += 2
	}
	return &res
}

func parseNEEntryTable(file *os.File, offset int64, length uint16) *Layout {

	res := Layout{
		Offset: offset,
		Length: int64(length),
		Info:   "NE entry table",
		Type:   Group}

	// The entry-table data is organized by bundle, each of which begins with
	// a 2-byte header. The first byte of the header specifies the number of
	// entries in the bundle (a value of 00h designates the end of the table).
	// The second byte specifies whether the corresponding segment is movable
	// or fixed. If the value in this byte is 0FFh, the segment is movable.
	// If the value in this byte is 0FEh, the entry does not refer to a segment
	// but refers, instead, to a constant defined within the module. If the
	// value in this byte is neither 0FFh nor 0FEh, it is a segment index.

	entryTableLen := 0
	for entryTableLen < int(length) {

		entries, _ := readUint8(file, offset)
		entryTableLen += 1

		segNumber, _ := readUint8(file, offset+1)
		entryTableLen += 1

		if entries == 0 {
			res.Childs = append(res.Childs, Layout{
				Offset: offset,
				Length: 2,
				Info:   "end marker",
				Type:   Uint16le})
			continue
		}

		res.Childs = append(res.Childs, []Layout{
			Layout{Offset: offset, Length: 1, Info: "items", Type: Uint8},
			Layout{Offset: offset + 1, Length: 1, Info: "segment", Type: Uint8},
		}...)

		offset += 2

		for i := 1; i <= int(entries); i++ {
			switch segNumber {
			case 0xff:
				int3f, _ := readUint16le(file, offset+1)

				entryTableLen += 6
				if int3f != 0x3fcd {
					panic("PARSE ERROR in NE - entry points. int3f == " + fmt.Sprintf("%04x", int3f))
				}

				id := fmt.Sprintf("%d", i)
				res.Childs = append(res.Childs, []Layout{
					Layout{Offset: offset, Length: 1, Info: "movable " + id + " flags", Type: Uint8, Masks: []Mask{
						Mask{Low: 0, Length: 1, Info: "exported"},
						Mask{Low: 1, Length: 1, Info: "global data segment"},
						Mask{Low: 2, Length: 1, Info: "reserved"},
						Mask{Low: 3, Length: 5, Info: "ring transition words"},
					}},
					Layout{Offset: offset + 1, Length: 2, Info: "movable " + id + " int3f", Type: Uint16le},
					Layout{Offset: offset + 3, Length: 1, Info: "movable " + id + " segment", Type: Uint8},
					Layout{Offset: offset + 4, Length: 2, Info: "movable " + id + " offset", Type: Uint16le},
				}...)
				offset += 6

			case 0xfe:
				panic("  TODO   refer to constant defined within module")

				// struct entry_tab_fixed_s
				// unsigned char flags;
				// unsigned short offset;

			default:
				if entries > 1 {
					// panic("sample please! entries > 1")
				}
				//Log("  TODO segment index " + nSegNumber);
				//NOTE: only sample i seen was empty here
				panic("xxx")
			}
		}
	}

	return &res
}

func parseNESegmentTable(file *os.File, offset int64, count uint16) *Layout {

	segmentLen := int64(8)

	res := Layout{
		Offset: offset,
		Length: segmentLen * int64(count),
		Info:   "NE segment table",
		Type:   Group}

	for i := 1; i <= int(count); i++ {
		id := fmt.Sprintf("%d", i)

		res.Childs = append(res.Childs, []Layout{
			Layout{Offset: offset, Length: 2, Info: "segment " + id + " offset", Type: Uint16le}, // in segments. 0 = no data exists
			Layout{Offset: offset + 2, Length: 2, Info: "segment " + id + " length", Type: Uint16le},
			Layout{Offset: offset + 4, Length: 2, Info: "segment " + id + " flags", Type: Uint16le, Masks: []Mask{
				Mask{Low: 0, Length: 1, Info: "segment " + id + " type"}, // 0=code, 1=data
				Mask{Low: 1, Length: 1, Info: "allocated"},
				Mask{Low: 2, Length: 1, Info: "loaded"},
				Mask{Low: 3, Length: 1, Info: "iterated"},
				Mask{Low: 4, Length: 1, Info: "1=moveable, 0=fixed"},
				Mask{Low: 5, Length: 1, Info: "shareable"},
				Mask{Low: 6, Length: 1, Info: "1=preload, 0=loadoncall"},
				Mask{Low: 7, Length: 1, Info: "execute only/read only"},
				Mask{Low: 8, Length: 1, Info: "reloc data"},
				Mask{Low: 9, Length: 3, Info: "reserved"},
				Mask{Low: 12, Length: 1, Info: "discardable"},
				Mask{Low: 13, Length: 3, Info: "reserved"},
			}},
			Layout{Offset: offset + 6, Length: 2, Info: "segment " + id + " min alloc size", Type: Uint16le}, // 0 = 64k
		}...)
		offset += segmentLen
	}

	return &res
}

/*

private Chunk ParseNEImportedTable(long baseOffset)
{
    var chunk = new Chunk("Imported Names Table");
    chunk.offset = baseOffset;

    //Log("Imported Names Table at 0x" + OffsetImportedNamesTableValue.ToString("x4"));
    BaseStream.Position = baseOffset;

    var unknown = ReadByte(); // FIXME reserved??
    if (unknown != 0)
        throw new Exception("Sample plz");

    byte len;
    uint importLen = 1; // first unknown byte
    do {
        long currOffset = BaseStream.Position;
        len = ReadByte();

        byte[] data = ReadBytes(len);

        var yo = new Chunk();
        yo.offset = currOffset;
        yo.length = (uint)(len + 1);

        if (len == 1 && (data[0] == 0 || data[0] == 0xFF)) {
            yo.Text = "End Marker";
        } else {
            string importName = ByteArrayToString(data);
            //Log(currOffset.ToString("x6") + ": import of len " + len + ": " + importName);
            yo.Text = importName;
        }
        chunk.Nodes.Add(yo);

        importLen += yo.length;

    } while (len > 1);

    chunk.length = importLen;
    return chunk;
}

private Chunk ParseNEResidentTable(long baseOffset)
{
    var chunk = new Chunk("Resident Names Table");
    chunk.offset = baseOffset;

    //Log("Resident Names Table at 0x" + OffsetResidentNamesTableValue.ToString("x4"));

    BaseStream.Position = baseOffset;
    //format: [byte lenght, string name, word ord]

    uint residentLen = 0;
    byte len;
    do {
        long currOffset = BaseStream.Position;

        len = ReadByte();

        var yo = new Chunk();
        yo.offset = currOffset;
        yo.length = (uint)(1 + len);

        if (len == 0) {
            yo.Text = "End Marker";
        } else {
            yo.length += 2;
            byte[] data = ReadBytes(len);

            string name = ByteArrayToString(data);
            short ord = ReadInt16();

            // Log(currOffset.ToString("x6") + ": import of len " + len + ", ord " + ord.ToString("x4") + ": " + name);
            yo.Text = name + " (ord " + ord.ToString("x4") + ")";
        }
        residentLen += yo.length;
        chunk.Nodes.Add(yo);
    } while (len > 0);

    chunk.length = residentLen;

    return chunk;
}

private Chunk ParseNENonResidentTable(long baseOffset)
{
    var chunk = new Chunk("Nonresident Names Table");
    chunk.offset = baseOffset;

    // Log("Nonresident Names Table at 0x" + OffsetNonresidentNamesTableValue.ToString("x4"));

    uint nonresidentLen = 0;
    BaseStream.Position = baseOffset;
    //format: [byte lenght, string name, word ord]

    byte len;
    do {
        long currOffset = BaseStream.Position;
        len = ReadByte();

        var yo = new Chunk();
        yo.offset = currOffset;
        yo.length = (uint)(1 + len + 2);


        if (len == 0) {
            yo.Text = "End Marker";
        } else {

            byte[] data = ReadBytes(len);

            string name = ByteArrayToString(data);
            short ord = ReadInt16();

            //Log(currOffset.ToString("x6") + ": import of len " + len + ", ord " + ord.ToString("x4") + ": " + xx);
            yo.Text = name + " (ord " + ord.ToString("x4") + ")";
        }
        nonresidentLen += yo.length;
        chunk.Nodes.Add(yo);

    } while (len > 0);

    chunk.length = nonresidentLen;

    return chunk;
}

private Chunk ParseNEResourceTable(long baseOffset)
{
    //Log("Resource Table at 0x" + OffsetResourceTableValue.ToString("x4"));
    var chunk = new Chunk("Resource Table");
    chunk.offset = baseOffset;

    BaseStream.Position = baseOffset;
    ushort shift = ReadUInt16();

    uint resourceLen = 2; // shift len

    do {
        long currOffset = BaseStream.Position;
        ushort type = ReadUInt16();

        string typeName = "UNKNOWN 0x" + type.ToString("x4");

        var yo = new Chunk();
        yo.offset = currOffset;
        yo.length = 2;

        if (type == 0) {
            yo.Text = "End Marker";

            resourceLen += yo.length;
            chunk.Nodes.Add(yo);
            break;
        }

        ushort count = ReadUInt16();
        yo.length += 2;

        switch (type) {
        case 0x8001:
            typeName = "Cursor";
            break;
        case 0x8002:
            typeName = "Bitmap";
            break;
        case 0x8003:
            typeName = "Icon";
            break;
        case 0x8004:
            typeName = "Menu";
            break;
        case 0x8005:
            typeName = "Dialog box";
            break;
        case 0x8006:
            typeName = "String table";
            break;
        case 0x8007:
            typeName = "Font directory";
            break;
        case 0x8008:
            typeName = "Font component";
            break;
        case 0x8009:
            typeName = "Accelerator table";
            break;
        case 0x800A:
            typeName = "Resource data";
            break;
        case 0x800C:
            typeName = "Cursor directory";
            break;
        case 0x800E:
                //tells wich icon to use for 16 colors and wich for 256 colors
            typeName = "Icon directory";
            break;
        case 0x8010:
            typeName = "Version";
            break;
        }
        yo.Text = typeName;

        //Log("Resource: " + typeName + ", " + count + " items");

        // skip unknown bytes (reserved?) bytes
        var r1 = ReadInt16();
        var r2 = ReadInt16();
        yo.length += 4;
        if (r1 != 0 || r2 != 0)
            throw new Exception("TODO sample-please: reserved assumed to be zero but wasnt #1: " + r1 + ", " + r2);


        for (int i = 0; i < count; i++) {
            var offset = ReadUInt16() << shift;
            var size = (uint)(ReadUInt16() << shift);
            var flags = ReadUInt16();
            var resource = ReadUInt16();

            //Log("   resource " + resource.ToString("x4") + ", offset=" + offset.ToString("x8") + ", size=" + size.ToString("x8") + ", flags=" + flags);
            var yoSub = new Chunk("resource " + resource.ToString("x4") + ", flags=" + flags.ToString("x4"));
            yoSub.offset = offset;
            yoSub.length = size;
            yo.Nodes.Add(yoSub);

            var res1 = ReadUInt16(); // skip 2 unknown bytes, 00
            var res2 = ReadUInt16(); // skip 2 more unknown bytes, 00
            if (res1 != 0 || res2 != 0)
                throw new Exception("TODO sample-please: reserved assumed to be zero wasnt #2: " + res1 + ", " + res2);

            yo.length += 12;
        }

        resourceLen += yo.length;
        chunk.Nodes.Add(yo);


    } while (true);

    chunk.length = resourceLen;

    return chunk;
}
*/
