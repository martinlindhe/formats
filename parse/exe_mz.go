package parse

// MS-DOS executable
// .exe; .sys; .dll; .ocx; .vxd; .cpl; .msi
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func MZ(file *os.File) (*ParsedLayout, error) {

	if !isMZ(file) {
		return nil, nil
	}
	return parseMZ(file)
}

func isMZ(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'M' || b[1] != 'Z' {
		return false
	}

	return true
}

func parseMZ(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 28, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: ASCII},
			Layout{Offset: 2, Length: 2, Info: "bytes in last page", Type: Uint16le},
			Layout{Offset: 4, Length: 2, Info: "pages", Type: Uint16le},
			Layout{Offset: 6, Length: 2, Info: "relocation items", Type: Uint16le},
			Layout{Offset: 8, Length: 2, Info: "header size in paragraphs", Type: Uint16le}, // 1 paragraph = group of 16 bytes
			Layout{Offset: 10, Length: 2, Info: "min mem", Type: Uint16le},
			Layout{Offset: 12, Length: 2, Info: "max mem", Type: Uint16le},
			Layout{Offset: 14, Length: 2, Info: "ss", Type: Uint16le},
			Layout{Offset: 16, Length: 2, Info: "sp", Type: Uint16le},
			Layout{Offset: 18, Length: 2, Info: "checksum", Type: Uint16le},
			Layout{Offset: 20, Length: 2, Info: "ip", Type: Uint16le},
			Layout{Offset: 22, Length: 2, Info: "cs", Type: Uint16le},

			// Offset of relocation table; 40h for new-(NE,LE,LX,W3,PE etc.) executable
			Layout{Offset: 24, Length: 2, Info: "reloc offset", Type: Uint16le},
			Layout{Offset: 26, Length: 2, Info: "overlay", Type: Uint16le},
		}})
	return &res, nil

	//	header.length = headerSizeValue * 16;

}

