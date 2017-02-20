package archive

// STATUS: borked

/*
public TarReader(FileStream fs) : base(fs)
{
    name = "Tar archive";
    extensions = ".tar";
}

override public bool IsRecognized()
{
    // NOTE tar does not have a "header", instead accept all files with .tar extension for now
    string ext = Path.GetExtension(filename);
    if (ext.ToLower() != ".tar")
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a tar");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 2;
    header.Text = "TAR identifier";
    res.Add(header);

    return res;
}
*/
