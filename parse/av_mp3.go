package parse

/*
public Mp3Reader(FileStream fs) : base(fs)
{
    name = "MP3 audio";
    extensions = ".mp3";
}

override public bool IsRecognized()
{
    if (BaseStream.Length < 100)
        return false;

    BaseStream.Position = 0;

    // TODO find mp3 stream start, ignore id3 tags

    if (ReadByte() != 'I' || ReadByte() != 'D' || ReadByte() != '3')
        return false;

    byte id3ver = ReadByte();
    if (id3ver != 3 && id3ver != 4)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a mp3");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 30; // XXX not right
    header.Text = "ID3 header";
    res.Add(header);

    var tag = new Chunk("ID3 identifier", 4);
    tag.offset = header.offset;
    header.Nodes.Add(tag);

    if (ReadByte(3) == 3) {
        tag.Text = "ID3 v3.3 identifier";
    } else if (ReadByte(3) == 4) {
        tag.Text = "ID3 v3.4 identifier";
    } else {
        throw new Exception("unknown id3 version");
    }

    var minorVer = tag.RelativeToByte("Minor version"); // usually 0x00
    header.Nodes.Add(minorVer);


    var id3Flags = minorVer.RelativeToByte("ID3 flags");
    byte id3FlagsValue = ReadByte(id3Flags.offset);
    header.Nodes.Add(id3Flags);

    var id3Size = id3Flags.RelativeToBigEndian32("ID3 size");
    var id3SizeValue = (uint)ReadInt32(id3Size.offset);
    header.Nodes.Add(id3Size);

    if (id3FlagsValue > 0) {
        // The ID3v2 tag size is the sum of the byte length of the extended
        // header, the padding and the frames after unsynchronisation. If a
        // footer is present this equals to ('total size' - 20) bytes, otherwise
        // ('total size' - 10) bytes.

        //TODO: id3 tag can have "extended header" etc
        //throw new Exception("ERROR unhandled id3 flags");

        // bit 7 =  indicates whether or not  unsynchronisation is applied on all frames (see section 6.1 for details)
        if ((id3FlagsValue & 0x80) > 0)
            Log("unsynchronisation");

        // bit 6 -  Extended header
        if ((id3FlagsValue & 0x40) > 0)
            Log("extended header"); // TODO parse

        // bit 5 - experimental indicator
        if ((id3FlagsValue & 0x20) > 0)
            Log("Experimental indicator");

        // bit 4 - footer present
        if ((id3FlagsValue & 0x10) > 0)
            Log("Footer present");

    }

    var Data = id3Size.RelativeTo("Data", id3SizeValue);
    header.Nodes.Add(Data);


    //uint32_t mpeg_header = read32be(f);
    //if (mpeg_header && 0xFFF00000 == 0xFFF00000) {
    //    if (!info) {
    //        printf("audio/mpeg\n"); //FF2 = 'audio/x-mpeg', IE7 = 'audio/mpeg'
    //    } else {
    //        printf("Format: MP3\n");
    //        printf("Mediatype: audio\n");
    //        printf("Mimetype: audio/mpeg\n");
    //    }
    //    return E_PROBESUCCESS;
    //}

    //printf("cant find mpeg frame start marker: %08x\n", mpeg_header);    //FFFA9344

    return res;
}
*/
