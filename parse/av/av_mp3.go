package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func MP3(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isMP3(file) {
		return nil, nil
	}
	return parseMP3(file, pl)
}

func isMP3(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO find mp3 stream start, ignore id3 tags

	if b[0] != 'I' || b[1] != 'D' || b[2] != '3' {
		return false
	}
	/*
	   byte id3ver = ReadByte();
	   if (id3ver != 3 && id3ver != 4)
	       return false;
	*/

	return true
}

func parseMP3(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
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
