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

	// reset file
	file.Seek(0, os.SEEK_SET)

	hex, err := layout.GetHex(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, "[60](fg-blue)[ea](fg-blue)[2b](fg-blue)[00](fg-blue)[22](fg-red)[0b](fg-red)[01](fg-red)[02](fg-red)[10](fg-red)[00](fg-red)[02](fg-red)[92](fg-red)[65](fg-red)[78](fg-red)[5e](fg-red)[52](fg-red)", hex)

	// reset file
	file.Seek(0, os.SEEK_SET)

	formatting.BetweenSymbols = " "
	formatting.GroupSize = 2

	hex, err = layout.GetHex(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, "[60](fg-blue)[ea](fg-blue) [2b](fg-blue)[00](fg-blue) [22](fg-red)[0b](fg-red) [01](fg-red)[02](fg-red) [10](fg-red)[00](fg-red) [02](fg-red)[92](fg-red) [65](fg-red)[78](fg-red) [5e](fg-red)[52](fg-red)", hex)
}
