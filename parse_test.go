package formats

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/martinlindhe/formats/parse"
	"github.com/stretchr/testify/assert"
)

// tests for the parse-folder

func TestParseARJ(t *testing.T) {

	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)
	spew.Dump(layout)
}

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
					parse.Layout{Offset: 6, Length: 4, Type: parse.Uint32le, Info: "reserved"},
					parse.Layout{Offset: 10, Length: 4, Type: parse.Uint32le, Info: "offset to image data"},
				},
			},
			parse.Layout{
				Offset: 14,
				Length: 40,
				Info:   "bmp info header V3 Win",
				Type:   parse.Group,
				Childs: []parse.Layout{
					parse.Layout{Offset: 40, Length: 4, Type: parse.Uint32le, Info: "info header size"},
					parse.Layout{Offset: 44, Length: 4, Type: parse.Uint32le, Info: "width"},
					parse.Layout{Offset: 48, Length: 4, Type: parse.Uint32le, Info: "height"},
					parse.Layout{Offset: 52, Length: 2, Type: parse.Uint16le, Info: "planes"},
					parse.Layout{Offset: 54, Length: 2, Type: parse.Uint16le, Info: "bpp"},
					parse.Layout{Offset: 56, Length: 4, Type: parse.Uint32le, Info: "compression"},
					parse.Layout{Offset: 60, Length: 4, Type: parse.Uint32le, Info: "size of picture"},
					parse.Layout{Offset: 64, Length: 4, Type: parse.Uint32le, Info: "horizontal resolution"},
					parse.Layout{Offset: 68, Length: 4, Type: parse.Uint32le, Info: "vertical resolution"},
					parse.Layout{Offset: 72, Length: 4, Type: parse.Uint32le, Info: "number of used colors"},
					parse.Layout{Offset: 76, Length: 4, Type: parse.Uint32le, Info: "number of important colors"},
				},
			},
			parse.Layout{
				Offset: 54,
				Length: 16,
				Info:   "image data",
				Type:   parse.Uint8,
			},
		},
	}, layout)
}
