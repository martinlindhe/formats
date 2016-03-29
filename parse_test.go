package formats

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// tests for the parse-folder

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
image data (0036), uint8
`, layout.PrettyPrint())
}
