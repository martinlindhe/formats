package parse

// STATUS: 60%
//   0% NE (win?!)
//   0% PE (win?!)

import (
	"encoding/binary"
	"fmt"
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

	offset := int64(0)
	mz := Layout{
		Offset: offset,
		Length: 28, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: offset, Length: 2, Info: "magic", Type: ASCII},
			Layout{Offset: offset + 2, Length: 2, Info: "extra bytes", Type: Uint16le},
			Layout{Offset: offset + 4, Length: 2, Info: "pages", Type: Uint16le},
			Layout{Offset: offset + 6, Length: 2, Info: "relocation items", Type: Uint16le},
			Layout{Offset: offset + 8, Length: 2, Info: "header size in paragraphs", Type: Uint16le}, // 1 paragraph = group of 16 bytes
			Layout{Offset: offset + 10, Length: 2, Info: "min allocation", Type: Uint16le},
			Layout{Offset: offset + 12, Length: 2, Info: "max allocation", Type: Uint16le},
			Layout{Offset: offset + 14, Length: 2, Info: "initial ss", Type: Uint16le},
			Layout{Offset: offset + 16, Length: 2, Info: "initial sp", Type: Uint16le},
			Layout{Offset: offset + 18, Length: 2, Info: "checksum", Type: Uint16le},
			Layout{Offset: offset + 20, Length: 2, Info: "initial ip", Type: Uint16le},
			Layout{Offset: offset + 22, Length: 2, Info: "initial cs", Type: Uint16le},

			// Offset of relocation table; 40h for new-(NE,LE,LX,W3,PE etc.) executable
			Layout{Offset: offset + 24, Length: 2, Info: "relocation offset", Type: Uint16le},
			Layout{Offset: offset + 26, Length: 2, Info: "overlay", Type: Uint16le},
		}}

	res.Layout = append(res.Layout, mz)

	custom := findCustomDOSHeaders(file)
	if custom != nil {
		res.Layout = append(res.Layout, *custom)
	}

	hdrSizeInParagraphs, _ := readUint16le(file, offset+8)
	ip, _ := readUint16le(file, offset+20)
	cs, _ := readUint16le(file, offset+22)
	relocOffset, _ := readUint16le(file, offset+24)

	if relocOffset == 0x40 {
		// XXX NE, PE etc

		/*

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
		*/
	} else {
		relocItems, _ := readUint16le(file, offset+6)
		if relocItems > 0 {
			offset = int64(relocOffset)
			reloc := Layout{
				Offset: offset,
				Length: int64(relocItems) * 4,
				Info:   "relocation table",
				Type:   Group}

			for i := 1; i <= int(relocItems); i++ {
				reloc.Childs = append(reloc.Childs, []Layout{
					Layout{Offset: offset, Length: 2, Info: "offset " + fmt.Sprintf("%d", i), Type: Uint16le},
					Layout{Offset: offset + 2, Length: 2, Info: "segment " + fmt.Sprintf("%d", i), Type: Uint16le},
				}...)
				offset += 4
			}
			res.Layout = append(res.Layout, reloc)
		}
	}

	exeStart := int64(((hdrSizeInParagraphs + cs) * 16) + ip)

	// XXX disasm until first ret or sth ???
	offset = exeStart
	codeChunk := Layout{
		Offset: offset,
		Length: 4,
		Info:   "program XXX",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: offset, Length: 4, Info: "XXX", Type: Bytes},
		}}

	res.Layout = append(res.Layout, codeChunk)

	return &res, nil
}

/*

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
