package parse

// OS/2 Lineear eXecutable
// STATUS: 10%

import (
	"os"
)

// http://www.virtualbox.org/svn/kstuff-mirror/trunk/include/k/kLdrFmts/lx.h

// parse OS/2 Linear eXecutable header
func parseMZ_LXHeader(file *os.File, pos int64) ([]Layout, error) {

	res := []Layout{}

	res = append(res, Layout{
		Offset: pos,
		Length: 196,
		Info:   "LX header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: pos, Length: 2, Info: "identifier", Type: ASCII},
			Layout{Offset: pos + 2, Length: 1, Info: "border", Type: Uint8},
			Layout{Offset: pos + 3, Length: 1, Info: "worder", Type: Uint8},
			Layout{Offset: pos + 4, Length: 4, Info: "level", Type: Uint32le},
			Layout{Offset: pos + 8, Length: 2, Info: "cpu", Type: Uint16le},
			Layout{Offset: pos + 10, Length: 2, Info: "os", Type: Uint16le},
			Layout{Offset: pos + 12, Length: 4, Info: "version", Type: MajorMinor32le},
			Layout{Offset: pos + 16, Length: 4, Info: "mflags", Type: Uint32le},
			Layout{Offset: pos + 20, Length: 4, Info: "mpages", Type: Uint32le},
			Layout{Offset: pos + 24, Length: 4, Info: "startobj", Type: Uint32le},
			Layout{Offset: pos + 28, Length: 4, Info: "eip", Type: Uint32le},
			Layout{Offset: pos + 32, Length: 4, Info: "stackobj", Type: Uint32le},
			Layout{Offset: pos + 36, Length: 4, Info: "esp", Type: Uint32le},
			Layout{Offset: pos + 40, Length: 4, Info: "pagesize", Type: Uint32le},
			Layout{Offset: pos + 44, Length: 4, Info: "pageshift", Type: Uint32le},
			Layout{Offset: pos + 48, Length: 4, Info: "fixup size", Type: Uint32le},
			Layout{Offset: pos + 52, Length: 4, Info: "fixup sum", Type: Uint32le},
			Layout{Offset: pos + 56, Length: 4, Info: "ldrsize", Type: Uint32le},
			Layout{Offset: pos + 60, Length: 4, Info: "ldrsum", Type: Uint32le},
			Layout{Offset: pos + 64, Length: 4, Info: "objtab", Type: Uint32le},
			Layout{Offset: pos + 68, Length: 4, Info: "objcnt", Type: Uint32le},
			Layout{Offset: pos + 72, Length: 4, Info: "objmap", Type: Uint32le},
			Layout{Offset: pos + 76, Length: 4, Info: "itermap", Type: Uint32le},
			Layout{Offset: pos + 80, Length: 4, Info: "rsrctab", Type: Uint32le},
			Layout{Offset: pos + 84, Length: 4, Info: "rsrccnt", Type: Uint32le},
			Layout{Offset: pos + 88, Length: 4, Info: "restab", Type: Uint32le},
			Layout{Offset: pos + 92, Length: 4, Info: "enttab", Type: Uint32le},
			Layout{Offset: pos + 96, Length: 4, Info: "dirtab", Type: Uint32le},
			Layout{Offset: pos + 100, Length: 4, Info: "dircnt", Type: Uint32le},
			Layout{Offset: pos + 104, Length: 4, Info: "fpagetab", Type: Uint32le},
			Layout{Offset: pos + 108, Length: 4, Info: "frectab", Type: Uint32le},
			Layout{Offset: pos + 112, Length: 4, Info: "impmod", Type: Uint32le},
			Layout{Offset: pos + 116, Length: 4, Info: "impmodcnt", Type: Uint32le},
			Layout{Offset: pos + 120, Length: 4, Info: "impproc", Type: Uint32le},
			Layout{Offset: pos + 124, Length: 4, Info: "pagesum", Type: Uint32le},
			Layout{Offset: pos + 128, Length: 4, Info: "datapage", Type: Uint32le},
			Layout{Offset: pos + 132, Length: 4, Info: "preload", Type: Uint32le},
			Layout{Offset: pos + 136, Length: 4, Info: "nrestab", Type: Uint32le},
			Layout{Offset: pos + 140, Length: 4, Info: "cbnrestab", Type: Uint32le},
			Layout{Offset: pos + 144, Length: 4, Info: "nressum", Type: Uint32le},
			Layout{Offset: pos + 148, Length: 4, Info: "autodata", Type: Uint32le},
			Layout{Offset: pos + 152, Length: 4, Info: "debuginfo", Type: Uint32le},
			Layout{Offset: pos + 156, Length: 4, Info: "debuglen", Type: Uint32le},
			Layout{Offset: pos + 160, Length: 4, Info: "instpreload", Type: Uint32le},
			Layout{Offset: pos + 164, Length: 4, Info: "instdemand", Type: Uint32le},
			Layout{Offset: pos + 168, Length: 4, Info: "heapsize", Type: Uint32le},
			Layout{Offset: pos + 172, Length: 4, Info: "stacksize", Type: Uint32le},
			Layout{Offset: pos + 176, Length: 20, Info: "reserved", Type: Bytes},
		}})

	return res, nil
}
