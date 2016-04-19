package archive

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func ZIP(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isZIP(file) {
		return nil, nil
	}
	return parseZIP(file, pl)
}

func isZIP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [6]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'P' || b[1] != 'K' || b[2] != 3 || b[3] != 4 {
		return false
	}
	return true
}

func parseZIP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 6, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}

/*

// PkZip Host OS table
private static string DecodeHostOs(byte b)
{
    if (b == 0)
        return "MS-DOS and OS/2 (FAT / VFAT / FAT32 file systems)";
    if (b == 1)
        return "Amiga";
    if (b == 2)
        return "OpenVMS";
    if (b == 3)
        return "UNIX";
    if (b == 4)
        return "VM/CMS";
    if (b == 5)
        return "Atari ST";
    if (b == 6)
        return "OS/2 H.P.F.S.";
    if (b == 7)
        return "Macintosh";
    if (b == 8)
        return "Z-System";
    if (b == 9)
        return "CP/M";
    if (b == 10)
        return "Windows NTFS";
    if (b == 11)
        return "MVS (OS/390 - Z/OS)";
    if (b == 12)
        return "VSE";
    if (b == 13)
        return "Acorn Risc";
    if (b == 14)
        return "VFAT";
    if (b == 15)
        return "alternate MVS";
    if (b == 16)
        return "BeOS";
    if (b == 17)
        return "Tandem";
    if (b == 18)
        return "OS/400";
    if (b == 19)
        return "OS X (Darwin)";

    // 20 thru 255 - unused
    return "Unknown " + b;
}

private static string DecodeCompressionMethod(ushort b)
{
    if (b == 0)
        return "Stored (no compression)";
    if (b == 1)
        return "Shrunk";
    if (b == 2)
        return "Reduced with compression factor 1";
    if (b == 3)
        return "Reduced with compression factor 2";
    if (b == 4)
        return "Reduced with compression factor 3";
    if (b == 5)
        return "Reduced with compression factor 4";
    if (b == 6)
        return "Imploded";
    if (b == 7)
        return "Reserved for Tokenizing compression algorithm";
    if (b == 8)
        return "Deflated";
    if (b == 9)
        return "Enhanced Deflating using Deflate64(tm)";
    if (b == 10)
        return "PKWARE Data Compression Library Imploding (old IBM TERSE)";
    if (b == 11)
        return "Reserved by PKWARE";
    if (b == 12)
        return "compressed using BZIP2 algorithm";
    if (b == 13)
        return "Reserved by PKWARE";
    if (b == 14)
        return "LZMA (EFS)";
    if (b == 15)
        return " Reserved by PKWARE";
    if (b == 16)
        return "Reserved by PKWARE";
    if (b == 17)
        return "Reserved by PKWARE";
    if (b == 18)
        return "File is compressed using IBM TERSE (new)";
    if (b == 19)
        return "IBM LZ77 z Architecture (PFS)";
    if (b == 97)
        return "WavPack compressed data";
    if (b == 98)
        return "PPMd version I, Rev 1";

    return "Unknown " + b;
}

private Chunk ParseZipHeader(uint offset)
{
    var header = new Chunk();
    header.offset = offset;

    header.Text = "ZIP header";

    var identifier = new LittleEndian32BitChunk("ZIP identifier");
    identifier.offset = header.offset;
    header.Nodes.Add(identifier);

    var version = identifier.RelativeToVersionMajorMinor16("Version");
    header.Nodes.Add(version);

    var bits = version.RelativeToLittleEndian16("Bitfield");
    header.Nodes.Add(bits);

    var method = bits.RelativeToLittleEndian16("Compression method");
    var methodValue = method.GetValue(BaseStream);
    method.Text += " = " + DecodeCompressionMethod(methodValue);
    header.Nodes.Add(method);

    var timestamp = method.RelativeToLittleEndianDateStamp("DateTime");
    header.Nodes.Add(timestamp);

    var crc = timestamp.RelativeToLittleEndian32("Crc");
    header.Nodes.Add(crc);

    var csize = crc.RelativeToLittleEndian32("Compressed size");
    var csizeValue = csize.GetValue(BaseStream);
    header.Nodes.Add(csize);

    var size = csize.RelativeToLittleEndian32("Uncompressed size");
    header.Nodes.Add(size);

    var nameLen = size.RelativeToLittleEndian16("Length of filename");
    header.Nodes.Add(nameLen);

    var extraLen = nameLen.RelativeToLittleEndian16("Length of extra field");
    var extraLenValue = extraLen.GetValue(BaseStream);
    header.Nodes.Add(extraLen);

    var nameLenValue = nameLen.GetValue(BaseStream);
    var nameData = extraLen.RelativeToZeroTerminatedString("Filename", nameLenValue);
    header.Nodes.Add(nameData);

    var extraData = nameData.RelativeTo("Extra field", extraLenValue);
    if (extraData.length > 0)
        header.Nodes.Add(extraData);

    var data = extraData.RelativeTo("Data", csizeValue);
    header.Nodes.Add(data);

    header.length = (uint)((data.offset + data.length) - identifier.offset);

    return header;
}

private Chunk ParseEndOfCds(uint offset)
{
    var header = new Chunk();
    header.offset = offset;

    header.Text = "End of CDS";

    var identifier = new LittleEndian32BitChunk("End of CDS identifier");
    identifier.offset = header.offset;
    header.Nodes.Add(identifier);

    var diskNo = identifier.RelativeToLittleEndian16("Disk number");
    header.Nodes.Add(diskNo);

    var startDiskNo = diskNo.RelativeToLittleEndian16("Start disk number");
    header.Nodes.Add(startDiskNo);

    var countFiles = startDiskNo.RelativeToLittleEndian16("Entries in this disk");
    header.Nodes.Add(countFiles);

    var countTotal = countFiles.RelativeToLittleEndian16("Entries in central dir");
    header.Nodes.Add(countTotal);

    var sizeDir = countTotal.RelativeToLittleEndian32("Size of central dir");
    header.Nodes.Add(sizeDir);

    // relative offset
    var offsetDir = sizeDir.RelativeToLittleEndian32("Offset to central dir");
    header.Nodes.Add(offsetDir);

    var commentLen = offsetDir.RelativeToLittleEndian16("Comment length");
    header.Nodes.Add(commentLen);
    var commentLenValue = commentLen.GetValue(BaseStream);

    var commentData = commentLen.RelativeTo("Comment", commentLenValue);
    if (commentData.length > 0)
        header.Nodes.Add(commentData);

    header.length = (uint)((commentData.offset + commentData.length) - identifier.offset);
    return header;
}

private Chunk ParseZipCds(uint offset)
{
    var header = new Chunk();
    header.offset = offset;

    header.Text = "Central Directory Structure";

    var identifier = new LittleEndian32BitChunk("CDS identifier");
    identifier.offset = header.offset;
    header.Nodes.Add(identifier);

    var version = identifier.RelativeTo("Version", 1);
    header.Nodes.Add(version);

    var hostOs = version.RelativeToByte("Host OS");
    var hostOsValue = hostOs.GetValue(BaseStream);
    hostOs.Text += " = " + DecodeHostOs(hostOsValue);
    header.Nodes.Add(hostOs);

    var minVersion = hostOs.RelativeTo("Minimum version", 1);
    header.Nodes.Add(minVersion);

    var targetOs = minVersion.RelativeToByte("Target OS");
    var targetOsValue = targetOs.GetValue(BaseStream);
    targetOs.Text += " = " + DecodeHostOs(targetOsValue);
    header.Nodes.Add(targetOs);

    var gpFlag = targetOs.RelativeToLittleEndian16("General purpose flag");
    header.Nodes.Add(gpFlag);

    var method = gpFlag.RelativeToLittleEndian16("Compression method");
    var methodValue = method.GetValue(BaseStream);
    method.Text += " = " + DecodeCompressionMethod(methodValue);
    header.Nodes.Add(method);

    var timestamp = method.RelativeToLittleEndianDateStamp("DateTime");
    header.Nodes.Add(timestamp);

    var crc = timestamp.RelativeToLittleEndian32("Crc");
    header.Nodes.Add(crc);

    var csize = crc.RelativeToLittleEndian32("Compressed size");
    header.Nodes.Add(csize);

    var size = csize.RelativeToLittleEndian32("Uncompressed size");
    header.Nodes.Add(size);

    var nameLen = size.RelativeToLittleEndian16("Length of filename");
    var nameLenValue = nameLen.GetValue(BaseStream);
    header.Nodes.Add(nameLen);

    var extraLen = nameLen.RelativeToLittleEndian16("Length of extra field");
    var extraLenValue = extraLen.GetValue(BaseStream);
    header.Nodes.Add(extraLen);

    var commentLen = extraLen.RelativeToLittleEndian16("Length of comment");
    var commentLenValue = commentLen.GetValue(BaseStream);
    header.Nodes.Add(commentLen);

    var diskNumber = commentLen.RelativeToLittleEndian16("Disk number");
    header.Nodes.Add(diskNumber);

    var fileAttrib = diskNumber.RelativeToLittleEndian16("Internal file attributes");
    header.Nodes.Add(fileAttrib);

    var extFileAttr = fileAttrib.RelativeToLittleEndian32("External file attributes");
    header.Nodes.Add(extFileAttr);

    var relOffset = extFileAttr.RelativeToLittleEndian32("Relative offset");
    header.Nodes.Add(relOffset);

    var nameData = relOffset.RelativeToZeroTerminatedString("Filename", nameLenValue);
    header.Nodes.Add(nameData);

    var extraData = nameData.RelativeTo("Extra field", extraLenValue);
    if (extraData.length > 0)
        header.Nodes.Add(extraData);

    var commentData = extraData.RelativeTo("Comment", commentLenValue);
    if (commentData.length > 0)
        header.Nodes.Add(commentData);

    header.length = (uint)((commentData.offset + commentData.length) - identifier.offset);

    return header;
}

override public List<Chunk> GetFileStructure()
{
    List<Chunk> res = new List<Chunk>();

    uint offset = 0;
    int count = 0;

    do {
        count++;
        BaseStream.Position = offset;

        // each ZIP chunk starts with the PK header
        if (ReadByte() != 'P' || ReadByte() != 'K') {
            Log("PARSE ERROR AT offset " + (BaseStream.Position - 2).ToString("x4"));
            break;
        }

        var b1 = ReadByte();
        var b2 = ReadByte();

        if (b1 == 1 && b2 == 2) {

            var header = ParseZipCds(offset);
            offset += header.length;
            res.Add(header);

        } else if (b1 == 3 && b2 == 4) {

            var header = ParseZipHeader(offset);
            offset += header.length;
            res.Add(header);

        } else if (b1 == 5 && b2 == 6) {

            var header = ParseEndOfCds(offset);
            offset += header.length;
            res.Add(header);

            Log("Stopped parser after End of CDS");
            break;
        }

    } while (count < 300);

    return res;
}
*/
