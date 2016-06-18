package exe

// OS/2 Lineear eXecutable
// http://www.virtualbox.org/svn/kstuff-mirror/trunk/include/k/kLdrFmts/lx.h

// STATUS: 10%

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

// parse OS/2 Linear eXecutable header
func parseMzLxHeader(file *os.File, pos int64) ([]parse.Layout, error) {

	res := []parse.Layout{{
		Offset: pos,
		Length: 196,
		Info:   "external header = LX",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
			{Offset: pos + 2, Length: 1, Info: "border", Type: parse.Uint8},
			{Offset: pos + 3, Length: 1, Info: "worder", Type: parse.Uint8},
			{Offset: pos + 4, Length: 4, Info: "level", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 2, Info: "cpu", Type: parse.Uint16le},
			{Offset: pos + 10, Length: 2, Info: "os", Type: parse.Uint16le},
			{Offset: pos + 12, Length: 4, Info: "version", Type: parse.MajorMinor32le},
			{Offset: pos + 16, Length: 4, Info: "mflags", Type: parse.Uint32le},
			{Offset: pos + 20, Length: 4, Info: "mpages", Type: parse.Uint32le},
			{Offset: pos + 24, Length: 4, Info: "startobj", Type: parse.Uint32le},
			{Offset: pos + 28, Length: 4, Info: "eip", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 4, Info: "stackobj", Type: parse.Uint32le},
			{Offset: pos + 36, Length: 4, Info: "esp", Type: parse.Uint32le},
			{Offset: pos + 40, Length: 4, Info: "pagesize", Type: parse.Uint32le},
			{Offset: pos + 44, Length: 4, Info: "pageshift", Type: parse.Uint32le},
			{Offset: pos + 48, Length: 4, Info: "fixup size", Type: parse.Uint32le},
			{Offset: pos + 52, Length: 4, Info: "fixup sum", Type: parse.Uint32le},
			{Offset: pos + 56, Length: 4, Info: "ldrsize", Type: parse.Uint32le},
			{Offset: pos + 60, Length: 4, Info: "ldrsum", Type: parse.Uint32le},
			{Offset: pos + 64, Length: 4, Info: "objtab", Type: parse.Uint32le},
			{Offset: pos + 68, Length: 4, Info: "objcnt", Type: parse.Uint32le},
			{Offset: pos + 72, Length: 4, Info: "objmap", Type: parse.Uint32le},
			{Offset: pos + 76, Length: 4, Info: "itermap", Type: parse.Uint32le},
			{Offset: pos + 80, Length: 4, Info: "rsrctab", Type: parse.Uint32le},
			{Offset: pos + 84, Length: 4, Info: "rsrccnt", Type: parse.Uint32le},
			{Offset: pos + 88, Length: 4, Info: "restab", Type: parse.Uint32le},
			{Offset: pos + 92, Length: 4, Info: "enttab", Type: parse.Uint32le},
			{Offset: pos + 96, Length: 4, Info: "dirtab", Type: parse.Uint32le},
			{Offset: pos + 100, Length: 4, Info: "dircnt", Type: parse.Uint32le},
			{Offset: pos + 104, Length: 4, Info: "fpagetab", Type: parse.Uint32le},
			{Offset: pos + 108, Length: 4, Info: "frectab", Type: parse.Uint32le},
			{Offset: pos + 112, Length: 4, Info: "impmod", Type: parse.Uint32le},
			{Offset: pos + 116, Length: 4, Info: "impmodcnt", Type: parse.Uint32le},
			{Offset: pos + 120, Length: 4, Info: "impproc", Type: parse.Uint32le},
			{Offset: pos + 124, Length: 4, Info: "pagesum", Type: parse.Uint32le},
			{Offset: pos + 128, Length: 4, Info: "datapage", Type: parse.Uint32le},
			{Offset: pos + 132, Length: 4, Info: "preload", Type: parse.Uint32le},
			{Offset: pos + 136, Length: 4, Info: "nrestab", Type: parse.Uint32le},
			{Offset: pos + 140, Length: 4, Info: "cbnrestab", Type: parse.Uint32le},
			{Offset: pos + 144, Length: 4, Info: "nressum", Type: parse.Uint32le},
			{Offset: pos + 148, Length: 4, Info: "autodata", Type: parse.Uint32le},
			{Offset: pos + 152, Length: 4, Info: "debuginfo", Type: parse.Uint32le},
			{Offset: pos + 156, Length: 4, Info: "debuglen", Type: parse.Uint32le},
			{Offset: pos + 160, Length: 4, Info: "instpreload", Type: parse.Uint32le},
			{Offset: pos + 164, Length: 4, Info: "instdemand", Type: parse.Uint32le},
			{Offset: pos + 168, Length: 4, Info: "heapsize", Type: parse.Uint32le},
			{Offset: pos + 172, Length: 4, Info: "stacksize", Type: parse.Uint32le},
			{Offset: pos + 176, Length: 20, Info: "reserved", Type: parse.Bytes},
		}}}

	return res, nil
}
