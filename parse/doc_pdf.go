package parse

/*
public PdfReader(FileStream fs) : base(fs)
{
    name = "Portable Document Format";
    extensions = ".pdf";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != '%' || ReadByte() != 'P' || ReadByte() != 'D' || ReadByte() != 'F')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a pdf");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "PDF identifier";
    res.Add(header);

    return res;
}
*/
