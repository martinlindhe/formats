package parse

// Python bytecode
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

var (
	pythonVersionMagic = map[uint32]string{
		0x00999902: "1.0",
		0x00999903: "1.1-1.2",
		0x0A0D2E89: "1.3",
		0x0A0D1704: "1.4",
		0x0A0D4E99: "1.5",
		0x0A0DC4FC: "1.6",
		0x0A0DC687: "2.0",
		0x0A0DEB2A: "2.1",
		0x0A0DED2D: "2.2",
		0x0A0DF23B: "2.3",
		0x0A0DF26D: "2.4",
		0x0A0DF2B3: "2.5",
		0x0A0DF2D1: "2.6",
		0x0A0DF303: "2.7",
		0x0A0D0C3A: "3.0",
		0x0A0D0C4E: "3.1",
		0x0A0D0C6C: "3.2",
		0x0A0D0C9E: "3.3",
		0x0A0D0CEE: "3.4",
	}
)

func PYTHON(file *os.File) (*ParsedLayout, error) {

	if !isPYTHON(file) {
		return nil, nil
	}
	return parsePYTHON(file)
}

func isPYTHON(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if _, ok := pythonVersionMagic[b]; ok {
		return true
	}

	return false
}

func parsePYTHON(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: Uint32le}, // XXX decode to python version
			}}}}

	return &res, nil
}

/*

// creates a timestamp from a Unix time (seconds since  00:00:00 GMT, Jan. 1, 1970)
private static DateTime MTimeToTimestamp(uint mtime)
{
    DateTime dt = new DateTime(1970, 1, 1, 0, 0, 0, 0);
    dt = dt.AddSeconds(mtime).ToLocalTime();
    return dt;
}

override public List<Chunk> GetFileStructure()
{
    List<Chunk> res = new List<Chunk>();

    var header = new Chunk("PYC header", 8);
    res.Add(header);

    var version = GetPythonVersion();

    var magic = new Chunk("Python version " + version, 4);
    header.Nodes.Add(magic);

    BaseStream.Position = 4;
    var dt = MTimeToTimestamp(ReadUInt32());

    var compileDate = new Chunk("Unix timestamp", 4);
    compileDate.offset = 0x04;
    compileDate.Text += " = " + dt.ToString("yyyy-MM-dd HH:mm");
    header.Nodes.Add(compileDate);

    var codeObject = new Chunk("Code object", (uint)BaseStream.Length - 8);
    codeObject.offset = 0x08;
    res.Add(codeObject);

    // TODO parse the code object

    // http://security.coverity.com/blog/2014/Nov/understanding-python-bytecode.html

    return res;
}
*/
