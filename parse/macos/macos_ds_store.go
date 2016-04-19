package macos

// http://search.cpan.org/~wiml/Mac-Finder-DSStore/DSStoreFormat.pod
// https://en.wikipedia.org/wiki/.DS_Store

// STATUS: 0%
// XXX borked

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func DSSTORE(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isDSSTORE(&hdr) {
		return nil, nil
	}
	return parseDSSTORE(file, pl)
}

func isDSSTORE(hdr *[0xffff]byte) bool {

	return false
	/*
		b := *hdr
		// XXX just guessing
		if b[0] != 0xff || b[1] != 0xfe || b[2] != 0x23 || b[3] != 0 {
			return false
		}
		return true
	*/
}

func parseDSSTORE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	header := parse.Layout{
		Offset: pos,
		Info:   "header",
		Type:   parse.Group}

	for {

		p, _ := parse.ReadUint32be(file, pos)
		//count, _ := readUint32be(file, offset+4)

		header.Childs = append(header.Childs, []parse.Layout{
			// XXX: Each node starts with two integers, P and count.
			{Offset: pos, Length: 4, Info: "p (rightmost child)", Type: parse.Uint32be},
			{Offset: pos + 4, Length: 4, Info: "count", Type: parse.Uint32be},
		}...)

		pos += 8           // XXX
		header.Length += 8 // XXX

		if p == 0 {
			// XXX: If P is 0, then this is a leaf node and count is immediately
			// followed by that many records.
		}

		if p != 0 {
			// XXX:  If P is nonzero, then this is an internal node, and count
			// is followed by the block number of the leftmost child, then a
			// record, then another block number, etc., for a total of count
			// child pointers and count records. P is itself the rightmost child
			// pointer, that is, it is logically at the end of the node.
			//for j := uint32(0); j < count; j++ {
			header.Childs = append(header.Childs, []parse.Layout{
				{Offset: pos, Length: 4, Info: "block num of leftmost child", Type: parse.Uint32be},
				{Offset: pos + 4, Length: 4, Info: "record", Type: parse.Uint32be},
			}...)
			pos += 8           // XXX
			header.Length += 8 // XXX
			//}
		}

		// XXX loop some times
		break

	}

	layout := []parse.Layout{header}

	pl.FileKind = parse.Executable
	pl.Layout = layout

	return &pl, nil
}
