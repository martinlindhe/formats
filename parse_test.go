package formats

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/martinlindhe/formats/parse"
	"github.com/stretchr/testify/assert"
)

// tests for the parse-folder

func TestParseBMP(t *testing.T) {

	file, err := os.Open("samples/bmp/bmp_003_WinV3.bmp")
	defer file.Close()
	assert.Equal(t, nil, err)

	b := parse.BMP(file)
	spew.Dump(b)
}
