package parse

/*
public PifReader(FileStream fs) : base(fs)
{
    name = "Windows PIF file";
    extensions = ".pif";
}

override public bool IsRecognized()
{
    string ext = Path.GetExtension(filename);
    if (ext.ToLower() != ".pif")
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a pif");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 2;
    header.Text = "PIF identifier";
    res.Add(header);

    return res;
}
*/
