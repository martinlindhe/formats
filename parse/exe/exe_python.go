package exe

// Python bytecode

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

var (
	pythonVersionMagic = map[uint32]string{
		0x00999902: "1.0",
		0x00999903: "1.1-1.2",
		0x0a0d2e89: "1.3",
		0x0a0d1704: "1.4",
		0x0a0d4e99: "1.5",
		0x0a0dc4fc: "1.6",
		0x0a0dc687: "2.0",
		0x0a0deb2a: "2.1",
		0x0a0ded2d: "2.2",
		0x0a0df23b: "2.3",
		0x0a0df26d: "2.4",
		0x0a0df2b3: "2.5",
		0x0a0df2d1: "2.6",
		0x0a0df303: "2.7",
		0x0a0d0c3a: "3.0",
		0x0a0d0c4e: "3.1",
		0x0a0d0c6c: "3.2",
		0x0a0d0c9e: "3.3",
		0x0a0d0cee: "3.4",
	}
)

func PythonBytecode(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPythonBytecode(c.File) {
		return nil, nil
	}
	return parsePythonBytecode(c.File, c.ParsedLayout)
}

func isPythonBytecode(file *os.File) bool {

	// TODO: in order to work with hdr []byte, we need to read uint32le from it ...
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

func parsePythonBytecode(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le}, // XXX decode to python version
		}}}

	return &pl, nil
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
