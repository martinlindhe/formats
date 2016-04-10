package parse

/*
public OtfReader(FileStream fs) : base(fs)
{
    name = "OpenType Font";
    extensions = ".otf";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != 'O' || ReadByte() != 'T' || ReadByte() != 'T' || ReadByte() != 'O')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a otf");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "OTF identifier";
    res.Add(header);

    return res;
}
*/
