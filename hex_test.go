package formats

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHex(t *testing.T) {

	// XXX how do we set up a mock reader with byte data for the test?
	file, _ := os.Open("samples/tiny.arj")
	defer file.Close()

	reader := io.Reader(file)

	// XXX get console screen height
	height := 2

	formatting.betweenSymbols = ""
	formatting.groupSize = 1

	hex, err := GetHex(&reader, height)
	assert.Equal(t, nil, err)

	assert.Equal(t, []string{
		"60ea2b00220b01021000029265785e52",
		"65785e52000000000000000000000000",
	}, hex)

	// reset file
	file.Seek(0, os.SEEK_SET)

	formatting.betweenSymbols = " "
	formatting.groupSize = 2

	hex, err = GetHex(&reader, height)
	assert.Equal(t, nil, err)

	assert.Equal(t, []string{
		"60ea 2b00 220b 0102 1000 0292 6578 5e52",
		"6578 5e52 0000 0000 0000 0000 0000 0000",
	}, hex)
}
