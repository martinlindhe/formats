package formats

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXXX(t *testing.T) {

	// XXX we fake result from structToFlatStruct() to test presentation
	x := map[uint64]Layout{0x0000: Layout{1, ASCIIZ, "hej"}}

	// XXXx
	fmt.Println(x)
}

func TestGetHex(t *testing.T) {

	// XXX how do we set up a mock reader with byte data for the test?
	file, _ := os.Open("samples/tiny.arj")
	defer file.Close()

	reader := io.Reader(file)

	formatting.BetweenSymbols = ""
	formatting.GroupSize = 1

	hex, err := GetHex(&reader)
	assert.Equal(t, nil, err)

	assert.Equal(t, "60ea2b00220b01021000029265785e52", hex)

	// reset file
	file.Seek(0, os.SEEK_SET)

	formatting.BetweenSymbols = " "
	formatting.GroupSize = 2

	hex, err = GetHex(&reader)
	assert.Equal(t, nil, err)

	assert.Equal(t, "60ea 2b00 220b 0102 1000 0292 6578 5e52", hex)
}
