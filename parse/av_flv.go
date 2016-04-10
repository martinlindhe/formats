package parse

/*
public FlvReader(FileStream fs) : base(fs)
{
    name = "FLV video";
    extensions = ".flv";
    mimetype = "video/x-flv";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != 'F' || ReadByte() != 'L' || ReadByte() != 'V')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a flv");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 3;
    header.Text = "FLV identifier";
    res.Add(header);

    return res;
}
*/
