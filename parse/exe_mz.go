package parse

// STATUS: 60%

import (
	"encoding/binary"
	"fmt"
	"os"
	"sort"
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

	res := ParsedLayout{
		FileKind: Executable,
	}

	pos := int64(0)
	mzHeaderLen := int64(28) // XXX
	mz := Layout{
		Offset: pos,
		Length: mzHeaderLen,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: ASCII},
			{Offset: pos + 2, Length: 2, Info: "extra bytes", Type: Uint16le},
			{Offset: pos + 4, Length: 2, Info: "pages", Type: Uint16le},
			{Offset: pos + 6, Length: 2, Info: "relocation items", Type: Uint16le},
			{Offset: pos + 8, Length: 2, Info: "header size in paragraphs", Type: Uint16le}, // 1 paragraph = group of 16 bytes
			{Offset: pos + 10, Length: 2, Info: "min allocation", Type: Uint16le},
			{Offset: pos + 12, Length: 2, Info: "max allocation", Type: Uint16le},
			{Offset: pos + 14, Length: 2, Info: "initial ss", Type: Uint16le},
			{Offset: pos + 16, Length: 2, Info: "initial sp", Type: Uint16le},
			{Offset: pos + 18, Length: 2, Info: "checksum", Type: Uint16le},
			{Offset: pos + 20, Length: 2, Info: "initial ip", Type: Uint16le},
			{Offset: pos + 22, Length: 2, Info: "initial cs", Type: Uint16le},
			{Offset: pos + 24, Length: 2, Info: "relocation offset", Type: Uint16le},
			{Offset: pos + 26, Length: 2, Info: "overlay", Type: Uint16le},
		}}

	res.Layout = append(res.Layout, mz)

	custom := findCustomDOSHeaders(file)
	if custom != nil {
		res.Layout = append(res.Layout, *custom)
	}

	hdrSizeInParagraphs, _ := readUint16le(file, pos+8)
	ip, _ := readUint16le(file, pos+20)
	cs, _ := readUint16le(file, pos+22)
	relocOffset, _ := readUint16le(file, pos+24)

	if relocOffset == 0x40 {
		// 0x40 for new-(NE,LE,LX,W3,PE etc.) executable
		pos += mzHeaderLen

		subHeaderLen := int64(36) // XXX
		res.Layout = append(res.Layout, Layout{
			Offset: pos,
			Length: subHeaderLen,
			Info:   "sub header", // XXX name
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 8, Info: "reserved", Type: Bytes},
				{Offset: pos + 8, Length: 2, Info: "oem id", Type: Uint16le},
				{Offset: pos + 10, Length: 2, Info: "oem info", Type: Uint16le},
				{Offset: pos + 12, Length: 20, Info: "reserved 2", Type: Uint16le},
				{Offset: pos + 32, Length: 4, Info: "start of ext header", Type: Uint32le},
			}})

		newHeaderPos, _ := readUint32le(file, pos+32)

		pos = int64(newHeaderPos)
		newHeaderId, _ := knownLengthASCII(file, pos, 2)

		switch newHeaderId {
		case "LX":
			// OS/2 (32-bit)
			header, _ := parseMZ_LXHeader(file, pos)
			res.Layout = append(res.Layout, header...)
		case "LE":
			// OS/2 (mixed 16/32-bit)
			panic("LE")
		case "NE":
			// Win16, OS/2
			header, _ := parseMZ_NEHeader(file, pos)
			res.Layout = append(res.Layout, header...)
		case "PE":
			// Win32, Win64
			header, _ := parseMZ_PEHeader(file, pos)
			res.Layout = append(res.Layout, header...)
		default:
			// XXX get samples of LE, W3 files
			panic("unknown newHeaderId =" + newHeaderId)
		}
	} else {
		relocItems, _ := readUint16le(file, pos+6)
		if relocItems > 0 {
			pos = int64(relocOffset)
			reloc := Layout{
				Offset: pos,
				Length: int64(relocItems) * 4,
				Info:   "relocation table",
				Type:   Group}

			for i := 1; i <= int(relocItems); i++ {
				reloc.Childs = append(reloc.Childs, []Layout{
					{Offset: pos, Length: 2, Info: "offset " + fmt.Sprintf("%d", i), Type: Uint16le},
					{Offset: pos + 2, Length: 2, Info: "segment " + fmt.Sprintf("%d", i), Type: Uint16le},
				}...)
				pos += 4
			}
			res.Layout = append(res.Layout, reloc)
		}
	}

	exeStart := int64(((hdrSizeInParagraphs + cs) * 16) + ip)

	// XXX disasm until first ret or sth ???
	pos = exeStart
	codeChunk := Layout{
		Offset: pos,
		Length: 4, // XXX
		Info:   "dos entry point",
		Type:   Group,
		Childs: []Layout{
			{Offset: pos, Length: 4, Info: "XXX", Type: Bytes},
		}}

	res.Layout = append(res.Layout, codeChunk)

	sort.Sort(ByLayout(res.Layout))

	return &res, nil
}

/*

override public List<Chunk> GetFileStructure()
{

    ## XXXXX new exes:


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

private static string ByteArrayToString(byte[] arr)
{
    var s = new StringBuilder();
    foreach (byte b in arr)
        s.Append((char)b);

    return s.ToString();
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
