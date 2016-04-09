package parse

// STATUS: 50%

import (
	"encoding/binary"
	"fmt"
	"os"
)

var (
	jpegChunkTypes = map[byte]string{ // "marker"
		0xC0: "baseline DCT (SOF0)",
		0xC1: "extended sequential DCT (SOF1)",
		0xC2: "progressive DCT (SOF2)",
		0xC3: "lossless (SOF3)",
		0xC4: "huffman table (DHT)",
		0xD8: "start of image (SOI)",
		0xD9: "end of image (EOI)",
		0xDA: "start of scan (SOS)",
		0xDB: "quantization table (DQT)",
		0xE0: "APP0",
		0xE1: "APP1",
		0xFE: "comment (COM)",
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

		magic, _ := readUint8(file, offset)
		marker, _ := readUint8(file, offset+1)

		fmt.Printf("Reading jpeg chunk at %04x. marker %02x\n", offset, marker)

		chunk.Info = jpegChunkTypes[marker]

		if magic != 0xff {
			fmt.Printf("jpeg parse error. expected ff, found %02x\n", magic)
			break
		}

		if marker == 0xd8 {
			// NOTE: this marker dont have any content
			chunk.Length = 2
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: "type", Type: Uint16le},
			}
			res.Layout = append(res.Layout, chunk)
			offset += chunk.Length
			continue
		}
		if marker == 0xd9 {
			chunk.Length = 2
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: "type", Type: Uint16be},
			}
			res.Layout = append(res.Layout, chunk)
			fmt.Println("Ending parser since EOI marker was detected")
			break
		}

		if marker == 0xda { // start of scan

			chunk.Length = 5 // XXX
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: "type", Type: Uint16be},
				Layout{Offset: offset + 2, Length: 2, Info: "length", Type: Uint16be},
				Layout{Offset: offset + 4, Length: 1, Info: "color components", Type: Uint8},
			}

			components, _ := readUint8(file, offset+4)
			fmt.Println("adding", components, "components")

			offset += chunk.Length
			for i := 0; i < int(components); i++ {
				chunk.Childs = append(chunk.Childs, []Layout{
					//        For each component
					//        An ID
					//        An AC table # (Low Nibble)
					//        An DC table # (High Nibble)
					Layout{Offset: offset, Length: 2, Info: "color", Type: Uint16be},
				}...)
				chunk.Length += 2
				offset += 2
			}

			chunk.Childs = append(chunk.Childs, []Layout{
				Layout{Offset: offset, Length: 3, Info: "unknown", Type: Bytes},
			}...)
			chunk.Length += 3

			res.Layout = append(res.Layout, chunk)
			offset += chunk.Length
			continue

			// TODO: now follows compressed image data, how to get length of data?
		}

		if marker == 0xe0 {
			// APP0
			chunk.Length = 18 // XXX
			chunk.Childs = []Layout{
				Layout{Offset: offset, Length: 2, Info: "type", Type: Uint16be},
				Layout{Offset: offset + 2, Length: 2, Info: "length", Type: Uint16be},
				Layout{Offset: offset + 4, Length: 5, Info: "identifier", Type: ASCII},
				Layout{Offset: offset + 9, Length: 2, Info: "revision", Type: MajorMinor16},
				Layout{Offset: offset + 11, Length: 1, Info: "units used", Type: Uint8},
				Layout{Offset: offset + 12, Length: 2, Info: "width", Type: Uint16be},
				Layout{Offset: offset + 14, Length: 2, Info: "height", Type: Uint16be},
				Layout{Offset: offset + 16, Length: 1, Info: "horizontal pixels", Type: Uint8},
				Layout{Offset: offset + 17, Length: 1, Info: "vertical pixels", Type: Uint8},
			}
			res.Layout = append(res.Layout, chunk)
			offset += chunk.Length
			continue
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
