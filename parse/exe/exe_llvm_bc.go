package exe

// STATUS: 1%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

// LLVMBitcode parses the LLVM Bit code format
func LLVMBitcode(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLLVMBitcode(c.Header) {
		return nil, nil
	}
	return parseLLVMBitcode(c.File, c.ParsedLayout)
}

func isLLVMBitcode(b []byte) bool {

	return binary.LittleEndian.Uint32(b) == 0x0b17c0de
}

func parseLLVMBitcode(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
