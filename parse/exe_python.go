package parse

/*
public PythonBytecodeReader(FileStream fs) : base(fs)
{
    name = "Python bytecode";
}

public string GetPythonVersion()
{
    BaseStream.Position = 0;

    var magic = ReadUInt32();

    if (magic == 0x00999902) {
        return "1.0";
    }
    if (magic == 0x00999903) {
        return "1.1-1.2";
    }
    if (magic == 0x0A0D2E89) {
        return "1.3";
    }
    if (magic == 0x0A0D1704) {
        return "1.4";
    }
    if (magic == 0x0A0D4E99) {
        return "1.5";
    }
    if (magic == 0x0A0DC4FC) {
        return "1.6";
    }
    if (magic == 0x0A0DC687) {
        return "2.0";
    }
    if (magic == 0x0A0DEB2A) {
        return "2.1";
    }
    if (magic == 0x0A0DED2D) {
        return "2.2";
    }
    if (magic == 0x0A0DF23B) {
        return "2.3";
    }
    if (magic == 0x0A0DF26D) {
        return "2.4";
    }
    if (magic == 0x0A0DF2B3) {
        return "2.5";
    }
    if (magic == 0x0A0DF2D1) {
        return "2.6";
    }
    if (magic == 0x0A0DF303) {
        return "2.7";
    }
    if (magic == 0x0A0D0C3A) {
        return "3.0";
    }
    if (magic == 0x0A0D0C4E) {
        return "3.1";
    }
    if (magic == 0x0A0D0C6C) {
        return "3.2";
    }
    if (magic == 0x0A0D0C9E) {
        return "3.3";
    }
    if (magic == 0x0A0D0CEE) {
        return "3.4";
    }
    return null;
}

override public bool IsRecognized()
{
    var ver = GetPythonVersion();
    if (ver == null) {
        return false;
    }
    return true;
}

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
