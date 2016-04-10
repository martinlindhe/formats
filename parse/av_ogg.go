package parse

/*
public OggReader(FileStream fs) : base(fs)
{
    name = "OGG container";
    extensions = ".ogg; .oga; .ogv";
}

override public bool IsRecognized()
{
    if (BaseStream.Length < 100)
        return false;

    BaseStream.Position = 0;

    if (ReadByte() != 'O' || ReadByte() != 'g' || ReadByte() != 'g')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a ogg");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 3;
    header.Text = "OGG identifier";
    res.Add(header);

    return res;
}
*/
