package bin

/*

// TODO: parse chunk lengths
// TODO: parse TRX v2, need sample!

public TrxReader(FileStream fs) : base(fs)
{
    name = "TRX firmware";
    extensions = ".trx";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;
    if (ReadByte() != 'H' || ReadByte() != 'D' || ReadByte() != 'R' || ReadByte() != '0')
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a trx");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk("TRX header", 28);
    header.offset = 0;
    res.Add(header);

    var identifier = new Chunk();
    identifier.offset = 0;
    identifier.length = 4;
    identifier.Text = "TRX identifier";
    header.Nodes.Add(identifier);

    var length = identifier.RelativeToLittleEndian32("Length");
    var lengthValue = length.GetValue(BaseStream);

    if (lengthValue != filesize) {
        Log("header Lenght value is wrong, it is " + lengthValue.ToString("x8") + " but should be " + filesize.ToString("x8"));
    }

    header.Nodes.Add(length);

    var crc = length.RelativeToLittleEndian32("CRC");
    header.Nodes.Add(crc);

    var flags = crc.RelativeToLittleEndian16("TRX flags");
    header.Nodes.Add(flags);

    var version = flags.RelativeToLittleEndian16("Version");
    var versionValue = version.GetValue(BaseStream);
    header.Nodes.Add(version);

    var offset0 = version.RelativeToLittleEndian32("Partition offset 0");
    var offset0value = offset0.GetValue(BaseStream);
    header.Nodes.Add(offset0);

    var offset1 = offset0.RelativeToLittleEndian32("Partition offset 1");
    var offset1value = offset1.GetValue(BaseStream);
    header.Nodes.Add(offset1);

    var offset2 = offset1.RelativeToLittleEndian32("Partition offset 2");
    var offset2value = offset2.GetValue(BaseStream);
    header.Nodes.Add(offset2);

    // Log("version = " + versionValue);

    if (versionValue > 1)
        Log("SAMPLE PLEASE! - TRX v2 parsing is untested");

    // TODO: how to get length of partition ?

    if (offset0value > 0) {
        var chunk0 = new Chunk("Partition 0 data - lzma-loader");
        chunk0.length = 10; // XXX FIX
        chunk0.offset = offset0value;
        res.Add(chunk0);
    }

    if (offset1value > 0) {
        var chunk1 = new Chunk("Partition 1 data - Linux kernel (squashfs?)");
        chunk1.length = 10; // XXX FIX
        chunk1.offset = offset1value;
        res.Add(chunk1);
    }

    if (offset2value > 0) {
        var chunk2 = new Chunk("Partition 2 data - rootfs");
        chunk2.length = 10; // XXX FIX
        chunk2.offset = offset2value;
        res.Add(chunk2);
    }

    // v2 only
    if (versionValue > 1) {
        var offset3 = offset2.RelativeToLittleEndian32("Partition offset 3");
        var offset3value = offset3.GetValue(BaseStream);
        header.Nodes.Add(offset3);

        if (offset3value > 0) {
            var chunk3 = new Chunk("Partition 3 data - bin-Header");
            chunk3.length = 10; // XXX FIX
            chunk3.offset = offset3value;
            res.Add(chunk3);
        }
    }

    return res;
}
*/