/*

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized()) {
        return new List<Chunk>();
    }

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.Text = "EXE header";

    var identifier = new LittleEndian16BitChunk("MZ identifier");
    identifier.offset = 0;

    header.Nodes.Add(identifier);

    var bytes = identifier.RelativeToLittleEndian16("Bytes in last page");
    header.Nodes.Add(bytes);

    var pages = bytes.RelativeToLittleEndian16("Pages");
    header.Nodes.Add(pages);

    var relocCnt = pages.RelativeToLittleEndian16("Relocation items");
    BaseStream.Position = relocCnt.offset;
    var relocCntValue = (uint)ReadInt16();
    header.Nodes.Add(relocCnt);

    /// 1 paragraph = group of 16 bytes
    var headerSize = relocCnt.RelativeToLittleEndian16("Header size in paragraphs");
    BaseStream.Position = headerSize.offset;
    var headerSizeValue = (uint)ReadInt16();
    header.Nodes.Add(headerSize);

    header.length = headerSizeValue * 16;
    this.ExeHeaderLength = header.length;

    var minPara = headerSize.RelativeToLittleEndian16("Min mem");
    header.Nodes.Add(minPara);

    var maxPara = minPara.RelativeToLittleEndian16("Max mem");
    header.Nodes.Add(maxPara);

    var ss = maxPara.RelativeToLittleEndian16("SS");
    header.Nodes.Add(ss);

    var sp = ss.RelativeToLittleEndian16("SP");
    header.Nodes.Add(sp);

    var checksum = sp.RelativeToLittleEndian16("Checksum");
    BaseStream.Position = checksum.offset;
    var mzChecksumValue = ReadUInt16();

    //Log("MZ checksum: 0x" + mzChecksumValue.ToString("x4"));
    //Log("calculated checksum: 0x" + CalculateChecksum16bit().ToString("x4"));

    header.Nodes.Add(checksum);

    var ip = checksum.RelativeToLittleEndian16("IP");
    BaseStream.Position = ip.offset;
    var ipValue = ReadUInt16();
    header.Nodes.Add(ip);

    var cs = ip.RelativeToLittleEndian16("CS");
    BaseStream.Position = cs.offset;
    var csValue = ReadUInt16();
    header.Nodes.Add(cs);

    // Offset of relocation table; 40h for new-(NE,LE,LX,W3,PE etc.) executable
    var reloc = cs.RelativeToLittleEndian16("Reloc offset");
    BaseStream.Position = reloc.offset;
    var relocValue = (uint)ReadInt16();
    header.Nodes.Add(reloc);

    // Overlay number (0h = main program)
    var overlay = reloc.RelativeToLittleEndian16("Overlay");
    BaseStream.Position = overlay.offset;
    var overlayValue = ReadInt16();
    header.Nodes.Add(overlay);


    if (overlayValue == 0x0000)
        overlay.Text += " = main program";
    else
        throw new Exception("SAMPLE PLZ- Unseen overlay value 0x" + overlayValue.ToString("x4"));


    // look for extended headers (PKLITE etc) for traditional MZ executables


    BaseStream.Position = 0x001C;
    if (ReadByte() == 0x01 && ReadByte() == 0x00 && ReadByte() == 0xFB) {

        // Borland TLINK
        // OFFSET              Count TYPE   Description
        // 001Ch                   2 byte   ?? (apparently always 01h 00h)
        // 001Eh                   1 byte   ID=0FBh
        // 001Fh                   1 byte   TLink version, major in high nybble
        // 0020h                   2 byte   ??

        Console.WriteLine("borland TLINK (DOS)");

        var tlink = overlay.RelativeTo("Borland TLINK header", 6);

        var tlinkHeader = overlay.RelativeTo("Identifier", 3);
        tlink.Nodes.Add(tlinkHeader);

        var tlinkVersion = tlinkHeader.RelativeToVersionMajorMinor8("Version");   // XXX hi & low nibble. 0x30  = 3.0
        tlink.Nodes.Add(tlinkVersion);

        var tlinkExtra = tlinkVersion.RelativeTo("Unknown", 2); // always "jr" ?
        tlink.Nodes.Add(tlinkExtra);

        header.Nodes.Add(tlink);
    }


    BaseStream.Position = 0x001C;
    var lzexeId = ReadStringZ();
    if (lzexeId.Length >= 4 && lzexeId.Substring(0, 2) == "LZ") {
        var lzexe = overlay.RelativeTo("LZEXE compressed executable header", 4);

        var lzexeIdentifier = overlay.RelativeTo("Identifier", 2);
        lzexe.Nodes.Add(lzexeIdentifier);

        string lzexeVerName = "UNKNOWN VERSION";
        switch (lzexeId.Substring(2, 2)) {
        case "09":
            lzexeVerName = "0.9";
            break;
        case "91":
            lzexeVerName = "0.91";
            break;
        }

        var lzexeVersion = lzexeIdentifier.RelativeTo("Version " + lzexeVerName, 2);
        lzexe.Nodes.Add(lzexeVersion);

        header.Nodes.Add(lzexe);
    }



    BaseStream.Position = 0x001E;
    var pkliteId = ReadStringZ();
    if (pkliteId.Length > 6 && pkliteId.Substring(0, 6) == "PKLITE") {
        var pklite = overlay.RelativeTo("PKLITE compressed executable header", (uint)(2 + pkliteId.Length));

        // 001Ch                   1 byte   Minor version number
        var pkliteMinVer = overlay.RelativeToByte("Minor version");
        pklite.Nodes.Add(pkliteMinVer);

        //001Dh                   1 byte   Bit mapped :
        //                                 0-3 - major version
        //                                 4 - Extra compression
        //                                 5 - Multi-segment file
        var pkliteMajorVer = pkliteMinVer.RelativeToByte("Major version");
        pklite.Nodes.Add(pkliteMajorVer);

        var pkliteString = pkliteMajorVer.RelativeTo("Identifier", (uint)pkliteId.Length);
        pklite.Nodes.Add(pkliteString);

        header.Nodes.Add(pklite);
    }



    res.Add(header);

    if (relocValue != 0x0040) {

        if (relocCntValue > 0) {
            // After the header, there follow the relocation items, which are used to span
            // multpile segments. The relocation items have the following format :
            // OFFSET              Count TYPE   Description
            // 0000h                   1 word   Offset within segment
            // 0002h                   1 word   Segment of relocation
            // To get the position of the relocation within the file, you have to compute the
            // physical adress from the segment:offset pair, which is done by multiplying the
            // segment by 16 and adding the offset and then adding the offset of the binary
            // start. Note that the raw binary code starts on a paragraph boundary within the
            // executable file. All segments are relative to the start of the executable in
            // memory, and this value must be added to every segment if relocation is done
            // manually

            BaseStream.Position = relocValue;

            var relocChunk = new Chunk("Relocation Table");
            relocChunk.offset = relocValue;
            relocChunk.length = relocCntValue * 4;
            header.Nodes.Add(relocChunk);

            for (int i = 1; i <= relocCntValue; i++) {
                ushort relocOffset = ReadUInt16();
                ushort relocSegment = ReadUInt16();
                var offset = ((headerSizeValue + relocSegment) * 16) + relocOffset - 1;

                string tmp = "Reloc " + i + " " + relocSegment.ToString("x4") + ":" + relocOffset.ToString("x4") + " => " + offset.ToString("x6");

                var relocItem = new Chunk(tmp);
                relocItem.offset = relocValue + ((i - 1) * 4);
                relocItem.length = 4;
                relocChunk.Nodes.Add(relocItem);
            }
        }

        Log("HeaderSize = " + headerSizeValue.ToString("x4"));
        Log("seg * 16 = " + (csValue * 16).ToString("x4"));
        EntryPoint = (((headerSizeValue + csValue) * 16) + ipValue);

        Log("EntryPoint CS:IP = " + csValue.ToString("x4") + ":" + ipValue.ToString("x4"));

        // decode to seg:offset "exact" address:
        int xx = (csValue * 16) + ipValue;
        Log("  tmp = 0x" + xx.ToString("x4"));

        Log("  = linear offset 0x" + EntryPoint.ToString("x6"));



    } else {
        // 40h for new-(NE,PE,LE,LX,W3 etc.) executable

        var subHead = ParseSubHeader(overlay);
        header.Nodes.Add(subHead);

        BaseStream.Position = ExtendedHeaderOffset;
        char b1 = ReadChar();
        char b2 = ReadChar();

        if (b1 == 'N' && b2 == 'E') {
            // Win16 / OS/2 file
            var neHead = ParseNEHeader();
            header.Nodes.Add(neHead);
        } else if (b1 == 'P' && b2 == 'E') {
            // Win32
            var peHead = ParsePEHeader();
            header.Nodes.Add(peHead);
        } else {
            throw new Exception("TODO unknown header at 0x" + ExtendedHeaderOffset.ToString("x4") + ": " + b1 + ", " + b2);
        }
    }

    // calculates real offset from virtual address
    foreach (var tmp in sections) {
        var chunk = new Chunk("Section " + tmp.Text);
        chunk.length = tmp.length;
        if (chunk.length > 0) {
            tmp.realOffset = FileOffsetFromVirtualAddress(tmp.virtualOffset);
            chunk.offset = tmp.realOffset;
            res.Add(chunk);
        }
    }

    // calculates real offset from virtual address
    foreach (var tmp in dataDirectory) {
        var chunk = new Chunk("DataDirectory " + tmp.Text);
        chunk.length = tmp.length;

        if (chunk.length > 0) {
            tmp.realOffset = FileOffsetFromVirtualAddress(tmp.virtualOffset);
            chunk.offset = tmp.realOffset;

            // TODO use ImportChunk class or soemthing
            if (tmp.Text == "Imports") {
                var OriginalFirstThunk = new LittleEndian32BitChunk("Original First Thunk");
                OriginalFirstThunk.offset = chunk.offset;
                BaseStream.Position = OriginalFirstThunk.offset;
                int OriginalFirstThunkValue = ReadInt32();

                if (OriginalFirstThunkValue > 0) {
                    long OriginalFirstThunkRealOffset = FileOffsetFromVirtualAddress(OriginalFirstThunkValue);
                    //OriginalFirstThunk.Text += " real offset " + OriginalFirstThunkRealOffset.ToString("x8");

                    var OriginalFirstData = new Chunk("Original First Data");
                    OriginalFirstData.offset = OriginalFirstThunkRealOffset;
                    OriginalFirstData.length = 6; // XXX empty-entry-terminated array

                    OriginalFirstThunk.Nodes.Add(OriginalFirstData);
                }


                chunk.Nodes.Add(OriginalFirstThunk);

                var TimeDateStamp = OriginalFirstThunk.RelativeToLittleEndianDateStamp("TimeDateStamp");
                chunk.Nodes.Add(TimeDateStamp);

                var ForwarderChain = TimeDateStamp.RelativeToLittleEndian32("Forwarder Chain");
                chunk.Nodes.Add(ForwarderChain);

                var Name = ForwarderChain.RelativeToLittleEndian32("Name");
                BaseStream.Position = Name.offset;
                int NameValue = ReadInt32();
                if (NameValue > 0) {
                    long realNameOffset = FileOffsetFromVirtualAddress(NameValue);

                    var NameData = new ZeroTerminatedStringChunk();
                    NameData.offset = realNameOffset;
                    NameData.length = 16;

                    string realName = "XX FIX FIX FIXME TODO NAME";  // NameData.GetString(d);

                    //Log("realName = " + realName);

                    NameData.length = (uint)(realName.Length + 1); // 0-terminated string
                    NameData.Text = realName;
                    Name.Nodes.Add(NameData);
                }

                chunk.Nodes.Add(Name);

                var FirstThunk = Name.RelativeToLittleEndian32("First Thunk");
                BaseStream.Position = FirstThunk.offset;
                int FirstThunkValue = ReadInt32();


                var FirstData = new Chunk("First Data");
                FirstData.offset = FileOffsetFromVirtualAddress(FirstThunkValue);
                FirstData.length = 6; // XXX empty-entry-terminated array
                FirstThunk.Nodes.Add(FirstData);


                chunk.Nodes.Add(FirstThunk);
            }

            res.Add(chunk);
        }
    }

    return res;
}

public class SectionPointer
{
    public long virtualOffset;
    public long realOffset;
    public uint length;
    public string Text;
}

public List<SectionPointer> sections = new List<SectionPointer>();
public List<SectionPointer> dataDirectory = new List<SectionPointer>();
public long EntryPoint;
long ExtendedHeaderOffset;
public long ExeHeaderLength;

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

private Chunk ParseSubHeader(Chunk previous)
{
    var subHead = new Chunk();

    subHead.offset = previous.offset + previous.length;
    subHead.Text = "Extended header";

    // New Executable header
    var unknown = previous.RelativeTo("Reserved", 8);
    subHead.Nodes.Add(unknown);

    var oemId = unknown.RelativeToLittleEndian16("OEM id");
    subHead.Nodes.Add(oemId);

    var oemInfo = oemId.RelativeToLittleEndian16("OEM info");
    subHead.Nodes.Add(oemInfo);

    var reserved = oemInfo.RelativeTo("Reserved", 20);
    subHead.Nodes.Add(reserved);

    // Offset of extended executable header from start of file (or 0 if plain MZ executable)
    var neHeader = reserved.RelativeToLittleEndian32("Offset of header");
    BaseStream.Position = neHeader.offset;
    this.ExtendedHeaderOffset = ReadUInt32();
    subHead.Nodes.Add(neHeader);

    subHead.length = (uint)((neHeader.offset + neHeader.length) - subHead.offset);

    // TODO: wrap "this program cant be run in dos mode" in a chunk, how to detect size? start offset is 0x40
    return subHead;
}

private Chunk ParseNEHeader()
{
    // used in win16 (Windows 3) and OS/2 executables
    this.name = "New Executable (Win16)";

    var neHead = new Chunk("NE header");
    neHead.offset = ExtendedHeaderOffset;
    neHead.length = 0x40;

    var neIdentifier = new LittleEndian16BitChunk("NE identifier");
    neIdentifier.offset = ExtendedHeaderOffset;
    neHead.Nodes.Add(neIdentifier);

    var LinkerVersion = neIdentifier.RelativeToVersionMajorMinor16("Linker version");
    neHead.Nodes.Add(LinkerVersion);


    var OffsetToEntryTable = LinkerVersion.RelativeToLittleEndian16("EntryTableOffset");
    BaseStream.Position = OffsetToEntryTable.offset;
    var OffsetToEntryTableValue = neHead.offset + ReadUInt16();
    neHead.Nodes.Add(OffsetToEntryTable);

    var EntryTableLength = OffsetToEntryTable.RelativeToLittleEndian16("EntryTableLength");
    BaseStream.Position = EntryTableLength.offset;
    var EntryTableLengthValue = ReadUInt16();
    neHead.Nodes.Add(EntryTableLength);


    var FileLoadCrc = EntryTableLength.RelativeToLittleEndian32("File Crc");
    neHead.Nodes.Add(FileLoadCrc);

    var FormatFlags = FileLoadCrc.RelativeToLittleEndian16("Format flags");
    BaseStream.Position = FormatFlags.offset;
    var FormatFlagsValue = ReadUInt16();
    neHead.Nodes.Add(FormatFlags);

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
    }

    var AutoDataSegmentIndex = FormatFlags.RelativeToLittleEndian16("Auto Data Segment Index");
    neHead.Nodes.Add(AutoDataSegmentIndex);

    var InitialLocalHeapSize = AutoDataSegmentIndex.RelativeToLittleEndian16("Initial Local Heap Size");
    neHead.Nodes.Add(InitialLocalHeapSize);

    var InitialStackSize = InitialLocalHeapSize.RelativeToLittleEndian16("Initial Stack Size");
    neHead.Nodes.Add(InitialStackSize);

    var EntryPointCSIP = InitialStackSize.RelativeToLittleEndian32("Entry Point CS:IP");
    neHead.Nodes.Add(EntryPointCSIP);
    Log("NE Entry point = CS:IP XXXX XFIXME READ PARSE DECODE");

    var InitialStackPointer = EntryPointCSIP.RelativeToLittleEndian32("Stack Pointer SS:SP");
    neHead.Nodes.Add(InitialStackPointer);

    // # of entries in segment table
    var SegmentCount = InitialStackPointer.RelativeToLittleEndian16("Segment table entries");
    BaseStream.Position = SegmentCount.offset;
    var SegmentCountValue = ReadUInt16();
    neHead.Nodes.Add(SegmentCount);

    // # of entries in module reference table
    var ModuleReferenceCount = SegmentCount.RelativeToLittleEndian16("Module reference table entires");
    BaseStream.Position = ModuleReferenceCount.offset;
    var ModuleReferenceCountValue = ReadUInt16();
    neHead.Nodes.Add(ModuleReferenceCount);

    var NonresidentNamesTableSize = ModuleReferenceCount.RelativeToLittleEndian16("NonresidentNamesTableSize");
    neHead.Nodes.Add(NonresidentNamesTableSize);

    var OffsetSegmentTable = NonresidentNamesTableSize.RelativeToLittleEndian16("OffsetSegmentTable");
    BaseStream.Position = OffsetSegmentTable.offset;
    var OffsetSegmentTableValue = neHead.offset + ReadInt16();
    neHead.Nodes.Add(OffsetSegmentTable);

    var OffsetResourceTable = OffsetSegmentTable.RelativeToLittleEndian16("OffsetResourceTable");
    BaseStream.Position = OffsetResourceTable.offset;
    var OffsetResourceTableValue = neHead.offset + ReadInt16();
    neHead.Nodes.Add(OffsetResourceTable);

    var OffsetResidentNamesTable = OffsetResourceTable.RelativeToLittleEndian16("OffsetResidentNamesTable");
    BaseStream.Position = OffsetResidentNamesTable.offset;
    var OffsetResidentNamesTableValue = neHead.offset + ReadInt16();
    neHead.Nodes.Add(OffsetResidentNamesTable);



    var OffsetModuleReferenceTable = OffsetResidentNamesTable.RelativeToLittleEndian16("OffsetModuleReferenceTable");
    BaseStream.Position = OffsetModuleReferenceTable.offset;
    var OffsetModuleReferenceTableValue = neHead.offset + ReadInt16();
    neHead.Nodes.Add(OffsetModuleReferenceTable);



    // (array of counted strings, terminated with a string of length 00h)
    var OffsetImportedNamesTable = OffsetModuleReferenceTable.RelativeToLittleEndian16("OffsetImportedNamesTable");
    BaseStream.Position = OffsetImportedNamesTable.offset;
    var OffsetImportedNamesTableValue = neHead.offset + ReadInt16();
    neHead.Nodes.Add(OffsetImportedNamesTable);



    // Offset from start of file to nonresident names table
    var OffsetNonresidentNamesTable = OffsetImportedNamesTable.RelativeToLittleEndian32("OffsetNonresidentNamesTable");
    BaseStream.Position = OffsetNonresidentNamesTable.offset;
    var OffsetNonresidentNamesTableValue = ReadInt32();
    neHead.Nodes.Add(OffsetNonresidentNamesTable);


    // Count of moveable entry point listed in entry table
    var MovableEntryPointsInEntryTable = OffsetNonresidentNamesTable.RelativeToLittleEndian16("MovableEntryPointsInEntryTable");
    BaseStream.Position = MovableEntryPointsInEntryTable.offset;
    var MovableEntryPointsInEntryTableValue = ReadInt32();
    neHead.Nodes.Add(MovableEntryPointsInEntryTable);

    //  File alignment size shift count, 0 is equivalent to 9 (default 512-byte pages)
    var FileAlignmentSizeShift = MovableEntryPointsInEntryTable.RelativeToLittleEndian16("FileAlignmentSizeShift");
    neHead.Nodes.Add(FileAlignmentSizeShift);

    // Number of resource table entries
    var ResourceTableEntries = FileAlignmentSizeShift.RelativeToLittleEndian16("Resources");
    neHead.Nodes.Add(ResourceTableEntries);

    var TargetOs = ResourceTableEntries.RelativeToByte("Target OS");

    BaseStream.Position = TargetOs.offset;
    var TargetOsValue = ReadByte();
    neHead.Nodes.Add(TargetOs);

    string TargetOsName = "Unknown";
    switch (TargetOsValue) {
    case 1:
        TargetOsName = "OS/2";
        break;
    case 2:
        TargetOsName = "Windows";
        break;
    case 3:
        TargetOsName = "European MS-DOS 4.x";
        break;
    case 4:
        TargetOsName = "Windows 386";
        break;
    case 5:
        TargetOsName = "BOSS (Borland Operating System Services)";
        break;
    }
    TargetOs.Text += " = " + TargetOsName;

    var AdditionalFlags = TargetOs.RelativeToByte("Additional Flags");
    BaseStream.Position = AdditionalFlags.offset;
    var AdditionalFlagsValue = ReadByte();
    neHead.Nodes.Add(AdditionalFlags);

    if ((AdditionalFlagsValue & 0x01) != 0)
        Log("TODO Long filename support?!?!");

    if ((AdditionalFlagsValue & 0x02) != 0) {
        // The executable file contains a Windows 2.x
        // application that runs in version 3.x protected mode.
        Log("TODO WIN2 protected mode");
    }

    if ((AdditionalFlagsValue & 0x04) != 0) {
        // If this bit is set, the  executable file contains a Windows 2.x
        // application that supports proportional fonts.
        Log("TODO WIN2 preoportional fonts");
    }
    if ((AdditionalFlagsValue & 0x08) != 0) {
        Log("TODO Executable has fastload area");
    }

    // NOTE: only used by windows
    var OffsetToFastload = AdditionalFlags.RelativeToLittleEndian16("OffsetToFastload");
    BaseStream.Position = OffsetToFastload.offset;
    var OffsetToFastloadValue = ReadInt16();
    neHead.Nodes.Add(OffsetToFastload);


    // NOTE: only used by windows
    // offset to segment reference thunks or length of gangload area.
    var LengthOfFastload = OffsetToFastload.RelativeToLittleEndian16("LengthOfFastload");
    BaseStream.Position = LengthOfFastload.offset;
    var LengthOfFastloadValue = ReadInt16();
    neHead.Nodes.Add(LengthOfFastload);
    Log("TODO offset to OffsetToFastload Area 0x" + OffsetToFastloadValue.ToString("x4") + ", length=0x" + LengthOfFastloadValue.ToString("x4"));

    var ReservedWord = LengthOfFastload.RelativeToLittleEndian16("Reserved");
    neHead.Nodes.Add(ReservedWord);

    // NOTE: only used by windows
    // TODO add to version datatype, MINOR.MAJOR byte order, 2 bytes
    var ExpectedWindowsVersion = ReservedWord.RelativeToLittleEndian16("ExpectedWindowsVersion");
    neHead.Nodes.Add(ExpectedWindowsVersion);



    neHead.Nodes.Add(ParseNEModuleReferenceTable(OffsetModuleReferenceTableValue, ModuleReferenceCountValue));

    neHead.Nodes.Add(ParseNEEntryTable(OffsetToEntryTableValue, EntryTableLengthValue));

    neHead.Nodes.Add(ParseNESegmentTable(OffsetSegmentTableValue, SegmentCountValue));

    neHead.Nodes.Add(ParseNEImportedTable(OffsetImportedNamesTableValue));

    neHead.Nodes.Add(ParseNEResidentTable(OffsetResidentNamesTableValue));

    neHead.Nodes.Add(ParseNENonResidentTable(OffsetNonresidentNamesTableValue));

    neHead.Nodes.Add(ParseNEResourceTable(OffsetResourceTableValue));

    return neHead;
}

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

private static string ByteArrayToString(byte[] arr)
{
    var s = new StringBuilder();
    foreach (byte b in arr)
        s.Append((char)b);

    return s.ToString();
}

private Chunk ParsePEHeader()
{
    Log("PE header found");

    this.name = "PE executable (Win32)";

    var peHead = new Chunk();
    peHead.offset = ExtendedHeaderOffset;
    peHead.Text = "PE header";
    //                    peHead.childs = new List<Chunk>();

    var peIdentifier = new LittleEndian32BitChunk("PE identifier");
    peIdentifier.offset = ExtendedHeaderOffset;
    peHead.Nodes.Add(peIdentifier);

    // COFF header
    var Machine = peIdentifier.RelativeToLittleEndian16("Machine = ");

    BaseStream.Position = Machine.offset;
    int MachineValue = ReadInt16();
    switch (MachineValue) {
    case 0x14c:
        Machine.Text += "Intel 386";
        break;

    default:
        Machine.Text += "Unknown";
        throw new Exception("Unrecognized machine: " + MachineValue.ToString("x4"));
    }
    peHead.Nodes.Add(Machine);

    // NumberOfSections
    var NumberOfSections = Machine.RelativeToLittleEndian16("Sections");
    peHead.Nodes.Add(NumberOfSections);

    BaseStream.Position = NumberOfSections.offset;
    var NumberOfSectionsValue = (uint)ReadInt16();

    var TimeDateStamp = NumberOfSections.RelativeToLittleEndianDateStamp("TimeDateStamp");
    peHead.Nodes.Add(TimeDateStamp);

    var PointerToSymbolTable = TimeDateStamp.RelativeToLittleEndian32("Pointer To Symbol Table");
    peHead.Nodes.Add(PointerToSymbolTable);

    // NumberOfSymbols
    var NumberOfSymbols = PointerToSymbolTable.RelativeToLittleEndian32("Symbols");
    peHead.Nodes.Add(NumberOfSymbols);

    var SizeOfOptionalHeader = NumberOfSymbols.RelativeToLittleEndian16("Size Of Optional Header");
    peHead.Nodes.Add(SizeOfOptionalHeader);

    BaseStream.Position = SizeOfOptionalHeader.offset;
    var SizeOfOptionalHeaderValue = (uint)ReadInt16();

    var Characteristics = SizeOfOptionalHeader.RelativeToLittleEndian16("Characteristics");
    peHead.Nodes.Add(Characteristics);

    var optHead = new Chunk("Optional header");
    optHead.length = SizeOfOptionalHeaderValue;
    optHead.offset = Characteristics.offset + Characteristics.length;

    var optSignature = new Chunk("Optional signature");
    optSignature.offset = optHead.offset;
    optSignature.length = 2;
    optHead.Nodes.Add(optSignature);

    BaseStream.Position = optSignature.offset;

    if (ReadByte() != 0x0b || ReadByte() != 0x01)
        throw new Exception("Unrecognized optional header value found");

    var LinkerVersion = optSignature.RelativeToVersionMajorMinor16("Linker version");
    LinkerVersion.Text = "Linker version";
    optHead.Nodes.Add(LinkerVersion);

    var SizeOfCode = LinkerVersion.RelativeToLittleEndian32("Size of Code");
    optHead.Nodes.Add(SizeOfCode);

    var SizeOfInitializedData = SizeOfCode.RelativeToLittleEndian32("Size Of Initialized Data");
    optHead.Nodes.Add(SizeOfInitializedData);

    var SizeOfUninitializedData = SizeOfInitializedData.RelativeToLittleEndian32("Size Of Uninitialized Data");
    optHead.Nodes.Add(SizeOfUninitializedData);

    //The RVA of the code entry point
    var AddressOfEntryPoint = SizeOfUninitializedData.RelativeToLittleEndian32("Address Of Entry Point");
    BaseStream.Position = AddressOfEntryPoint.offset;
    int AddressOfEntryPointValue = ReadInt32();
    Log("XXXXX entry offset = " + AddressOfEntryPointValue.ToString("x6"));
    this.EntryPoint = AddressOfEntryPointValue;

    optHead.Nodes.Add(AddressOfEntryPoint);

    var BaseOfCode = AddressOfEntryPoint.RelativeToLittleEndian32("Base of Code");
    optHead.Nodes.Add(BaseOfCode);

    var BaseOfData = BaseOfCode.RelativeToLittleEndian32("Base of Data");
    optHead.Nodes.Add(BaseOfData);

    var ImageBase = BaseOfData.RelativeToLittleEndian32("Image Base");
    optHead.Nodes.Add(ImageBase);

    var SectionAlignment = ImageBase.RelativeToLittleEndian32("Section Alignment");
    optHead.Nodes.Add(SectionAlignment);

    var FileAlignment = SectionAlignment.RelativeToLittleEndian32("File Alignment");
    optHead.Nodes.Add(FileAlignment);

    var OSVersion = FileAlignment.RelativeToVersionMajorMinor32("OS Version");
    optHead.Nodes.Add(OSVersion);

    var ImageVersion = OSVersion.RelativeToVersionMajorMinor32("Image Version");
    optHead.Nodes.Add(ImageVersion);

    var SubsystemVersion = ImageVersion.RelativeToVersionMajorMinor32("Subsystem Version");
    optHead.Nodes.Add(SubsystemVersion);

    var Reserved = SubsystemVersion.RelativeToLittleEndian32("Reserved");
    optHead.Nodes.Add(Reserved);

    var SizeOfImage = Reserved.RelativeToLittleEndian32("Size of Image");
    optHead.Nodes.Add(SizeOfImage);

    var SizeOfHeaders = SizeOfImage.RelativeToLittleEndian32("Size of Headers");
    optHead.Nodes.Add(SizeOfHeaders);


    var Checksum = SizeOfHeaders.RelativeToLittleEndian32("Checksum");
    optHead.Nodes.Add(Checksum);

    var Subsystem = Checksum.RelativeToLittleEndian16("Subsystem = ");
    optHead.Nodes.Add(Subsystem);

    BaseStream.Position = Subsystem.offset;
    int SubsystemValue = ReadInt16();

    switch (SubsystemValue) {
    case 0x0001:
        Subsystem.Text += "Native";
        break;
    case 0x0002:
        Subsystem.Text += "GUI";
        break;
    case 0x0003:
        Subsystem.Text += "Console";
        break;
    case 0x0005:
        Subsystem.Text += "OS/2";
        break;
    case 0x0007:
        Subsystem.Text += "POSIX";
        break;
    default:
        Log("Unknown " + SubsystemValue.ToString("x4"));
        break;
    }

    var DLLCharacteristics = Subsystem.RelativeToLittleEndian16("DLL Characteristics");
    optHead.Nodes.Add(DLLCharacteristics);

    var SizeOfStackReserve = DLLCharacteristics.RelativeToLittleEndian32("Size Of Stack Reserve");
    optHead.Nodes.Add(SizeOfStackReserve);

    var SizeOfStackCommit = SizeOfStackReserve.RelativeToLittleEndian32("Size Of Stack Commit");
    optHead.Nodes.Add(SizeOfStackCommit);

    var SizeOfHeapReserve = SizeOfStackCommit.RelativeToLittleEndian32("Size Of Heap Reserve");
    optHead.Nodes.Add(SizeOfHeapReserve);

    var SizeOfHeapCommit = SizeOfHeapReserve.RelativeToLittleEndian32("Size Of Heap Commit");
    optHead.Nodes.Add(SizeOfHeapCommit);

    var LoaderFlags = SizeOfHeapCommit.RelativeToLittleEndian32("Loader Flags");
    optHead.Nodes.Add(LoaderFlags);

    var NumberOfRvaAndSizes = LoaderFlags.RelativeToLittleEndian32("Number of RVA");
    optHead.Nodes.Add(NumberOfRvaAndSizes);

    BaseStream.Position = NumberOfRvaAndSizes.offset;
    var NumberOfRvaAndSizesValue = (uint)ReadInt32();
    Log("NumberOfRvaAndSizesValue = " + NumberOfRvaAndSizesValue);

    if (NumberOfRvaAndSizesValue != 16)
        throw new Exception("odd number of RVA:s = " + NumberOfRvaAndSizesValue);

    var DataDirectory = new Chunk("Data Directory");
    DataDirectory.offset = NumberOfRvaAndSizes.offset + NumberOfRvaAndSizes.length;
    DataDirectory.length = NumberOfRvaAndSizesValue * 8;

    for (int i = 0; i < NumberOfRvaAndSizesValue; i++) {
        var RVAChunk = new Chunk();
        RVAChunk.length = 8;
        RVAChunk.offset = DataDirectory.offset + (i * RVAChunk.length);
        switch (i) {
        case 0:
            RVAChunk.Text = "Exports";
            break;
        case 1:
            RVAChunk.Text = "Imports";
            break;
        case 2:
            RVAChunk.Text = "Resources";
            break;
        case 5:
            RVAChunk.Text = "Base reolcations";
            break;
        case 9:
            RVAChunk.Text = "Thread Local Storage";
            break;
        case 12:
            RVAChunk.Text = "Import Address Table";
            break;
        case 14:
            RVAChunk.Text = "CLR Header";
            break;
        default:
            RVAChunk.Text = "Data Directory # " + i;
            break;
        }

        var VirtualAddress = new Chunk();
        VirtualAddress.length = 4;
        VirtualAddress.offset = RVAChunk.offset;
        BaseStream.Position = VirtualAddress.offset;
        int VirtualAddressValue = ReadInt32();
        VirtualAddress.Text = "Virtual Address = " + VirtualAddressValue.ToString("x8");
        RVAChunk.Nodes.Add(VirtualAddress);

        var RVASize = new Chunk();
        RVASize.length = 4;
        RVASize.offset = VirtualAddress.offset + VirtualAddress.length;
        BaseStream.Position = RVASize.offset;
        var RVASizeValue = (uint)ReadInt32();

        if (RVASizeValue == 0)
            RVAChunk.Text = "Empty";

        RVASize.Text = "Size = " + RVASizeValue;
        RVAChunk.Nodes.Add(RVASize);

        // TODO: create list of DataDirectory entries. later; calculate their physical offsets from virtual offsets

        var dd = new SectionPointer();
        dd.Text = RVAChunk.Text;
        dd.virtualOffset = VirtualAddressValue;
        dd.length = RVASizeValue;
        dd.realOffset = 0; // TODO calc later

        if (RVASizeValue > 0) {
            dataDirectory.Add(dd);
            DataDirectory.Nodes.Add(RVAChunk);
        }
    }

    optHead.Nodes.Add(DataDirectory);

    peHead.length = (uint)((optHead.offset + optHead.length) - peHead.offset);

    var SectionsOverview = new Chunk();
    SectionsOverview.length = NumberOfSectionsValue * 40;
    SectionsOverview.offset = peHead.offset + peHead.length;
    SectionsOverview.Text = "Sections";

    if (SizeOfOptionalHeaderValue > 0)
        peHead.Nodes.Add(optHead);

    if (NumberOfSectionsValue > 0)
        peHead.Nodes.Add(SectionsOverview);

    for (int i = 0; i < NumberOfSectionsValue; i++) {
        var SectionChunk = new Chunk();
        SectionChunk.length = 40;
        SectionChunk.offset = peHead.offset + peHead.length + (i * SectionChunk.length);

        // Section Name - common names are .text .data .bss
        var SectionName = new ZeroTerminatedStringChunk("Section name", 8);
        SectionName.offset = SectionChunk.offset;
        SectionChunk.Nodes.Add(SectionName);

        var SectioNameValue = "XXX FIXME AGAIN"; // SectionName.GetString(d);

        SectionChunk.Text = "Section " + SectioNameValue;

        // Size of the section once it is loaded to memory
        var SectionSize = SectionName.RelativeToLittleEndian32("Section size loaded");
        SectionChunk.Nodes.Add(SectionSize);

        // RVA (location) of section once it is loaded to memory
        var RVALocation = SectionSize.RelativeToLittleEndian32("RVA location of section");
        BaseStream.Position = RVALocation.offset;
        int RVALocationValue = ReadInt32();
        SectionChunk.Nodes.Add(RVALocation);

        // Physical size of section on disk
        var PhysSize = RVALocation.RelativeToLittleEndian32("Physical size of section");
        BaseStream.Position = PhysSize.offset;
        var PhysSizeValue = (uint)ReadInt32();
        SectionChunk.Nodes.Add(PhysSize);

        // Physical location of section on disk (from start of disk image)
        var PhysOffset = PhysSize.RelativeToLittleEndian32("Physical offset of section");
        BaseStream.Position = PhysOffset.offset;
        int PhysOffsetValue = ReadInt32();
        SectionChunk.Nodes.Add(PhysOffset);

        // Reserved (usually zero) (used in object formats)
        var Reserved12 = PhysOffset.RelativeTo("Reserved", 12);
        SectionChunk.Nodes.Add(Reserved12);

        // Section flags
        var SectionFlags = Reserved12.RelativeToLittleEndian32("Section flags");
        BaseStream.Position = SectionFlags.offset;
        int SectionFlagsValue = ReadInt32();

        if ((SectionFlagsValue & 0x00000020) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000020 = Code";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x00000040) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000040 = Initialized data";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x00000080) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000080 = Uninitialized data";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x00000200) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000200 = Info";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x02000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x02000000 = Discardable";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x10000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x10000000 = Shared";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x20000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x20000000 = Executable";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x40000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x40000000 = Readable";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x80000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x80000000 = Writeable";
            SectionFlags.Nodes.Add(note);
        }

        SectionChunk.Nodes.Add(SectionFlags);


        var sectionPointer = new SectionPointer();
        sectionPointer.length = PhysSizeValue;
        sectionPointer.realOffset = PhysOffsetValue;
        sectionPointer.virtualOffset = RVALocationValue;
        sectionPointer.Text = SectioNameValue;
        if (sectionPointer.length > 0)
            sections.Add(sectionPointer);

        SectionsOverview.Nodes.Add(SectionChunk);
    }

    peHead.length += SectionsOverview.length;

    return peHead;
}

// Calculates the 16-bit checksum used in the orginal MZ header
public ushort CalculateChecksum16bit()
{
    // based on code from http://support.microsoft.com/KB/71971
    BaseStream.Position = 0;

    ushort sum16 = 0;

    // NOTE if we skip offset 0x0012, we get 0x0000 ???

    for (int x = 0; x < BaseStream.Length / 2; x++) {
        //if (x == 0x0006)
        //    continue;
        sum16 += ReadUInt16();
    }

    // make sure and get the last byte if odd size...
    if (BaseStream.Length % 2 != 0) {
        sum16 += ReadByte();
    }

    return sum16;
}
*/
