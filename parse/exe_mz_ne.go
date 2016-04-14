package parse

// 16-bit NE exe (Win16, OS/2)
// STATUS: 10%

import (
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
func parseMZ_NEHeader(file *os.File) (*Layout, error) {

	offset := int64(0x400)

	targetOSId, _ := readUint8(file, offset+54)
	targetOS := "unknown"

	if val, ok := neTargetOS[targetOSId]; ok {
		targetOS = val
	}

	res := Layout{
		Offset: offset,
		Length: 64, // XXX
		Info:   "NE header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: offset, Length: 2, Info: "identifier", Type: ASCII},
			Layout{Offset: offset + 2, Length: 2, Info: "linker version", Type: MajorMinor16},
			Layout{Offset: offset + 4, Length: 2, Info: "entry table offset", Type: Uint16le},
			Layout{Offset: offset + 6, Length: 2, Info: "entry table length", Type: Uint16le},
			Layout{Offset: offset + 8, Length: 4, Info: "file crc", Type: Uint32le},
			Layout{Offset: offset + 12, Length: 2, Info: "format flags", Type: Uint16le},
			// XXX bit map:
			/*
			   if ((FormatFlagsValue & 0x0001) != 0) {
			       // The  linker  sets  this  bit  if  the executable-file format is
			       // SINGLEDATA. An  executable file with  this format contains  one
			       // data segment.  This bit is  set if the  file is a  dynamic-link
			       // library (DLL).
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_SINGLEDATA"));
			   }
			   if ((FormatFlagsValue & 0x0002) != 0) {
			       // The  linker  sets  this  bit  if  the executable-file format is
			       // MULTIPLEDATA.  An  executable  file  with  this format contains
			       // multiple  data segments.  This bit  is  set  if the  file is  a
			       // Windows application.
			       // If neither bit  0 nor bit 1 is  set, the executable-file format
			       // is  NOAUTODATA. An  executable file  with this  format does not
			       // contain an automatic data segment.
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_MULTIPLEDATA"));
			   }
			   if ((FormatFlagsValue & 0x0010) != 0)
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_WIN32"));
			   if ((FormatFlagsValue & 0x0020) != 0)// Wine built-in module
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_BUILTIN"));
			   if ((FormatFlagsValue & 0x0800) != 0) {
			       // If this  bit is set, the  first segment in the  executable file
			       // contains code that loads the application.
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_SELFLOAD"));
			   }
			   if ((FormatFlagsValue & 0x2000) != 0) {
			       // If this bit is set, the  linker detects errors at link time but
			       // still creates an executable file.
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_LINKERROR"));
			   }
			   if ((FormatFlagsValue & 0x4000) != 0)
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_CALLWEP"));
			   if ((FormatFlagsValue & 0x8000) != 0) {
			       // If this bit is set, the executable file is a library module.
			       // If   bit  15   is  set,   the  CS:IP   registers  point  to  an
			       // initialization  procedure  called  with  the  value  in  the AX
			       // register  equal  to  the   module  handle.  The  initialization
			       // procedure  must execute  a far   return to  the caller.  If the
			       // procedure is successful, the value in AX is nonzero. Otherwise,
			       // the value in AX is zero. The value in the DS register is set to
			       // the library's data segment if  SINGLEDATA is set. Otherwise, DS
			       // is set  to the data segment  of the application that  loads the
			       // library.
			       FormatFlags.Nodes.Add(new Chunk("NE_FFLAGS_LIBMODULE"));
			   }*/

			Layout{Offset: offset + 14, Length: 2, Info: "auto data segment index", Type: Uint16le},
			Layout{Offset: offset + 16, Length: 2, Info: "initial local heap size", Type: Uint16le},
			Layout{Offset: offset + 18, Length: 2, Info: "initial stack size", Type: Uint16le},
			Layout{Offset: offset + 20, Length: 4, Info: "entry point CS:IP", Type: Uint32le},   // XXX type CS:IP,  XXX XFIXME READ PARSE DECODE
			Layout{Offset: offset + 24, Length: 4, Info: "stack pointer SS:SP", Type: Uint32le}, // XXX type
			Layout{Offset: offset + 28, Length: 2, Info: "segment table entries", Type: Uint16le},
			Layout{Offset: offset + 30, Length: 2, Info: "module reference table entires", Type: Uint16le},
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
			Layout{Offset: offset + 56, Length: 2, Info: "OffsetToFastload", Type: Uint16le}, // XXX only used by windows
			Layout{Offset: offset + 58, Length: 2, Info: "LengthOfFastload", Type: Uint16le}, // XXX only used by windows, offset to segment reference thunks or length of gangload area.
			Layout{Offset: offset + 60, Length: 2, Info: "reserved", Type: Uint16le},
			Layout{Offset: offset + 62, Length: 2, Info: "ExpectedWindowsVersion", Type: MinorMajor16}, // XXX only used by windows
		}}

	/*

	   neHead.Nodes.Add(ParseNEModuleReferenceTable(OffsetModuleReferenceTableValue, ModuleReferenceCountValue));

	   neHead.Nodes.Add(ParseNEEntryTable(OffsetToEntryTableValue, EntryTableLengthValue));

	   neHead.Nodes.Add(ParseNESegmentTable(OffsetSegmentTableValue, SegmentCountValue));

	   neHead.Nodes.Add(ParseNEImportedTable(OffsetImportedNamesTableValue));

	   neHead.Nodes.Add(ParseNEResidentTable(OffsetResidentNamesTableValue));

	   neHead.Nodes.Add(ParseNENonResidentTable(OffsetNonresidentNamesTableValue));

	   neHead.Nodes.Add(ParseNEResourceTable(OffsetResourceTableValue));
	*/

	return &res, nil
}

