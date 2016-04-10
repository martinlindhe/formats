package parse

/*
// TODO  video/x-ms-wmv, audio/x-ms-wma

public AsfReader(FileStream fs) : base(fs)
{
    name = "ASF container";
    extensions = ".asf; .wmv; .wma";
}

public static readonly byte[] ObjectSignature = new byte[] {
    0x30, 0x26, 0xB2, 0x75, 0x8E, 0x66, 0xCF, 0x11, 0xA6, 0xD9, 0x00, 0xAA, 0x00, 0x62, 0xCE, 0x6C
};
public static readonly byte[] ObjectStreamProperties = new byte[] {
    0x91, 0x07, 0xDC, 0xB7, 0xB7, 0xA9, 0xCF, 0x11, 0x8E, 0xE6, 0x00, 0xC0, 0x0C, 0x20, 0x53, 0x65
};
public static readonly byte[] StreamPropertyAudio = new byte[] {
    0x40, 0x9E, 0x69, 0xF8, 0x4D, 0x5B, 0xCF, 0x11, 0xA8, 0xFD, 0x00, 0x80, 0x5F, 0x5C, 0x44, 0x2B
};
public static readonly byte[] StreamPropertyVideo = new byte[] {
    0xC0, 0xEF, 0x19, 0xBC, 0x4D, 0x5B, 0xCF, 0x11, 0xA8, 0xFD, 0x00, 0x80, 0x5F, 0x5C, 0x44, 0x2B
};

public bool HasSignature(long offset, byte[] sig)
{
    BaseStream.Position = offset;

    for (int i = 0; i < sig.Length; i++)
        if (ReadByte() != sig[i])
            return false;

    return true;
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // FIXME what is minimal size?
    if (BaseStream.Length < 16 * 4)
        return false;

    if (!HasSignature(0, ObjectSignature))
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a ASF");

    var res = new List<Chunk>();

    var header = ParseAsfHeader();
    res.Add(header);

    return res;
}

private Chunk ParseAsfHeader()
{
    var header = new Chunk("ASF header");
    header.offset = 0;
    header.length = 16 + 14;

    var signature = new Chunk("ASF signature");
    signature.offset = header.offset;
    signature.length = 16;

    header.Nodes.Add(signature);

    var size = signature.RelativeToLittleEndian64("Size");
    header.Nodes.Add(size);

    var Objects = size.RelativeToLittleEndian32("Objects");
    var ObjectsValue = ReadInt32(Objects.offset);
    header.Nodes.Add(Objects);

    var Reserved1 = Objects.RelativeToByte("Reserved 1");
    header.Nodes.Add(Reserved1);

    var Reserved2 = Reserved1.RelativeToByte("Reserved 2");
    header.Nodes.Add(Reserved2);

    Log("Parsing " + ObjectsValue + " objects");

    long offset = Reserved2.offset + Reserved2.length;
    for (int i = 0; i < ObjectsValue; i++) {
        var subHead = new Chunk("Object # " + (i + 1));
        subHead.offset = offset;
        header.Nodes.Add(subHead);

        var guid = new Chunk("GUID", 16);
        guid.offset = subHead.offset;
        subHead.Nodes.Add(guid);

        // string hex = ByteArrayToString(d, guid.offset, 16);
        string hex = "XXX TODO";
        Log("Object guid = " + hex);

        var len = guid.RelativeToLittleEndian64("Length");
        var lenValue = ReadInt64(len.offset);
        subHead.Nodes.Add(len);

        var Data = len.RelativeTo("Data", (uint)lenValue - (guid.length + len.length));
        subHead.length = (uint)lenValue;
        subHead.Nodes.Add(Data);

        if (HasSignature(guid.offset, ObjectStreamProperties)) {
            // TODO parse remaining of stream properites object
            var streamGuid = len.RelativeTo("Stream prop GUID", 16);
            subHead.Nodes.Add(streamGuid);

            //string streamHex = ByteArrayToString(d, streamGuid.offset, 16);
            string streamHex = "XXX streamHex";

            if (HasSignature(streamGuid.offset, StreamPropertyAudio)) {
                Log("Audio");
            } else if (HasSignature(streamGuid.offset, StreamPropertyVideo)) {
                Log("Video");
            } else {
                Log("Unknown stream props guid = " + streamHex);
            }
        }

        offset += (int)lenValue;
    }

    return header;
}

public static string ByteArrayToString(byte[] ba, long offset, int length)
{
    StringBuilder hex = new StringBuilder(ba.Length * 2);

    int count = 0;

    for (long i = offset; i < ba.Length; i++) {
        hex.Append(" ");
        hex.AppendFormat("{0:x2}", ba[i]);
        count++;
        if (count >= length)
            break;
    }
    return hex.ToString().Trim();
}
*/
