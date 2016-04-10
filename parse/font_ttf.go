package parse

/*
public TtfReader(FileStream fs) : base(fs)
{
    name = "TrueType Font";
    extensions = ".ttf";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != 0 || ReadByte() != 1 || ReadByte() != 0 || ReadByte() != 0)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a ttf");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 5;
    header.Text = "TTF identifier";
    res.Add(header);

    return res;
}
*/
