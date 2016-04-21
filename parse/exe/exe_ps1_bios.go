package exe

/*
public PlaystationBiosReader(FileStream fs) : base(fs)
{
    name = "Playstation BIOS image";
}

override public bool IsRecognized()
{
    var identifier = Encoding.ASCII.GetBytes("Sony Computer Entertainment Inc.\0");

    if (BaseStream.Length < (0x108 + identifier.Length)) {
        return false;
    }

    BaseStream.Position = 0x108;

    var tmp = ReadBytes(identifier.Length);

    if (tmp.SequenceEqual(identifier))
        return true;

    return false;
}

override public List<Chunk> GetFileStructure()
{
    List<Chunk> res = new List<Chunk>();
    return res;
}
*/
