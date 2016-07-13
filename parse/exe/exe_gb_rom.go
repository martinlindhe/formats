package exe

// http://gbdev.gg8.se/wiki/articles/The_Cartridge_Header

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	gbLogo = []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
		0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
		0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E}

	gbSgbModes = map[byte]string{
		0x80: "no",
		0xC0: "yes",
	}
	gbCartTypes = map[byte]string{
		0:    "ROM ONLY",
		1:    "MBC1",
		2:    "MBC1+RAM",
		3:    "MBC1+RAM+BATTERY",
		5:    "MBC2",
		6:    "MBC2+BATTERY",
		8:    "ROM+RAM",
		9:    "ROM+RAM+BATTERY",
		0xb:  "MMM01",
		0xc:  "MMM01+RAM",
		0xd:  "MMM01+RAM+BATTERY",
		0xf:  "MBC3+TIMER+BATTERY",
		0x10: "MBC3+TIMER+RAM+BATTERY",
		0x11: "MBC3",
		0x12: "MBC3+RAM",
		0x13: "MBC3+RAM+BATTERY",
		0x15: "MBC4",
		0x16: "MBC4+RAM",
		0x17: "MBC4+RAM+BATTERY",
		0x19: "MBC5",
		0x1a: "MBC5+RAM",
		0x1b: "MBC5+RAM+BATTERY",
		0x1c: "MBC5+RUMBLE",
		0x1d: "MBC5+RUMBLE+RAM",
		0x1e: "MBC5+RUMBLE+RAM+BATTERY",
		0xfc: "POCKET CAMERA",
		0xfd: "BANDAI TAMA5",
		0xfe: "HuC3",
		0xff: "HuC1+RAM+BATTERY",
	}
	gbROMSizes = map[byte]string{
		0:    "32KByte (no ROM banking)",
		1:    "64KByte (4 banks)",
		2:    "128KByte (8 banks)",
		3:    "256KByte (16 banks)",
		4:    "512KByte (32 banks)",
		5:    "1MByte (64 banks) - only 63 banks used by MBC1",
		6:    "2MByte (128 banks) - only 125 banks used by MBC1",
		7:    "4MByte (256 banks)",
		0x52: "1.1MByte (72 banks)",
		0x53: "1.2MByte (80 banks)",
		0x54: "1.5MByte (96 banks)",
	}
	gbRAMSizes = map[byte]string{
		0: "none",
		1: "2 KBytes",
		2: "8 Kbytes",
		3: "32 KBytes (4 banks of 8KBytes each)",
		4: "128 KBytes (16 banks of 8KBytes each)",
		5: "64 KBytes (8 banks of 8KBytes each)",
	}
	gbDestCodes = map[byte]string{
		0: "Japanese",
		1: "Non-Japanese",
	}
)

// GameboyROM parses the Game Boy ROM image format
func GameboyROM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGameboyROM(c.Header) {
		return nil, nil
	}
	return parseGameboyROM(c.File, c.ParsedLayout)
}

func isGameboyROM(b []byte) bool {

	if !parse.HasSignatureInHeader(b, 0x104, gbLogo) {
		return false
	}
	return true
}

func parseGameboyROM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0x100)
	sgb, _ := parse.ReadToMap(file, parse.Uint8, pos+70, gbSgbModes)
	cartType, _ := parse.ReadToMap(file, parse.Uint8, pos+71, gbCartTypes)
	romSize, _ := parse.ReadToMap(file, parse.Uint8, pos+72, gbROMSizes)
	ramSize, _ := parse.ReadToMap(file, parse.Uint8, pos+73, gbRAMSizes)
	destCode, _ := parse.ReadToMap(file, parse.Uint8, pos+74, gbDestCodes)

	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 80, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "entry point", Type: parse.Bytes}, // XXX mark as Code or sth
			{Offset: pos + 4, Length: 48, Info: "nintendo logo", Type: parse.Bytes},
			{Offset: pos + 52, Length: 16, Info: "title", Type: parse.ASCIIZ},
			{Offset: pos + 68, Length: 2, Info: "new licensee code", Type: parse.ASCII},
			{Offset: pos + 70, Length: 1, Info: "sgb flag = " + sgb, Type: parse.Uint8},
			{Offset: pos + 71, Length: 1, Info: "cartridge type = " + cartType, Type: parse.Uint8},
			{Offset: pos + 72, Length: 1, Info: "rom size = " + romSize, Type: parse.Uint8},
			{Offset: pos + 73, Length: 1, Info: "ram size = " + ramSize, Type: parse.Uint8},
			{Offset: pos + 74, Length: 1, Info: "destination = " + destCode, Type: parse.Uint8},
			{Offset: pos + 75, Length: 1, Info: "old license code", Type: parse.Uint8},
			{Offset: pos + 76, Length: 1, Info: "mask rom version", Type: parse.Uint8},
			{Offset: pos + 77, Length: 1, Info: "header checksum", Type: parse.Uint8},
			{Offset: pos + 78, Length: 2, Info: "gloal checksum", Type: parse.Uint16le},
		}}}

	return &pl, nil
}
