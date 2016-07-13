package archive

// Extensions: .xar .pkg
// https://en.wikipedia.org/wiki/Xar_%28archiver%29
// https://github.com/mackyle/xar/wiki/xarformat

// STATUS: 5%

import (
	//"bytes"
	//"compress/zlib"
	"fmt"
	//"io"
	"os"

	"github.com/martinlindhe/formats/parse"
)

// XAR parses the xar format
func XAR(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isXAR(c.Header) {
		return nil, nil
	}
	return parseXAR(c.File, c.ParsedLayout)
}

func isXAR(b []byte) bool {

	if b[0] != 'x' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}
	return true
}

var (
	checksumAlgos = map[uint32]string{
		0: "none",
		1: "sha1",
		2: "md5",
		3: "xxx special mode", // XXX some special mode "If cksum_alg is 3 then the size field MUST be a multiple of 4 and at least 32"
	}
)

func parseXAR(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	hdrLen, _ := parse.ReadUint16be(file, pos+4)
	tocLen, _ := parse.ReadUint64be(file, pos+8)

	chksumAlgo, _ := parse.ReadToMap(file, parse.Uint32be, pos+24, checksumAlgos)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 28, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 2, Info: "header size", Type: parse.Uint16be},
			{Offset: pos + 6, Length: 2, Info: "version", Type: parse.Uint16be},
			{Offset: pos + 8, Length: 8, Info: "toc length compressed", Type: parse.Uint64be},
			{Offset: pos + 16, Length: 8, Info: "toc length uncompressed", Type: parse.Uint64be},
			{Offset: pos + 24, Length: 4, Info: "checksum algorithm = " + chksumAlgo, Type: parse.Uint32be},
		}}}

	if hdrLen != 28 {
		fmt.Println("warning: xar header len. expected 28, found", hdrLen)
	}

	pos += int64(hdrLen)

	pl.Layout = append(pl.Layout, parse.Layout{
		Offset: pos,
		Length: int64(tocLen),
		Info:   "toc",
		Type:   parse.Group,

		Childs: []parse.Layout{
			{Offset: pos, Length: int64(tocLen), Info: "toc data", Type: parse.Bytes},
		}})

	// XXX toc is zlib compressed, so extract it. contains xml with rest of file struct

	/*
		toc := parse.ReadBytesFrom(file, pos, int64(tocLen))

		b := bytes.NewReader(toc)
		r, err := zlib.NewReader(b)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, r)
		r.Close()
	*/

	return &pl, nil
}
