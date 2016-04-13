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

	searchDir := "./samples/archives/cab"

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

		layout, err := ParseLayout(f)
		assert.Equal(t, nil, err)

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
			if l.Type == parse.Uint8 && l.Length != 1 {
				t.Fatalf("Uint8 field must be %d bytes, was %d", 1, l.Length)
			}
			if l.Type == parse.Bytes && l.Length == 1 {
				t.Fatalf("Bytes field should never be used for single-byte fields")
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
			if l.Offset+l.Length > layout.FileSize {
				t.Fatalf("%s child extends above end of file with %d bytes", l.Info, layout.FileSize-(l.Offset+l.Length))
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
*/

func TestParseGIF89a(t *testing.T) {

	file, err := os.Open("samples/gif/gif_89a_001.gif")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, `header (0000), Group
  signature (0000), ASCII
  version (0003), ASCII
logical screen descriptor (0006), Group
  width (0006), uint16-le
  height (0008), uint16-le
  packed (000a), uint8
  background color (000b), uint8
  aspect ratio (000c), uint8
global color table (000d), Group
  color 1 (000d), RGB
  color 2 (0010), RGB
  color 3 (0013), RGB
  color 4 (0016), RGB
extension (0019), Group
  block id (extension) (0019), uint8
  graphic control (001a), uint8
  byte size (001b), uint8
  packed #2 (001c), uint8
  delay time (001d), uint16-le
  transparent color index (001f), uint8
  block terminator (0020), uint8
image descriptor (0021), Group
  image separator (0021), uint8
  image left (0022), uint16-le
  image top (0024), uint16-le
  image width (0026), uint16-le
  image height (0028), uint16-le
  packed #3 (002a), uint8
image data (002b), Group
  lzw code size (002b), uint8
  block length (002c), uint8
  block (002d), uint8
  block length (0043), uint8
trailer (0044), Group
  trailer (0044), uint8
`, layout.PrettyPrint())

	// perform some bitfield tests on this known file
	assert.Equal(t, uint32(1), layout.DecodeBitfieldFromInfo(file, "global color table flag"))
	assert.Equal(t, uint32(0), layout.DecodeBitfieldFromInfo(file, "local color table flag"))
}

func TestParseARJ(t *testing.T) {

	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, `main header (0035), Group
  magic (0035), uint16-le
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

	file, err := os.Open("samples/bmp/bmp_v3-001.bmp")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)

	assert.Equal(t, `file header (0000), Group
  magic (0000), ASCII
  file size (0002), uint32-le
  reserved (0006), uint32-le
  offset to image data (000a), uint32-le
info header V3 (000e), Group
  info header size (000e), uint32-le
  width (0012), uint32-le
  height (0016), uint32-le
  planes (001a), uint16-le
  bpp (001c), uint16-le
  compression = rgb (001e), uint32-le
  size of picture (0022), uint32-le
  horizontal resolution (0026), uint32-le
  vertical resolution (002a), uint32-le
  number of used colors (002e), uint32-le
  number of important colors (0032), uint32-le
image data (0036), Group
  image data (0036), uint8
`, layout.PrettyPrint())
}
