package parse

/*
public MkvReader(FileStream fs) : base(fs)
{
    name = "MKV container";
    extensions = ".mkv";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is magic sequence?
    if (ReadByte() != 0x1a || ReadByte() != 0x45 || ReadByte() != 0xdf || ReadByte() != 0xa3)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a mkv");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "MKV identifier";
    res.Add(header);

    return res;
}
*/
