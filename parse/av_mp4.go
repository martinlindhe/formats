package parse

/*
public Mp4Reader(FileStream fs) : base(fs)
{
    name = "MP4 audio";
    extensions = ".mp4; .m4a; .m4r";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // TODO what is right magic bytes? just guessing
    if (ReadByte() != 0 || ReadByte() != 0 || ReadByte() != 0 || ReadByte() != 0x18)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a mp4");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "MP4 identifier";
    res.Add(header);

    return res;
}
*/
