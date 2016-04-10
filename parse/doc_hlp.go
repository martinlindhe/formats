package parse

/*
public HlpReader(FileStream fs) : base(fs)
{
    name = "HLP help file (Windows)";
    extensions = ".hlp";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is right magic bytes? just guessing
    if (ReadByte() != 0x3F || ReadByte() != 0x5F || ReadByte() != 3 || ReadByte() != 0)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a HLP");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "HLP identifier";
    res.Add(header);

    return res;
}
*/
