package windows

// Direct3D shader bytecode

// http://timjones.tw/blog/archive/2015/09/02/parsing-direct3d-shader-bytecode

// STATUS: 30%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func DXBC(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isDXBC(&c.Header) {
		return nil, nil
	}
	return parseDXBC(c.File, c.ParsedLayout)
}

func isDXBC(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'D' || b[1] != 'X' || b[2] != 'B' || b[3] != 'C' {
		return false
	}
	return true
}

func parseDXBC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource

	chunkCount, _ := parse.ReadUint32le(file, pos+28)
	header := dxbcHeader(pos)
	pl.Layout = []parse.Layout{header}
	pos += header.Length

	if chunkCount > 0 {

		chunkOffsets := parse.Layout{
			Offset: pos,
			Length: int64(chunkCount) * 4,
			Info:   "chunk offsets",
			Type:   parse.Group}
		pl.Layout = append(pl.Layout, chunkOffsets)

		for i := 0; i < int(chunkCount); i++ {
			id := fmt.Sprintf("%d", i)
			chunkOffsets.Childs = append(chunkOffsets.Childs, parse.Layout{
				Offset: pos,
				Length: 4,
				Info:   "chunk " + id + " offset",
				Type:   parse.Uint32le})

			chunkPos, _ := parse.ReadUint32le(file, pos)
			chunkLen, _ := parse.ReadUint32le(file, int64(chunkPos)+4)

			pos += 4
			pl.Layout = append(pl.Layout, parse.Layout{
				Offset: int64(chunkPos),
				Length: 8 + int64(chunkLen),
				Info:   "chunk :)",
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: int64(chunkPos), Length: 4, Info: "type", Type: parse.ASCII},
					{Offset: int64(chunkPos) + 4, Length: 4, Info: "chunk size", Type: parse.Uint32le},

					// XXX decode fields depending on chunk id
					{Offset: int64(chunkPos) + 8, Length: int64(chunkLen), Info: "data", Type: parse.Bytes},
				}})
		}
		pos += int64(chunkCount) * 4
	}

	return &pl, nil
}

func dxbcHeader(pos int64) parse.Layout {
	return parse.Layout{
		Offset: pos,
		Length: 32, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 4, Info: "checksum 1", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "checksum 2", Type: parse.Uint32le},
			{Offset: pos + 12, Length: 4, Info: "checksum 3", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 4, Info: "checksum 4", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "unknown", Type: parse.Uint32le}, // always 1 ?
			{Offset: pos + 24, Length: 4, Info: "total size", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "chunk count", Type: parse.Uint32le},
		}}
}
