package formats

import (
	"os"
	"testing"

	"github.com/martinlindhe/formats/parse"
	"github.com/stretchr/testify/assert"
)

// tests for the parse-folder

func TestParseBMP(t *testing.T) {

	file, err := os.Open("samples/bmp/bmp_003_WinV3.bmp")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)

	assert.Equal(t, &parse.ParsedLayout{
		FormatName: "bmp",
		FileSize:   70,
		Layout: []parse.Layout{
			parse.Layout{
				Length: 14,
				Info:   "bitmap file header",
				Type:   parse.Group,
				Childs: []parse.Layout{
					parse.Layout{Offset: 0, Length: 2, Type: parse.ASCII, Info: "magic (BMP image)"},
					parse.Layout{Offset: 2, Length: 4, Type: parse.Uint32le, Info: "file size"},
					parse.Layout{Offset: 6, Length: 2, Type: parse.Uint16le, Info: "reserved 1"},
					parse.Layout{Offset: 8, Length: 2, Type: parse.Uint16le, Info: "reserved 2"},
					parse.Layout{Offset: 10, Length: 4, Type: parse.Uint32le, Info: "offset to pixel data"},
				},
			},
			parse.Layout{
				Length: 14,
				Info:   "bitmap info header",
				Type:   parse.Group,
				Childs: []parse.Layout{
					parse.Layout{Offset: 14, Length: 4, Type: parse.Uint32le, Info: "info header size"},
				},
			},
		},
	}, layout)
}
