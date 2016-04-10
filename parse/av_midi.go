package parse

/*
public MidiReader(FileStream fs) : base(fs)
{
    name = "MIDI file";
    extensions = ".mid; .midi";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != 'M' || ReadByte() != 'T' || ReadByte() != 'h' || ReadByte() != 'd')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a midi");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 4;
    header.Text = "MIDI identifier";
    res.Add(header);

    return res;
}
*/
