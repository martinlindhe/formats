package parse

/*
public CafReader(FileStream fs) : base(fs)
{
    name = "CAF audio";
    extensions = ".caf";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (BaseStream.Length < 100)
        return false;

    if (ReadByte() != 'c' || ReadByte() != 'a' || ReadByte() != 'f' || ReadByte() != 'f')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    BaseStream.Position = 0;

    if (!IsRecognized())
        throw new Exception("not a caff");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "CAFF identifier";
    res.Add(header);

    return res;
}
*/
