package parse

/*
public WriReader(FileStream fs) : base(fs)
{
    name = "WRI document (Win16)";
    extensions = ".wri";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is right magic bytes? just guessing FIXME IT IS     if data.find(b'\xBE\x00\x00\x00\xAB\x00\x00\x00\x00\x00\x00\x00\x00') == 1
    if (ReadByte() != 0x31 || ReadByte() != 0xBE || ReadByte() != 0 || ReadByte() != 0)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a ẂRI");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "WRI identifier";
    res.Add(header);

    return res;
}
*/
