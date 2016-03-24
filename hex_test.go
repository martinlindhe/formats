package formats

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHex(t *testing.T) {

	// XXX how do we set up a mock reader with byte data for the test?
	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	formatting.BetweenSymbols = ""
	formatting.GroupSize = 1

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)

	hex, err := layout.GetHex(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, "60ea2b00220b01021000029265785e52", hex)

	// reset file
	file.Seek(0, os.SEEK_SET)

	formatting.BetweenSymbols = " "
	formatting.GroupSize = 2

	hex, err = layout.GetHex(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, "60ea 2b00 220b 0102 1000 0292 6578 5e52", hex)
}
