package parse

// STATUS: 80%. TODO need sample with exif data
// XXX samples/jpg/jpeg_002.jpg parse OK!
// XXX samples/jpg/jpeg_001.jpg loops forever

import (
	"encoding/binary"
	"fmt"
	"io"
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

	pos := int64(0)
	res := ParsedLayout{
		FileKind: Image,
	}

	for {
		magic, _ := ReadUint8(file, pos)
		marker, _ := ReadUint8(file, pos+1)

		// fmt.Printf("Reading jpeg chunk at %04x. marker %02x\n", offset, marker)
		if magic != 0xff {
			fmt.Printf("jpeg parse error at %04x. expected ff, found %02x\n", pos, magic)
			break
		}

		if marker == 0xd8 { // start of image
			// NOTE: this marker dont have any content
			res.Layout = append(res.Layout, Layout{
				Offset: pos,
				Type:   Group,
				Length: 2,
				Info:   jpegChunkTypes[marker],
				Childs: []Layout{
					{Offset: pos, Length: 2, Info: "type", Type: Uint16le},
				},
			})
			pos += 2
			continue
		}
		if marker == 0xd9 {
			res.Layout = append(res.Layout, Layout{
				Offset: pos,
				Type:   Group,
				Length: 2,
				Info:   jpegChunkTypes[marker],
				Childs: []Layout{
					{Offset: pos, Length: 2, Info: "type", Type: Uint16le},
				},
			})
			// fmt.Println("Ending parser since EOI marker was detected")
			break
		}

		if marker == 0xda { // start of scan

			components, _ := ReadUint8(file, pos+4)
			chunk := Layout{
				Offset: pos,
				Type:   Group,
				Length: 5,
				Info:   jpegChunkTypes[marker],
				Childs: []Layout{
					{Offset: pos, Length: 2, Info: "type", Type: Uint16be},
					{Offset: pos + 2, Length: 2, Info: "length", Type: Uint16be},
					{Offset: pos + 4, Length: 1, Info: "color components", Type: Uint8},
				},
			}
			pos += chunk.Length

			for i := 0; i < int(components); i++ {
				chunk.Childs = append(chunk.Childs, []Layout{
					{Offset: pos, Length: 1, Info: "color id", Type: Uint8},

					// XXX decode values:
					// An AC table # (Low Nibble)
					// An DC table # (High Nibble)
					{Offset: pos + 1, Length: 1, Info: "ac,dc tables", Type: Uint8}, // XXX hi/lo nibbles type
				}...)
				chunk.Length += 2
				pos += 2
			}

			chunk.Childs = append(chunk.Childs, []Layout{
				{Offset: pos, Length: 3, Info: "unknown", Type: Bytes},
			}...)
			chunk.Length += 3
			pos += 3

			file.Seek(pos, os.SEEK_SET)

			// fmt.Printf("starting at %04x\n", offset)
			dataStart := pos
			rewind := false
			res.Layout = append(res.Layout, chunk)

			for {
				var b uint16
				err := binary.Read(file, binary.BigEndian, &b)
				if err == io.EOF {
					fmt.Println(err)
					break
				}

				// fmt.Printf("from %04x: %04x\n", offset, b)
				pos += 2

				marker := b & 0xff
				if b&0xff00 == 0xff00 && marker != 0 && (marker < 0xd0 || marker > 0xd8) {
					// eoi
					dataLen := pos - dataStart - 2

					res.Layout = append(res.Layout, Layout{
						Offset: dataStart,
						Length: dataLen,
						Type:   Group,
						Info:   "image data",
						Childs: []Layout{
							{Offset: dataStart, Length: dataLen, Info: "image data", Type: Bytes},
						}})

					rewind = true
					pos -= 2
					// fmt.Printf("rewinded offset to %04x\n", offset)
					file.Seek(pos, os.SEEK_SET)
					break
				}
			}

			if !rewind {
				pos += chunk.Length
			}
			continue
		}

		if marker == 0xe0 {
			// APP0
			res.Layout = append(res.Layout, Layout{
				Offset: pos,
				Type:   Group,
				Length: 18,
				Info:   jpegChunkTypes[marker],
				Childs: []Layout{
					{Offset: pos, Length: 2, Info: "type", Type: Uint16be},
					{Offset: pos + 2, Length: 2, Info: "length", Type: Uint16be},
					{Offset: pos + 4, Length: 5, Info: "identifier", Type: ASCII},
					{Offset: pos + 9, Length: 2, Info: "revision", Type: MajorMinor16be},
					{Offset: pos + 11, Length: 1, Info: "units used", Type: Uint8},
					{Offset: pos + 12, Length: 2, Info: "width", Type: Uint16be},
					{Offset: pos + 14, Length: 2, Info: "height", Type: Uint16be},
					{Offset: pos + 16, Length: 1, Info: "horizontal pixels", Type: Uint8},
					{Offset: pos + 17, Length: 1, Info: "vertical pixels", Type: Uint8},
				}})
			pos += 18
			continue
		}

		chunkLen, _ := ReadUint16be(file, pos+2)

		res.Layout = append(res.Layout, Layout{
			Offset: pos,
			Type:   Group,
			Length: 2 + int64(chunkLen),
			Info:   jpegChunkTypes[marker],
			Childs: []Layout{
				{Offset: pos, Length: 2, Info: "type", Type: Uint16be},
				{Offset: pos + 2, Length: 2, Info: "length", Type: Uint16be},
				{Offset: pos + 4, Length: int64(chunkLen) - 2, Info: "data", Type: Bytes}, // XXX
			}})
		pos += 2 + int64(chunkLen)
	}

	return &res, nil
}
