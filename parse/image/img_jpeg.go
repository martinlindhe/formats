package image

// TODO need sample with exif data
// XXX samples/jpg/jpeg_002.jpg parse OK!
// XXX samples/jpg/jpeg_001.jpg loops forever

// STATUS: 80%

import (
	"encoding/binary"
	"fmt"
	"github.com/martinlindhe/formats/parse"
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

func JPEG(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isJPEG(&hdr) {
		return nil, nil
	}
	return parseJPEG(file, pl)
}

func isJPEG(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0xff || b[1] != 0xd8 {
		return false
	}
	if b[6] != 'J' || b[7] != 'F' || b[8] != 'I' || b[9] != 'F' || b[10] != 0 {
		return false
	}
	return true
}

func parseJPEG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image

	for {
		magic, _ := parse.ReadUint8(file, pos)
		marker, _ := parse.ReadUint8(file, pos+1)

		// fmt.Printf("Reading jpeg chunk at %04x. marker %02x\n", offset, marker)
		if magic != 0xff {
			fmt.Printf("jpeg parse error at %04x. expected ff, found %02x\n", pos, magic)
			break
		}

		if marker == 0xd8 { // start of image
			// NOTE: this marker dont have any content
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Length: 2,
				Info:   jpegChunkTypes[marker],
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16le},
				}})
			pos += 2
			continue
		}
		if marker == 0xd9 {
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Type:   parse.Group,
				Length: 2,
				Info:   jpegChunkTypes[marker],
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16le},
				}})
			// fmt.Println("Ending parser since EOI marker was detected")
			break
		}

		if marker == 0xda { // start of scan

			components, _ := parse.ReadUint8(file, pos+4)
			chunk := parse.Layout{
				Offset: pos,
				Length: 5,
				Info:   jpegChunkTypes[marker],
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
					{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
					{Offset: pos + 4, Length: 1, Info: "color components", Type: parse.Uint8},
				}}
			pos += chunk.Length

			for i := 0; i < int(components); i++ {
				chunk.Childs = append(chunk.Childs, []parse.Layout{
					{Offset: pos, Length: 1, Info: "color id", Type: parse.Uint8},

					// XXX decode values:
					// An AC table # (Low Nibble)
					// An DC table # (High Nibble)
					{Offset: pos + 1, Length: 1, Info: "ac,dc tables", Type: parse.Uint8}, // XXX hi/lo nibbles type
				}...)
				chunk.Length += 2
				pos += 2
			}

			chunk.Childs = append(chunk.Childs, []parse.Layout{
				{Offset: pos, Length: 3, Info: "unknown", Type: parse.Bytes},
			}...)
			chunk.Length += 3
			pos += 3

			file.Seek(pos, os.SEEK_SET)

			// fmt.Printf("starting at %04x\n", offset)
			dataStart := pos
			rewind := false
			pl.Layout = append(pl.Layout, chunk)

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

					pl.Layout = append(pl.Layout, parse.Layout{
						Offset: dataStart,
						Length: dataLen,
						Info:   "image data",
						Type:   parse.Group,
						Childs: []parse.Layout{
							{Offset: dataStart, Length: dataLen, Info: "image data", Type: parse.Bytes},
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
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: pos,
				Length: 18,
				Info:   jpegChunkTypes[marker],
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
					{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
					{Offset: pos + 4, Length: 5, Info: "identifier", Type: parse.ASCII},
					{Offset: pos + 9, Length: 2, Info: "revision", Type: parse.MajorMinor16be},
					{Offset: pos + 11, Length: 1, Info: "units used", Type: parse.Uint8},
					{Offset: pos + 12, Length: 2, Info: "width", Type: parse.Uint16be},
					{Offset: pos + 14, Length: 2, Info: "height", Type: parse.Uint16be},
					{Offset: pos + 16, Length: 1, Info: "horizontal pixels", Type: parse.Uint8},
					{Offset: pos + 17, Length: 1, Info: "vertical pixels", Type: parse.Uint8},
				}})
			pos += 18
			continue
		}

		chunkLen, _ := parse.ReadUint16be(file, pos+2)

		pl.Layout = append(pl.Layout, parse.Layout{
			Offset: pos,
			Length: 2 + int64(chunkLen),
			Info:   jpegChunkTypes[marker],
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Info: "type", Type: parse.Uint16be},
				{Offset: pos + 2, Length: 2, Info: "length", Type: parse.Uint16be},
				{Offset: pos + 4, Length: int64(chunkLen) - 2, Info: "data", Type: parse.Bytes}, // XXX
			}})
		pos += 2 + int64(chunkLen)
	}

	return &pl, nil
}
