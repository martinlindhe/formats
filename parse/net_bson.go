package parse

/*
public BsonReader(FileStream fs) : base(fs)
{
    name = "Bson";
}

override public bool IsRecognized()
{
    // TODO improve id method to get less false positives

    // first uint32 is total file length
    BaseStream.Position = 0;

    var dataLen = ReadUInt32();
    if (dataLen == BaseStream.Length) {
        return true;
    }

    return false;
}

 // Decodes cstring;
 // Zero or more modified UTF-8 encoded characters followed by '\x00'.
 // The (byte*) MUST NOT contain '\x00', hence it is not full UTF-8.
protected string ReadCString()
{
    var res = new StringBuilder();

    while (true) {

        var b = BaseStream.ReadByte();
        if (b == -1) {
            throw new Exception("end of file");
        }

        if (b == 0) {
            return res.ToString();
        }

        res.Append((char)b);
    }
}

 // String - The int32 is the number bytes in the (byte*) + 1 (for the trailing '\x00').
 // The (byte*) is zero or more UTF-8 encoded characters.
protected string ReadString()
{
    var res = new StringBuilder();
    var len = ReadInt32();

    // TODO does this really work with utf8 strings?

    for (var pos = 0; pos < len; pos++) {
        var x = ReadByte();
        res.Append((char)x);
    }

    return res.ToString();
}

override public List<Chunk> GetFileStructure()
{
    List<Chunk> res = new List<Chunk>();

    var header = new Chunk("Bson data", (uint)BaseStream.Length);
    res.Add(header);

    BaseStream.Position = 4;

    while (true) {
        var next = DecodeNext(header);
        if (next == false) {
            break;
        }
    }

    if (BaseStream.Length - BaseStream.Position != 1) {
        Console.WriteLine("ERROR something failed");
        return res;
    }

    var last = new Chunk();
    last.offset = BaseStream.Position;
    last.length = 1;
    last.Text = "end of document marker";

    var lastByte = ReadByte();
    if (lastByte != 0) {
        Console.WriteLine("ERROR last byte is not 0");
    }

    header.Nodes.Add(last);

    return res;
}

private bool DecodeNext(Chunk parent)
{
    if (BaseStream.Position >= parent.offset + parent.length - 1) {
        Console.WriteLine("reached end of block");
        return false;
    }

    var next = new Chunk();
    next.offset = BaseStream.Position;

    var b = ReadByte();

    var chunkType = new Chunk {
        Text = "type: ",
        offset = next.offset,
        length = 1
    };
    next.Nodes.Add(chunkType);

    var str = ReadCString();

    var chunkStr = new Chunk {
        Text = "name: " + str,
        offset = next.offset + 1,
        length = (uint)str.Length + 1
    };
    next.Nodes.Add(chunkStr);


    if (b == 0x01) {
        // floating point

        var val = ReadDouble();

        chunkType.Text += "float";

        next.Text = "float: " + str + " = " + val;
        next.length = 1 + (uint)str.Length + 1 + 8;
        parent.Nodes.Add(next);

        var chunkData = new Chunk {
            Text = "value: " + val.ToString(),
            offset = next.offset + 1 + (uint)str.Length + 1,
            length = 8
        };
        next.Nodes.Add(chunkData);

        return true;

    } else if (b == 0x02) {
        // UTF-8 string;
        var val = ReadString();

        chunkType.Text += "string";

        next.Text = "string: " + str + " = " + val;
        next.length = 1 + (uint)str.Length + 1 + 4 + (uint)val.Length;
        parent.Nodes.Add(next);

        var chunkData = new Chunk {
            Text = "value: " + val,
            offset = next.offset + 1 + (uint)str.Length + 1,
            length = 4 + (uint)val.Length
        };
        next.Nodes.Add(chunkData);

        return true;

    } else if (b == 0x04) {
        // array
        var len = (uint)ReadInt32();

        chunkType.Text += "array";

        next.Text = "array: " + str;
        next.length = 1 + (uint)str.Length + 1 + len;
        parent.Nodes.Add(next);

        // loop thru content and add to current node
        while (BaseStream.Position < next.offset + next.length - 1) {
            //Console.WriteLine("pos = " + BaseStream.Position.ToString("x2") + ", max pos = " + (next.offset + next.length).ToString("x2"));
            var sub = DecodeNext(next);
            if (sub == false) {
                Console.WriteLine("XXX ERROR - too much read in array decoder");
                return false;
            }
        }

        var last = new Chunk();
        last.offset = BaseStream.Position;
        last.length = 1;
        last.Text = "end of array marker";

        var lastByte = ReadByte();
        if (lastByte != 0) {
            Console.WriteLine("ERROR last byte is not 0");
        }

        next.Nodes.Add(last);

        return true;
    } else if (b == 0x08) {
        // boolean
        var val = ReadByte();

        chunkType.Text += "bool";

        next.Text = "bool: " + str + " = " + (val == 1 ? "true" : "false");
        next.length = 1 + (uint)str.Length + 1 + 1;
        parent.Nodes.Add(next);

        var chunkData = new Chunk {
            Text = "value: " + val,
            offset = next.offset + 1 + (uint)str.Length + 1,
            length = 1
        };
        next.Nodes.Add(chunkData);

        return true;
    } else if (b == 0x0A) {
        // null

        chunkType.Text += "null";

        next.Text = "null: " + str;
        next.length = 1 + (uint)str.Length + 1;
        parent.Nodes.Add(next);
        return true;
    } else if (b == 0x12) {
        // int64
        var val = ReadInt64();

        chunkType.Text += "int64";

        next.Text = "int64: " + str + " = " + val;
        next.length = 1 + (uint)str.Length + 1 + 8;
        parent.Nodes.Add(next);

        var chunkData = new Chunk {
            Text = "value: " + val,
            offset = next.offset + 1 + (uint)str.Length + 1,
            length = 8
        };
        next.Nodes.Add(chunkData);

        return true;
    } else {
        Console.WriteLine("unhandled at " + BaseStream.Position.ToString("x2") + ", val " + b.ToString("x2"));
        return false;
    }
}
*/
