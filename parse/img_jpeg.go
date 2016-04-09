package parse

// STATUS: 0%

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

	fileHeader := Layout{
		Offset: offset,
		Info:   "header",
		Type:   Group,
	}

	for {
		fmt.Printf("Reading at %04x\n", offset)
		magic, _ := readUint8(file, offset)
		marker, _ := readUint8(file, offset+1)

		if magic != 0xff {
			fmt.Printf("jpeg parse error, found %02x\n", marker)
			break
		}

		if marker == 0xd8 {
			// NOTE: this marker dont have any content
			fileHeader.Childs = append(fileHeader.Childs, []Layout{
				Layout{Offset: offset, Length: 2, Info: jpegChunkTypes[0xd8], Type: Uint16le},
			}...)
			offset += 2
			continue
		}
		if marker == 0xd9 {
			fileHeader.Childs = append(fileHeader.Childs, []Layout{
				Layout{Offset: offset, Length: 2, Info: jpegChunkTypes[0xd9], Type: Uint16be},
			}...)
			fmt.Println("Ending parser since EOI marker was detected")
			break
		}

		chunkLen, _ := readUint16be(file, offset+2)
		fmt.Println("adding", chunkLen)

		childs := []Layout{
			Layout{Offset: offset, Length: 2, Info: "type " + jpegChunkTypes[marker], Type: Uint16be},
			Layout{Offset: offset + 2, Length: 2, Info: "length", Type: Uint16be},
			Layout{Offset: offset + 4, Length: int64(chunkLen) - 2, Info: "data", Type: Bytes}, // XXX
		}

		offset += 2 + int64(chunkLen)

		fileHeader.Childs = append(fileHeader.Childs, childs...)
	}

	res.Layout = append(res.Layout, fileHeader)
	return &res, nil
}

/*


   var length = type.RelativeToBigEndian16("Length");
   marker.Nodes.Add(length);

   uint lenghtValue = (uint)ReadInt16BE(length.offset);

   marker.length = 2 + lenghtValue;
   Log("len = " + marker.length);


   var data = length.RelativeTo("Data", lenghtValue - 2);
   marker.Nodes.Add(data);

   switch (marker1) {
   case 0xC0:
       marker.Text = "SOF0 - Baseline DCT";
       break;
   case 0xC1:
       marker.Text = "SOF1 - Extended sequential DCT";
       break;
   case 0xC2:
       marker.Text = "SOF2 - Progressive DCT";
       break;
   case 0xC3:
       marker.Text = "SOF3 - Lossless (sequential)";
       break;
   case 0xC4:
       marker.Text = "DHT - Huffman table";
       break;

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

   case 0xDB:
       marker.Text = "DQT - Quantization table";
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

   case 0xE1:
       marker.Text = "APP1 - Application Use";
       // XXXX
       break;

   case 0xFE:
       marker.Text = "COM - Comment";
       break;

   default:
       throw new Exception("TODO " + marker1.ToString("x2"));
   }

   BaseStream.Position += lenghtValue - 2;
   res.Add(marker);


*/
