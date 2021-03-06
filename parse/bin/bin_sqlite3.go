package bin

// https://www.sqlite.org/fileformat2.html

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	sqlite3TextEncodings = map[uint32]string{
		1: "utf-8",
		2: "utf-16le",
		3: "utf-16be",
	}
)

// SQLITE3 parses the sqlite3 format
func SQLITE3(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isSQLITE3(c.Header) {
		return nil, nil
	}
	return parseSQLITE3(c.File, c.ParsedLayout)
}

func isSQLITE3(b []byte) bool {

	s := string(b[0:16])
	return s == "SQLite format 3"
}

func parseSQLITE3(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	textEncodingName, _ := parse.ReadToMap(file, parse.Uint32be, pos+56, sqlite3TextEncodings)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 100, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 16, Info: "magic", Type: parse.ASCIIZ},
			{Offset: pos + 16, Length: 2, Info: "page size", Type: parse.Uint16be},
			{Offset: pos + 18, Length: 1, Info: "write version", Type: parse.Int8},
			{Offset: pos + 19, Length: 1, Info: "read version", Type: parse.Int8},
			{Offset: pos + 20, Length: 1, Info: "reserved", Type: parse.Int8},
			{Offset: pos + 21, Length: 1, Info: "max embedded payload fraction", Type: parse.Int8}, // Must be 64
			{Offset: pos + 22, Length: 1, Info: "min embedded payload fraction", Type: parse.Int8}, // Must be 32
			{Offset: pos + 23, Length: 1, Info: "leaf payload fraction", Type: parse.Int8},         // Must be 32
			{Offset: pos + 24, Length: 4, Info: "file change counter", Type: parse.Uint32be},
			{Offset: pos + 28, Length: 4, Info: "size of database file in pages", Type: parse.Uint32be},
			{Offset: pos + 32, Length: 4, Info: "page number of the first freelist trunk page", Type: parse.Uint32be},
			{Offset: pos + 36, Length: 4, Info: "total number of freelist pages", Type: parse.Uint32be},
			{Offset: pos + 40, Length: 4, Info: "schema cookie", Type: parse.Uint32be},
			{Offset: pos + 44, Length: 4, Info: "schema format number", Type: parse.Uint32be}, // allowed 1-4
			{Offset: pos + 48, Length: 4, Info: "default page cache size", Type: parse.Uint32be},
			{Offset: pos + 52, Length: 4, Info: "page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.", Type: parse.Uint32be},
			{Offset: pos + 56, Length: 4, Info: "text encoding = " + textEncodingName, Type: parse.Uint32be},
			{Offset: pos + 60, Length: 4, Info: "user version", Type: parse.Uint32be}, // XXX decode, as used by https://www.sqlite.org/pragma.html#pragma_schema_version
			{Offset: pos + 64, Length: 4, Info: "True (non-zero) for incremental-vacuum mode. False (zero) otherwise. ", Type: parse.Uint32be},
			{Offset: pos + 68, Length: 4, Info: "application id", Type: parse.Uint32be}, // set by PRAGMA application_id
			{Offset: pos + 72, Length: 20, Info: "reserved", Type: parse.Uint32be},
			{Offset: pos + 92, Length: 4, Info: "version valid for", Type: parse.Uint32be},     // XXX https://www.sqlite.org/fileformat2.html#validfor
			{Offset: pos + 96, Length: 4, Info: "SQLITE_VERSION_NUMBER", Type: parse.Uint32be}, // XXX https://www.sqlite.org/c3ref/c_source_id.html
		}}}

	return &pl, nil
}
