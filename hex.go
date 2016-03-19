package formats

import (
	"fmt"
	"reflect"
	"strings"
)

// HexFormatting ...
type HexFormatting struct {
	BetweenSymbols string
	GroupSize      byte
}

var (
	formatting = HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      1,
	}
)

// Formatting ...
func Formatting(fmt HexFormatting) { formatting = fmt }

// CombineHexRow ...
func CombineHexRow(symbols []string) string {

	group := []string{}
	row := []string{}
	cur := byte(0)

	for _, sym := range symbols {
		cur++
		group = append(group, sym)
		if cur == formatting.GroupSize {
			row = append(row, strings.Join(group, ""))
			group = nil
			cur = 0
		}
	}
	return strings.Join(row, formatting.BetweenSymbols)
}

// Layout represents a parsed file structure layout as a flat list
type Layout struct {
	Offset int64
	Length byte
	Type   DataType
	Info   string
}

// DataType ...
type DataType int

func (dt DataType) String() string {

	m := map[DataType]string{
		ASCIIZ:   "ASCIIZ",
		Byte:     "byte",
		Uint16le: "uint16-le",
		Uint32le: "uint32-le",
		Int16le:  "int16-le",
		Int32le:  "int32-le",
	}

	if val, ok := m[dt]; ok {
		return val
	}

	// NOTE should only be able to panic during dev (as in:
	// adding a new datatype and forgetting to add it to the map)
	panic(dt)
}

// ...
const (
	_               = iota
	ASCIIZ DataType = iota
	Byte
	Uint16le
	Uint32le
	Int16le
	Int32le
)

func structToFlatStruct(obj interface{}) []Layout { // XXX implement

	res := []Layout{}

	//	spew.Dump(x)

	// XXX iterate over struct, create a 2d rep of the structure mapping

	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		// XXX is it a struct ?
		// fmt.Println(f.Tag)

		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	return res
}
