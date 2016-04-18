package parse

// STATUS: 0% - borked!
// http://search.cpan.org/~wiml/Mac-Finder-DSStore/DSStoreFormat.pod
// https://en.wikipedia.org/wiki/.DS_Store

import (
	"encoding/binary"
	"os"
)

func DSSTORE(file *os.File) (*ParsedLayout, error) {

	if !isDSSTORE(file) {
		return nil, nil
	}
	return parseDSSTORE(file)
}

func isDSSTORE(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX just guessing
	//	if b[0] != 0xff || b[1] != 0xfe || b[2] != 0x23 || b[3] != 0 {
	//		return false
	//	}
	return false
}

func parseDSSTORE(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	header := Layout{
		Offset: pos,
		Info:   "header",
		Type:   Group}

	for {

		p, _ := readUint32be(file, pos)
		//count, _ := readUint32be(file, offset+4)

		header.Childs = append(header.Childs, []Layout{
			// XXX: Each node starts with two integers, P and count.
			{Offset: pos, Length: 4, Info: "p (rightmost child)", Type: Uint32be},
			{Offset: pos + 4, Length: 4, Info: "count", Type: Uint32be},
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
			header.Childs = append(header.Childs, []Layout{
				{Offset: pos, Length: 4, Info: "block num of leftmost child", Type: Uint32be},
				{Offset: pos + 4, Length: 4, Info: "record", Type: Uint32be},
			}...)
			pos += 8           // XXX
			header.Length += 8 // XXX
			//}
		}

		// XXX loop some times
		break

	}

	layout := []Layout{header}

	res := ParsedLayout{
		FileKind: Executable,
		Layout:   layout}

	return &res, nil
}
