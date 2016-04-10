package parse

/*
// TODO detect audio ( WAVE) or video format

public RiffReader(FileStream fs) : base(fs)
{
    name = "RIFF format (WAV, AVI)";
    extensions = ".wav; .avi";
}

override public bool IsRecognized()
{
    if (BaseStream.Length < 100)
        return false;

    BaseStream.Position = 0;

    if (ReadByte() != 'R' || ReadByte() != 'I' || ReadByte() != 'F' || ReadByte() != 'F')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a riff");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "WAV identifier";
    res.Add(header);

    return res;
}
*/
