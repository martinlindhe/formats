package formats

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGif87a(t *testing.T) {

	// NOTE this, among others, tests the "87a,89a" matching of possible values in formats/gif.yml

	file, err := os.Open("samples/gif/gif_001_87a.gif")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, &ParsedLayout{
		FormatName: "gif",
		Layout: []Layout{
			Layout{0, 3, ASCII, "magic (GIF image)"},
			Layout{3, 3, ASCII, "version"},
			Layout{6, 2, Uint16le, "width"},
			Layout{8, 2, Uint16le, "height"},
			Layout{10, 1, Uint8, "packed"},
			Layout{11, 1, Uint8, "background color"},
			Layout{12, 1, Uint8, "aspect ratio"},
		},
	}, layout)
}

func TestParseGif89a(t *testing.T) {

	file, err := os.Open("samples/gif/gif_002_89a.gif")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, &ParsedLayout{
		FormatName: "gif",
		Layout: []Layout{
			Layout{0, 3, ASCII, "magic (GIF image)"},
			Layout{3, 3, ASCII, "version"},
			Layout{6, 2, Uint16le, "width"},
			Layout{8, 2, Uint16le, "height"},
			Layout{10, 1, Uint8, "packed"},
			Layout{11, 1, Uint8, "background color"},
			Layout{12, 1, Uint8, "aspect ratio"},
		},
	}, layout)
}
