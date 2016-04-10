package parse

/*
public XzReader(FileStream fs) : base(fs)
{
    name = "XZ archive";
    extensions = ".xz";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    if (ReadByte() != 0xFD || ReadByte() != '7' || ReadByte() != 'z' || ReadByte() != 'X' || ReadByte() != 'Z' || ReadByte() != 0x00)
        return false;

    return true;
}

string DecodeFlagsValue(ushort flags)
{
    if (flags == 0x0000)
        return "None";

    if (flags == 0x0100)
        return "CRC32";

    if (flags == 0x0400)
        return "CRC64";

    if (flags == 0x0A00)
        return "SHA-256";

    return "Unknown " + flags.ToString("x4");
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a xz");

    List<Chunk> res = new List<Chunk>();

    var identifier = new Chunk();
    identifier.offset = 0;
    identifier.length = 6;
    identifier.Text = "XZ identifier";
    res.Add(identifier);

    var flags = identifier.RelativeToLittleEndian16("Flags");
    var flagsValue = flags.GetValue(BaseStream);

    flags.Text += " = " + DecodeFlagsValue(flagsValue);

    res.Add(flags);

    var crc32 = flags.RelativeToLittleEndian32("CRC32");
    res.Add(crc32);

    return res;
}
*/
