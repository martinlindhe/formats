package formats

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHex(t *testing.T) {

	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)

	// reset file
	file.Seek(0, os.SEEK_SET)

	hex, err := layout.GetHex(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, "[60](fg-cyan) [ea](fg-cyan) [2b](fg-cyan) [00](fg-cyan) [22](fg-red) [0b](fg-red) [01](fg-red) [02](fg-red) [10](fg-red) [00](fg-red) [02](fg-red) [92](fg-red) [65](fg-red) [78](fg-red) [5e](fg-red) [52](fg-red)", hex)
}
