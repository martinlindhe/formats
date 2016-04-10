package parse

/*
public DocReader(FileStream fs) : base(fs)
{
    name = "MS Word document";
    extensions = ".doc";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is right magic bytes? just guessing
    if (ReadByte() != 0xD0 || ReadByte() != 0xCF || ReadByte() != 0x11 || ReadByte() != 0xE0)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a DOC");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "DOC identifier";
    res.Add(header);

    return res;
}
*/
