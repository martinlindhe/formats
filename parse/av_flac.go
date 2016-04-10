package parse

/*
public FlacReader(FileStream fs) : base(fs)
{
    name = "FLAC audio";
    extensions = ".fla; .flac";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (BaseStream.Length < 100)
        return false;

    if (ReadByte() != 'f' || ReadByte() != 'L' || ReadByte() != 'a' || ReadByte() != 'C')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a flac");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "FLAC identifier";
    res.Add(header);

    return res;
}
*/
