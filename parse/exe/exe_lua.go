package exe

// Lua bytecode
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func LUA(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isLUA(file) {
		return nil, nil
	}
	return parseLUA(file, pl)
}

func isLUA(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b == 0x61754c1b {
		// Lua 5.1 and 5.2 identifer
		return true
	}

	return false
}

func parseLUA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}

/*
override public List<Chunk> GetFileStructure()
{
    // NOTE: first 12 bytes are same for lua 5.1 and lua 5.2

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk{ Text = "Lua header", length = 12 };
    res.Add(header);


    var identifier = new Chunk{ Text = "Lua identifier", length = 4 };
    header.Nodes.Add(identifier);

    BaseStream.Position = 4;

    var versionNumber = ReadByte();
    if (versionNumber >= 0x52) {
        header.length = 18;
    }

    // TODO need samples for more versions, like lua 5.0
    if (versionNumber != 0x51 && versionNumber != 0x52) {
        Console.WriteLine("unknown LUA version " + versionNumber + ", sample please!");
        return res;
    }

    var version = new Chunk("Version = " + versionNumber.ToString("x2"), 1);  // TODO decode version number
    version.offset = 4;
    header.Nodes.Add(version);

    var officialCode = ReadByte();
    var official = version.RelativeTo("Official = " + (officialCode == 0 ? "yes" : "no"), 1);
    header.Nodes.Add(official);

    // TODO decode meaning of system params
    var systemParam = official.RelativeTo("System params", 6);
    header.Nodes.Add(systemParam);

    if (versionNumber >= 0x52) {
        var conversionErr = systemParam.RelativeTo("data to catch conversation errors", 6);
        header.Nodes.Add(conversionErr);
    }

    // TODO after header comes: function definitions, then opcodes, then constants, then function prototypes, upvalues, debug info,

    return res;
}
*/
