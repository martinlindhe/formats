package parse

/*
public AiffReader(FileStream fs) : base(fs)
{
    name = "AIFF audio";
    extensions = ".aiff";
}

override public bool IsRecognized()
{
    // TODO also detect "AIFF" string
    BaseStream.Position = 0;

    if (ReadByte() != 'F' || ReadByte() != 'O' || ReadByte() != 'R' || ReadByte() != 'M')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a aiff");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "AIFF identifier";
    res.Add(header);

    return res;
}
*/
