package formats

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/martinlindhe/formats/parse"
	"github.com/stretchr/testify/assert"
)

// some tests to see that parsed files look ok
func TestParsedLayout(t *testing.T) {

	searchDir := "./samples"

	err := filepath.Walk(searchDir, func(path string, fi os.FileInfo, err error) error {

		if fi.IsDir() {
			return nil
		}

		fmt.Println("OPEN", path)
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}

		layout := ParseLayout(f)
		if layout == nil {
			// NOTE since not all samples is currently supported
			fmt.Println("FAIL to parse layout from", path)
			return nil
		}

		for _, l := range layout.Layout {
			if l.Type != parse.Group {
				t.Fatalf("root level must be group %v, %s", l, path)
			}

			if l.Type == parse.RGB && l.Length != 3 {
				t.Fatalf("RGB field must be %d bytes, was %d", 3, l.Length)
			}

			if l.Type == parse.Uint16le && l.Length != 2 {
				t.Fatalf("Uint16le field must be %d bytes, was %d", 2, l.Length)
			}

			if l.Type == parse.Uint32le && l.Length != 4 {
				t.Fatalf("Uint16le field must be %d bytes, was %d", 4, l.Length)
			}

			if len(l.Childs) > 0 && l.Childs[0].Offset != l.Offset {
				t.Fatalf("%s child 0 offset should be same as parent %04x, but is %04x", l.Info, l.Offset, l.Childs[0].Offset)
			}

			sum := int64(0)
			for _, child := range l.Childs {
				sum += child.Length
				if child.Type == parse.Group {
					t.Fatalf("child level cant be group %v, %s", l, path)
				}
			}
			if sum != l.Length {
				t.Fatalf("child sum for %s, field %s is %d, but group length is %d, %v", path, l.Info, sum, l.Length, l)
			}
		}

		return nil
	})
	assert.Equal(t, nil, err)
}

/*
func TestParseGIF87a(t *testing.T) {

	file, err := os.Open("samples/gif/gif_001_87a.gif")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)
	assert.Equal(t, true, layout != nil)
	assert.Equal(t, `XXX`, layout)
}

func TestParseGIF89a(t *testing.T) {

	file, err := os.Open("samples/gif/gif_002_89a.gif")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)
	assert.Equal(t, true, layout != nil)
	assert.Equal(t, `XXX`, layout)
}
*/

func TestParseARJ(t *testing.T) {

	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)
	assert.Equal(t, `arj main header (0035), Group
  magic (ARJ archive) (0035), uint16-le
  basic header size (0037), uint16-le
  size up to and including 'extra data' (0039), uint8
  archiver version number (003a), uint8
  minimum archiver version to extract (003b), uint8
  host OS (003c), uint8
  arj flags (003d), uint8
  security version (003e), uint8
  file type (003f), uint8
  created time (0040), uint32-le
  modified time (0044), uint32-le
  archive size for secured archive (0048), uint32-le
  security envelope file position (004c), uint32-le
  filespec position in filename (0050), uint32-le
  length in bytes of security envelope data (0054), uint16-le
  encryption version (0056), uint8
  last chapter (0057), uint8
  archive name (0035), ASCIIZ
  comment (0035), ASCIIZ
  crc32 (0035), uint32-le
  ext header size (0039), uint32-le
`, layout.PrettyPrint())
}

func TestParseBMP(t *testing.T) {

	file, err := os.Open("samples/bmp/bmp_003_WinV3.bmp")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout := ParseLayout(file)

	assert.Equal(t, `bitmap file header (0000), Group
  magic (BMP image) (0000), ASCII
  file size (0002), uint32-le
  reserved (0006), uint32-le
  offset to image data (000a), uint32-le
bmp info header V3 Win (000e), Group
  info header size (0028), uint32-le
  width (002c), uint32-le
  height (0030), uint32-le
  planes (0034), uint16-le
  bpp (0036), uint16-le
  compression (0038), uint32-le
  size of picture (003c), uint32-le
  horizontal resolution (0040), uint32-le
  vertical resolution (0044), uint32-le
  number of used colors (0048), uint32-le
  number of important colors (004c), uint32-le
image data (0036), Group
  image data (0036), uint8
`, layout.PrettyPrint())
}
