package parse

import (
	"encoding/binary"
	"os"
)

func XZ(file *os.File) (*ParsedLayout, error) {

	if !isXZ(file) {
		return nil, nil
	}
	return parseXZ(file)
}

func isXZ(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [6]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 0xFD || b[1] != '7' || b[2] != 'z' || b[3] != 'X' || b[4] != 'Z' || b[5] != 0x00 {
		return false
	}

	return true
}

func parseXZ(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 6, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 6, Info: "magic", Type: Bytes},
		}})
	return &res, nil
}

/*
string DecodeFlagsValue(ushort flags)
{
    if (flags == 0x0000)
        return "None";

    if (flags == 0x0100)
        return "CRC32";

    if (flags == 0x0400)
        return "CRC64";

    if (flags == 0x0A00)
        return "SHA-256";

    return "Unknown " + flags.ToString("x4");
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a xz");

    List<Chunk> res = new List<Chunk>();

    var identifier = new Chunk();
    identifier.offset = 0;
    identifier.length = 6;
    identifier.Text = "XZ identifier";
    res.Add(identifier);

    var flags = identifier.RelativeToLittleEndian16("Flags");
    var flagsValue = flags.GetValue(BaseStream);

    flags.Text += " = " + DecodeFlagsValue(flagsValue);

    res.Add(flags);

    var crc32 = flags.RelativeToLittleEndian32("CRC32");
    res.Add(crc32);

    return res;
}
*/
