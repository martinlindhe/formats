package parse

// STATUS: 50%

import (
	"encoding/binary"
	"fmt"
	"os"
)

var (
	jpegChunkTypes = map[byte]string{ // "marker"
		0xC0: "SOF0 - Baseline DCT",
		0xC1: "SOF1 - Extended sequential DCT",
		0xC2: "SOF2 - Progressive DCT",
		0xC3: "SOF3 - Lossless (sequential)",
		0xC4: "DHT - Huffman table",
		0xD8: "SOI - start of image",
		0xD9: "EOI - End of Image",
		0xDA: "SOS - Start of scan",
		0xDB: "DQT - Quantization table",
		0xE0: "APP0 - Application Use",
		0xE1: "APP1 - Application Use",
		0xFE: "COM - Comment",
	}
)

func JPEG(file *os.File) (*ParsedLayout, error) {

	if !isJPEG(file) {
		return nil, nil
	}
	return parseJPEG(file)
}

func isJPEG(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [12]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 0xff || b[1] != 0xd8 {
		return false
	}

	if b[6] != 'J' || b[7] != 'F' || b[8] != 'I' || b[9] != 'F' || b[10] != 0 {
		return false
	}
	return true
}

func parseJPEG(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	offset := int64(0)

	for {

		chunk := Layout{
			Offset: offset,
			Type:   Group,
		}

		fmt.Printf("Reading jpeg chunk at %04x\n", offset)
		magic, _ := readUint8(file, offset)
		marker, _ := readUint8(file, offset+1)

		chunk.Info = jpegChunkTypes[marker]

		if magic != 0xff {
			fmt.Printf("jpeg parse error, found %02x\n", marker)
			break
		}

		if marker == 0xd8 {
			// NOTE: this marker dont have any content
			chunk.Length = 2
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: jpegChunkTypes[0xd8], Type: Uint16le},
			}
			res.Layout = append(res.Layout, chunk)
			offset += chunk.Length
			continue
		}
		if marker == 0xd9 {
			chunk.Length = 2
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: jpegChunkTypes[0xd9], Type: Uint16be},
			}
			res.Layout = append(res.Layout, chunk)
			fmt.Println("Ending parser since EOI marker was detected")
			break
		}

		chunkLen, _ := readUint16be(file, offset+2)

		chunk.Length = 2 + int64(chunkLen)
		chunk.Childs = []Layout{
			Layout{Offset: offset, Length: 2, Info: "type", Type: Uint16be},
			Layout{Offset: offset + 2, Length: 2, Info: "length", Type: Uint16be},
			Layout{Offset: offset + 4, Length: int64(chunkLen) - 2, Info: "data", Type: Bytes}, // XXX
		}

		offset += 2 + int64(chunkLen)

		res.Layout = append(res.Layout, chunk)
	}

	return &res, nil
}

/*

case 0xDA:
   marker.Text = "SOS - Start of scan";

   int component_count = ReadByte(data.offset);

   var Components = length.RelativeToByte("Color components");
   data.Nodes.Add(Components);

   for (int i = 0; i < component_count; i++) {
       //        For each component
       //        An ID
       //        An AC table # (Low Nibble)
       //        An DC table # (High Nibble)

       var kex = new BigEndian16BitChunk("Color");
       kex.offset = data.offset + 1 + (i * kex.length);
       data.Nodes.Add(kex);
   }

   var Unknown = new Chunk();
   Unknown.offset = data.offset + 1 + (2 * component_count);
   Unknown.length = 3;
   Unknown.Text = "Unknown";
   data.Nodes.Add(Unknown);

   // TODO: now follows compressed image data, how to get length of data?

   break;

case 0xE0:
   marker.Text = "APP0 - Application Use";

   var Identifier = length.RelativeToZeroTerminatedString("Id String", 5);
   marker.Nodes.Add(Identifier);

   var Version = Identifier.RelativeToVersionMajorMinor16("Revision");
   marker.Nodes.Add(Version);

   // Units used for Resolution
   var Units = Version.RelativeToByte("Units used");
   marker.Nodes.Add(Units);

   var Width = Units.RelativeToBigEndian16("Width");
   marker.Nodes.Add(Width);

   var Height = Width.RelativeToBigEndian16("Height");
   marker.Nodes.Add(Height);

   var XThumbnail = Height.RelativeToByte("Horizontal Pixels");
   marker.Nodes.Add(XThumbnail);

   var YThumbnail = XThumbnail.RelativeToByte("Vertical Pixels");
   marker.Nodes.Add(YThumbnail);
   break;

*/