/*
private Chunk ParseNEModuleReferenceTable(long baseOffset, uint count)
{
    var chunk = new Chunk("Module Reference Table");
    chunk.offset = baseOffset;
    chunk.length = count * 2;

    //Log("Module Reference Table at 0x" + OffsetModuleReferenceTableValue.ToString("x4"));
    BaseStream.Position = baseOffset;

    // The module-reference table contains offsets for
    // module names stored in the imported-name table.
    for (int i = 0; i < count; i++) {
        long currOffset = BaseStream.Position;
        ushort offset = ReadUInt16();
        //Log("  module reference: " + offset.ToString("x4"));

        var sub = new Chunk("Module reference = " + offset.ToString("x4"));
        sub.offset = currOffset;
        sub.length = 2;
        chunk.Nodes.Add(sub);
    }
    return chunk;
}

private Chunk ParseNEEntryTable(long baseOffset, uint length)
{
    //Log("EntryTable at 0x" + OffsetToEntryTableValue.ToString("x4") + ", Length " + EntryTableLengthValue.ToString("x4"));
    // Log(" XXXX TODO care about MovableEntryPointsInEntryTableValue = " + MovableEntryPointsInEntryTableValue);
    var chunk = new Chunk("Entry Table");
    chunk.offset = baseOffset;
    chunk.length = length;

    BaseStream.Position = baseOffset;


    // The entry-table data is organized  by bundle,  each of which
    // begins with a 2-byte header. The first  byte of the header specifies the number
    // of entries in the bundle (a value of  00h designates the end of the table). The
    // second byte specifies whether the corresponding segment is movable or fixed. If
    // the value in  this byte is 0FFh, the  segment is movable. If the  value in this
    // byte is 0FEh, the  entry does not refer to a segment  but refers, instead, to a
    // constant defined within  the module. If the value in  this byte is neither 0FFh
    // nor 0FEh, it is a segment index.

    int entryTableLen = 0;
    do {
        var currOffset = BaseStream.Position;
        var nEntries = ReadByte();
        entryTableLen += 1;

        // Log("   entry point entries = " + nEntries);

        if (nEntries == 0) {
            // Log("   this is end of table marker");
            var sub2 = new Chunk("End Marker");
            sub2.offset = currOffset;
            sub2.length = 1;
            chunk.Nodes.Add(sub2);

            continue;
        }

        var nSegNumber = ReadByte();
        entryTableLen += 1;


        var sub = new Chunk("Header: " + nEntries + " items, segment number 0x" + nSegNumber.ToString("x2"));
        sub.offset = currOffset;
        sub.length = 2;
        chunk.Nodes.Add(sub);


        for (int i = 0; i < nEntries; i++) {
            currOffset = BaseStream.Position;
            if (nSegNumber == 0xFF) {
                byte flags = ReadByte();
                ushort int3f = ReadUInt16();
                byte segment = ReadByte();
                ushort offset = ReadUInt16();
                entryTableLen += 6;
                if (int3f != 0x3fcd) {
                    Log("PARSE ERROR in NE - entry points. int3f == " + int3f.ToString("x4"));
                    break;
                }

                //Log("[" + currOffset.ToString("x4") + "] movable segment, flags = " + flags.ToString("x2") + " , offset = " + segment.ToString("x4") + ":" + offset.ToString("x4"));

                var sub2 = new Chunk("Movable segment at " + segment.ToString("x4") + ":" + offset.ToString("x4") + ", flags = " + flags.ToString("x2"));
                sub2.offset = currOffset;
                sub2.length = 6;
                chunk.Nodes.Add(sub2);
            } else if (nSegNumber == 0xFE) {
                Log("  TODO   refer to constant defined within module");

                // struct entry_tab_fixed_s
                // unsigned char flags;
                // unsigned short offset;

            } else {
                if (nEntries > 1)
                    throw new Exception("sample please! entries = " + nEntries);
                //Log("  TODO segment index " + nSegNumber);
                //NOTE: only sample i seen was empty here
            }
        }

    } while (entryTableLen < length);
    return chunk;
}

private Chunk ParseNESegmentTable(long baseOffset, uint count)
{
    //Log("Segment Table at 0x" + OffsetSegmentTableValue.ToString("x4"));
    var chunk = new Chunk("Segment Table");
    chunk.offset = baseOffset;
    chunk.length = 8 * count;

    BaseStream.Position = baseOffset;

    for (int i = 0; i < count; i++) {
        var offset = ReadUInt16() * 16;   // Sector offset (in segments) of segment data. 0 = no data exists
        if (offset == 0)
            throw new Exception("sample plz");


        var length = ReadUInt16();   // Length of segment data
        var flags = ReadUInt16();    // Flags associated with this segment
        var minAlloc = ReadUInt16(); // Minimum allocation size for table. 0 = 64k, unless offset also is 0

        var sub = new Chunk("offset=0x" + offset.ToString("x4") + ", flags=0x" + flags.ToString("x4") + ", minAlloc=0x" + minAlloc.ToString("x4"));
        sub.offset = offset;
        sub.length = length;

        if ((flags & 0x0001) != 0) {
            // If this bit  is set, the segment is a data segment. Otherwise,
            // the segment is a code segment.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_DATA"));
        }

        if ((flags & 0x0002) != 0) {
            // If this  bit is set,  the loader has  allocated memory for  the segment.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_ALLOCATED"));
        }

        if ((flags & 0x0004) != 0) {
            // If this bit is set, the segment is loaded.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_LOADED"));
        }

        if ((flags & 0x0008) != 0) {
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_ITERATED"));
        }

        if ((flags & 0x0010) != 0) {
            // If this bit is set, the segment type is MOVABLE. Otherwise, the
            // segment type is FIXED.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_MOVEABLE"));
        }
        if ((flags & 0x0020) != 0) {
            // If  this bit  is set,  the segment  type is  PURE or SHAREABLE.
            // Otherwise, the segment type is IMPURE or NONSHAREABLE.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_SHAREABLE"));
        }
        if ((flags & 0x0040) != 0) {
            // If this bit is set, the segment type is PRELOAD. Otherwise, the
            // segment type is LOADONCALL.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_PRELOAD"));
        }
        if ((flags & 0x0080) != 0) {
            // If  this bit  is set  and the  segment is  a code  segment, the
            // segment type is EXECUTEONLY. If this bit is set and the segment
            // is a data segment, the segment type is READONLY.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_EXECUTEONLY"));
        }
        if ((flags & 0x0100) != 0) {
            // If this bit is set, the segment contains relocation data.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_RELOC_DATA"));
        }
        if ((flags & 0x1000) != 0) {
            // If this bit is set, the segment is discardable.
            sub.Nodes.Add(new Chunk("NE_SEGFLAGS_DISCARDABLE"));
        }

        chunk.Nodes.Add(sub);
    }
    return chunk;
}

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
