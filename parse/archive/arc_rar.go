package archive

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RAR(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isRAR(c.Header) {
		return nil, nil
	}
	return parseRAR(c.File, c.ParsedLayout)
}

func isRAR(b []byte) bool {

	if b[0] != 'R' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}

	// RAR 4.x signature
	//if (ReadByte() != 0x1A || ReadByte() != 0x07 || ReadByte() != 0x00)
	//    return false;

	// RAR 5.0 signature
	//if (ReadByte() != 0x1A || ReadByte() != 0x07 || ReadByte() != 0x01 || ReadByte() != 0x00)
	//    return false;

	return true
}

func parseRAR(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}

/*
// NOTE: internal field naming from unrar sources (headers.hpp, arcread.cpp)
// TODO: support RAR 5.0 format, need samples
// RAR 5.0 header types.
// HEAD_MARK=0x00, HEAD_MAIN=0x01, HEAD_FILE=0x02, HEAD_SERVICE=0x03, HEAD_CRYPT=0x04, HEAD_ENDARC=0x05, HEAD_UNKNOWN=0xff,

class RarVolumeHeader
{
    public ushort crc;
    public byte type;
    public ushort flags;
    public ushort size;
}


void DecodeMainHeaderFlags(ushort flags)
{
    Log("main head flags = 0x" + flags.ToString("x4"));

    if ((flags & 0x0001) != 0)
        Log("MHD_VOLUME");
    if ((flags & 0x0002) != 0)
        Log("MHD_COMMENT");
    if ((flags & 0x0004) != 0)
        Log("MHD_LOCK");
    if ((flags & 0x0008) != 0)
        Log("MHD_SOLID");
    if ((flags & 0x0010) != 0)
        Log("MHD_PACK_COMMENT or MHD_NEWNUMBERING");
    if ((flags & 0x0020) != 0)
        Log("MHD_AV");
    if ((flags & 0x0040) != 0)
        Log("MHD_PROTECT");
    if ((flags & 0x0080) != 0)
        Log("MHD_PASSWORD");
    if ((flags & 0x0100) != 0)
        Log("MHD_FIRSTVOLUME");
    if ((flags & 0x0200) != 0)
        Log("MHD_ENCRYPTVER");
}

void DecodeFileHeaderFlags(ushort flags)
{
    Log("file head flags = 0x" + flags.ToString("x4"));

    if ((flags & 0x0001) != 0)
        Log("LHD_SPLIT_BEFORE");
    if ((flags & 0x0002) != 0)
        Log("LHD_SPLIT_AFTER");
    if ((flags & 0x0004) != 0)
        Log("LHD_PASSWORD");
    if ((flags & 0x0008) != 0)
        Log("LHD_COMMENT");
    if ((flags & 0x0010) != 0)
        Log("LHD_SOLID");
    if ((flags & 0x0100) != 0)
        Log("LHD_LARGE");
    if ((flags & 0x0200) != 0)
        Log("LHD_UNICODE");
    if ((flags & 0x0400) != 0)
        Log("LHD_SALT");
    if ((flags & 0x0800) != 0)
        Log("LHD_VERSION");
    if ((flags & 0x1000) != 0)
        Log("LHD_EXTTIME");
    if ((flags & 0x2000) != 0)
        Log("LHD_EXTFLAGS");
}

private string DecodeHostOs(byte b)
{
    if (b == 0)
        return "MS DOS";
    if (b == 1)
        return "OS/2";
    if (b == 2)
        return "Win32";
    if (b == 3)
        return "Unix";
    if (b == 4)
        return "Mac OS";
    if (b == 5)
        return "BeOS";

    return "Unknown";
}

private string DecodeMethod(byte b)
{
    if (b == 0x30)
        return "storing";
    if (b == 0x31)
        return "fastest compression";
    if (b == 0x32)
        return "fast compression";
    if (b == 0x33)
        return "normal compression";
    if (b == 0x34)
        return "good compression";
    if (b == 0x35)
        return "best compression";

    return "Unknown";
}

Chunk ParseExtTime(long baseOffset)
{
    var ExtTime = new Chunk("ExtTime");
    ExtTime.offset = baseOffset;
    ExtTime.length = 2;

    var ExtTimeFlags = new LittleEndian16BitChunk("ExtTime Flags");
    ExtTimeFlags.offset = baseOffset;
    var ExtTimeFlagsValue = ExtTimeFlags.GetValue(BaseStream);
    ExtTime.Nodes.Add(ExtTimeFlags);

    var offset = BaseStream.Position;


    uint rmode = (uint)(ExtTimeFlagsValue >> 12);
    uint count = 0;

    // FIXME verify that other than mtype is decoded properly, need samples

    if ((rmode & 8) != 0) {
        count = rmode & 0x3;
        Log("mtime_count = " + count);

        var mtime = new Chunk("mtime", count);
        mtime.offset = offset;
        ExtTime.Nodes.Add(mtime);

        ExtTime.length += count;
        offset += count;
    }

    rmode = (uint)(ExtTimeFlagsValue >> 8);
    if ((rmode & 8) != 0) {
        //  Set ctime to readBits(16)
        count = rmode & 0x3;
        Log("ctime_count = " + count);
        //var ctime_reminder = readBytes(ctime_count);

        var ctime = new Chunk("ctime", count);
        ctime.offset = offset;
        ExtTime.Nodes.Add(ctime);

        ExtTime.length += count;
        offset += count;
    }

    rmode = (uint)(ExtTimeFlagsValue >> 4);
    if ((rmode & 8) != 0) {
        // Set atime to readBits(16)
        count = rmode & 0x3;
        Log("atime_count = " + count);
        //var atime_reminder = readBytes(atime_count);

        var atime = new Chunk("atime", count);
        atime.offset = offset;
        ExtTime.Nodes.Add(atime);

        ExtTime.length += count;
        offset += count;
    }

    rmode = ExtTimeFlagsValue;
    if ((rmode & 8) != 0) {
        // Set arctime to readBits(16)
        count = rmode & 0x3;
        Log("arctime_count = " + count);
        //var arctime_reminder = readBytes(arctime_count);

        var arctime = new Chunk("arctime", count);
        arctime.offset = offset;
        ExtTime.Nodes.Add(arctime);

        ExtTime.length += count;
        offset += count;
    }

    if (count == 0)
        throw new Exception("sample please, time_count == 0");

    return ExtTime;
}

// Parses 4.x RAR header
private Chunk ParseVolumeHeader(uint baseOffset)
{
    BaseStream.Position = baseOffset;

    var volhdr = new RarVolumeHeader();
    volhdr.crc = ReadUInt16();
    volhdr.type = ReadByte();
    volhdr.flags = ReadUInt16();
    volhdr.size = ReadUInt16();

    var chunk = new Chunk("Volume header");
    chunk.offset = baseOffset;
    chunk.length = volhdr.size;

    var crc = new LittleEndian16BitChunk("Crc");
    crc.offset = chunk.offset;
    chunk.Nodes.Add(crc);

    var type = crc.RelativeToByte("Type = 0x" + volhdr.type.ToString("x2"));
    chunk.Nodes.Add(type);

    var flags = type.RelativeToLittleEndian16("Flags");
    chunk.Nodes.Add(flags);

    var size = flags.RelativeToLittleEndian16("Size");
    chunk.Nodes.Add(size);

    if (volhdr.type == 0x72) {
        chunk.Text = "Volume HEAD3_MARK";
    } else if (volhdr.type == 0x73) {
        chunk.Text = "Volume HEAD3_MAIN";

        DecodeMainHeaderFlags(volhdr.flags);

        var HighPosAv = size.RelativeToLittleEndian16("HighPosAv");
        chunk.Nodes.Add(HighPosAv);

        var PosAv = HighPosAv.RelativeToLittleEndian32("PosAv");
        chunk.Nodes.Add(PosAv);

        //TODO EncryptVer              1 byte (only present if MHD_ENCRYPTVER is set)

        if ((volhdr.flags & 0x0080) != 0) {
            // MHD_PASSWORD

            // XXX is it always 28 bytes?
            var Password = PosAv.RelativeTo("Password", 28);
            chunk.Nodes.Add(Password);

            chunk.length += Password.length;
        }

    } else if (volhdr.type == 0x74) {
        chunk.Text = "Volume HEAD3_FILE";

        var PackSize = size.RelativeToLittleEndian32("PackSize");
        var PackSizeValue = PackSize.GetValue(BaseStream);
        chunk.Nodes.Add(PackSize);

        chunk.length += PackSizeValue;

        DecodeFileHeaderFlags(volhdr.flags);

        var UnpSize = PackSize.RelativeToLittleEndian32("UnpSize");
        chunk.Nodes.Add(UnpSize);

        var HostOs = UnpSize.RelativeToByte("Host OS");
        var HostOsValue = HostOs.GetValue(BaseStream);
        HostOs.Text += " = " + DecodeHostOs(HostOsValue);
        chunk.Nodes.Add(HostOs);

        var FileCRC = HostOs.RelativeToLittleEndian32("FileCRC");
        chunk.Nodes.Add(FileCRC);

        var FileTime = FileCRC.RelativeToLittleEndian32("FileTime");
        chunk.Nodes.Add(FileTime);

        var UnpVer = FileTime.RelativeToByte("UnpVer");
        chunk.Nodes.Add(UnpVer);

        var Method = UnpVer.RelativeToByte("Method");
        var MethodValue = Method.GetValue(BaseStream);

        Method.Text += " = " + DecodeMethod(MethodValue);
        chunk.Nodes.Add(Method);

        var NameSize = Method.RelativeToLittleEndian16("NameSize");
        var NameSizeValue = NameSize.GetValue(BaseStream);

        chunk.Nodes.Add(NameSize);

        var FileAttr = NameSize.RelativeToLittleEndian32("FileAttr");
        chunk.Nodes.Add(FileAttr);

        var FileName = FileAttr.RelativeTo("FileName", NameSizeValue);
        chunk.Nodes.Add(FileName);


        long offset = FileName.offset + FileName.length;

        if ((volhdr.flags & 0x0400) != 0) {
            // LHD_SALT

            var salt = new Chunk("Salt", 8);
            salt.offset = offset;
            offset += salt.length;
            chunk.Nodes.Add(salt);
        }

        if ((volhdr.flags & 0x1000) != 0) {
            // LHD_EXTTIME

            var ExtTime = ParseExtTime(offset);
            offset += ExtTime.length;
            chunk.Nodes.Add(ExtTime);
        }

        var Data = new Chunk("Data", PackSizeValue);
        Data.offset = offset;
        chunk.Nodes.Add(Data);

        //    HighPackSize            4 bytes (only present if LHD_LARGE is set)
        //    HighUnpSize             4 bytes (only present if LHD_LARGE is set)

        // TODO If the LHD_LARGE flag is set, then the archive is large and 64-bits are needed to
        // represent the packed and unpacked size. HighPackSize is used as the upper 32-bits and
        // PackSize is used as the lower 32-bits for the packed size in bytes. HighUnpSize is used
        // as the upper 32-bits and UnpSize is used as the lower 32-bits for the unpacked size in bytes.

    } else if (volhdr.type == 0x75) {
        chunk.Text = "Volume HEAD3_CMT - TODO";

          //CommHead.UnpSize=Raw.Get2();
          //CommHead.UnpVer=Raw.Get1();
          //CommHead.Method=Raw.Get1();
          //CommHead.CommCRC=Raw.Get2();

    } else if (volhdr.type == 0x76) {
        chunk.Text = "Volume HEAD3_AV - TODO";

        //AVHead.UnpVer=Raw.Get1();
        //AVHead.Method=Raw.Get1();
        //AVHead.AVVer=Raw.Get1();
        //AVHead.AVInfoCRC=Raw.Get4();

    } else if (volhdr.type == 0x77) {
        chunk.Text = "Volume HEAD3_OLDSERVICE - TODO";

    } else if (volhdr.type == 0x78) {
        chunk.Text = "Volume HEAD3_PROTECT - WIP";

        //ProtectHead.DataSize=Raw.Get4();
        //ProtectHead.Version=Raw.Get1();
        //ProtectHead.RecSectors=Raw.Get2();
        //ProtectHead.TotalBlocks=Raw.Get4();
        //Raw.GetB(ProtectHead.Mark,8);
        //NextBlockPos+=ProtectHead.DataSize;
        //RecoverySize=ProtectHead.RecSectors*512;

        var DataSize = size.RelativeToLittleEndian32("DataSize");
        chunk.Nodes.Add(DataSize);

        var Version = DataSize.RelativeToByte("Version");
        chunk.Nodes.Add(Version);

        var RecSectors = Version.RelativeToLittleEndian16("RecSectors");
        chunk.Nodes.Add(RecSectors);

        var TotalBlocks = RecSectors.RelativeToLittleEndian32("TotalBlocks");
        chunk.Nodes.Add(TotalBlocks);

    } else if (volhdr.type == 0x79) {
        chunk.Text = "Volume HEAD3_SIGN - TODO";

          //SignHead.CreationTime=Raw.Get4();
          //SignHead.ArcNameSize=Raw.Get2();
          //SignHead.UserNameSize=Raw.Get2();

    } else if (volhdr.type == 0x7A) {
        chunk.Text = "Volume HEAD3_SERVICE - TODO";

    } else if (volhdr.type == 0x7B) {
        chunk.Text = "Volume HEAD3_ENDARC - WIP";

        //EndArcHead.NextVolume=(EndArcHead.Flags & EARC_NEXT_VOLUME)!=0;
        //EndArcHead.DataCRC=(EndArcHead.Flags & EARC_DATACRC)!=0;
        //EndArcHead.RevSpace=(EndArcHead.Flags & EARC_REVSPACE)!=0;
        //EndArcHead.StoreVolNumber=(EndArcHead.Flags & EARC_VOLNUMBER)!=0;
        //if (EndArcHead.DataCRC)
        //    EndArcHead.ArcDataCRC=Raw.Get4();
        //if (EndArcHead.StoreVolNumber)
        //    VolNumber=EndArcHead.VolNumber=Raw.Get2();
    }

    return chunk;
}

override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a rar");

    List<Chunk> res = new List<Chunk>();

    uint offset = 0;
    int count = 0;

    do {
        count++;

        if (offset >= BaseStream.Length) {
            Log("Reached end of file");
            break;
        }

        var chunk = ParseVolumeHeader(offset);

        offset += chunk.length;

        res.Add(chunk);

    } while (count < 300);

    return res;
}

*/
