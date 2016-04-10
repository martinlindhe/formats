package parse

/*
public RtfReader(FileStream fs) : base(fs)
{
    name = "Rich Type File";
    extensions = ".rtf";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != '{' || ReadByte() != '\\' || ReadByte() != 'r' || ReadByte() != 't' || ReadByte() != 'f')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a rtf");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 5;
    header.Text = "RTF identifier";
    res.Add(header);

    return res;
}
*/
