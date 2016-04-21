package archive

// Windows Imaging format

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WIM(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isWIM(&c.Header) {
		return nil, nil
	}
	return parseWIM(c.File, c.ParsedLayout)
}

func isWIM(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'M' || b[1] != 'S' || b[2] != 'W' || b[3] != 'I' ||
		b[4] != 'M' {
		return false
	}
	return true
}

func parseWIM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 5, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 5, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
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