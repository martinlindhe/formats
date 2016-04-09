package parse

// handles PNG and MNG images
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func PNG(file *os.File) (*ParsedLayout, error) {

	if !isPNG(file) {
		return nil, nil
	}
	return parsePNG(file)
}

func getPNGHeader(file *os.File) ([8]byte, error) {
	file.Seek(0, os.SEEK_SET)

	var b [8]uint8
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
}
func isPNG(file *os.File) bool {

	b, err := getPNGHeader(file)
	if err != nil {
		return false
	}

	if (b[0] == 0x89 && b[1] == 'P') || // png
		(b[0] == 0x8a && b[1] == 'M') { // mng
		if b[2] == 'N' && b[3] == 'G' && b[4] == 0xd &&
			b[5] == 0xa && b[6] == 0x1a && b[7] == 0xa {
			return true
		}
	}

	return false
}

func parsePNG(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)
	res := ParsedLayout{}

	b, err := getPNGHeader(file)
	if err != nil {
		return nil, err
	}
	fileType := "PNG"
	if b[1] == 'M' {
		fileType = "MNG"
	}

	fileHeader := Layout{
		Offset: offset,
		Length: 8,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 8, Info: "magic = " + fileType, Type: Bytes},
		},
	}

	res.Layout = append(res.Layout, fileHeader)

	offset = 8

	chunks := []Layout{}
	for {
		l := Layout{
			Offset: offset,
			Length: 8,
			Info:   "XXX chunk",
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: offset, Length: 4, Info: "length", Type: Uint32be},
				Layout{Offset: offset + 4, Length: 4, Info: "type", Type: ASCII}, // XXX
			},
		}

		chunks = append(chunks, l)
		break // XXX

		/*

			chunk.Text = "Chunk " + typeStr;
			chunk.length = lengthVal + 4 + 4 + 4;  // "length" (4 byte) + "type" (4 byte) + data + crc (4 byte)

			var data = type.RelativeTo("Data", lengthVal);

			if (lengthVal > 0) {
			    if (typeStr == "IHDR") {
			        var width = type.RelativeToBigEndian32("Width");
			        data.Nodes.Add(width);

			        var height = width.RelativeToBigEndian32("Height");
			        data.Nodes.Add(height);

			        var bd = height.RelativeToByte("Bit depth");
			        data.Nodes.Add(bd);

			        var color = bd.RelativeToByte("Color type");
			        data.Nodes.Add(color);

			        var compression = color.RelativeToByte("Compression method");
			        data.Nodes.Add(compression);

			        var filter = compression.RelativeToByte("Filter method");
			        data.Nodes.Add(filter);

			        var interlace = filter.RelativeToByte("Interlace method");
			        data.Nodes.Add(interlace);
			    }

			    chunk.Nodes.Add(data);
			}

			var crc = data.RelativeToBigEndian32("Crc");
			chunk.Nodes.Add(crc);

			offset += chunk.length;

			res.Add(chunk);

			if (typeStr == "IEND") {
			    Log("Stopped parser after IEND chunk");
			    break;
			}

			//} while (offset < BaseStream.Length);
		*/
	}

	res.Layout = append(res.Layout, chunks...)

	return &res, nil
}

/*
private string ReadString(long offset, int length)
{
    BaseStream.Position = offset;

    string res = "";

    for (int i = 0; i < length; i++) {
        res += ReadChar();
    }

    return res;
}
*/
