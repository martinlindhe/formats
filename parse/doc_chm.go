package parse

/*
public ChmReader(FileStream fs) : base(fs)
{
    name = "CHM help file (Windows)";
    extensions = ".chm";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is right magic bytes? just guessing
    if (ReadByte() != 'I' || ReadByte() != 'T' || ReadByte() != 'S' || ReadByte() != 'F')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a CHM");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "CHM identifier";
    res.Add(header);

    return res;
}
*/
