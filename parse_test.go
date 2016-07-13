package formats

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/martinlindhe/formats/parse"
	"github.com/stretchr/testify/assert"
)

// some tests to see that parsed files look ok
func TestParsedLayout(t *testing.T) {

	searchDir := "./samples/font"

	err := filepath.Walk(searchDir, func(path string, fi os.FileInfo, err error) error {

		if fi == nil {
			t.Fatalf("invalid path " + searchDir)
		}

		if fi.IsDir() {
			return nil
		}

		if fi.Size() == 0 {
			return nil
		}

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}

		layout, err := ParseLayout(f)
		if err != nil {
			t.Fatalf("got error %v", err)
		}

		assert.Equal(t, true, layout != nil)
		assert.Equal(t, true, layout.FileKind != 0)

		if layout == nil {
			// NOTE since not all samples is currently supported
			t.Log("warning: failed to parse layout from", path)
			return nil
		}

		// make sure all parsed layouts got FileKind set
		assert.Equal(t, false, layout.FileKind == 0)

		if layout.MimeType == "" {
			/*
				// ask "file" about mime type
				filemagicMime, _ := runCommandReturnStdout("file --mime-type " + path)
				filemagicMime = strings.TrimSpace(filemagicMime)
				res := strings.Split(filemagicMime, " ")
				mime := ""
				if len(res) > 1 {
					mime = res[1]
				} else {
					mime = filemagicMime
				}
				if filemagicMime != "" {
					t.Log("warning:", layout.FormatName, "has no mime. file suggests", mime)
				} else {
					t.Log("warning:", layout.FormatName, "has no mime")
				}
			*/
		}

		for _, l := range layout.Layout {
			if l.Type != parse.Group {
				t.Fatalf("%s:%s in header %s: root level must be group", layout.FormatName, path, l.Info)
			}
			if len(l.Masks) > 0 {
				t.Fatalf("%s:%s in header %s: can not have masks on root level group %v", layout.FormatName, path, l.Info, l)
			}
			if l.Type == parse.RGB && l.Length != 3 {
				t.Fatalf("%s:%s in header %s: RGB field must be %d bytes, was %d", layout.FormatName, path, l.Info, 3, l.Length)
			}
			if l.Type == parse.Uint8 && l.Length != 1 {
				t.Fatalf("%s:%s in header %s: Uint8 field must be %d bytes, was %d", layout.FormatName, path, l.Info, 1, l.Length)
			}
			if l.Type == parse.Bytes && l.Length == 1 {
				t.Fatalf("%s:%s in header %s: Bytes field should never be used for single-byte fields", layout.FormatName, path, l.Info)
			}
			if l.Type == parse.Uint16le && l.Length != 2 {
				t.Fatalf("%s:%s in header %s: Uint16le field must be %d bytes, was %d", layout.FormatName, path, l.Info, 2, l.Length)
			}
			if l.Type == parse.Uint32le && l.Length != 4 {
				t.Fatalf("%s:%s in header %s: Uint16le field must be %d bytes, was %d", layout.FormatName, path, l.Info, 4, l.Length)
			}
			if len(l.Childs) > 0 && l.Childs[0].Offset != l.Offset {
				t.Fatalf("%s:%s in header %s: %s child 0 offset should be same as parent %04x, but is %04x", layout.FormatName, path, l.Info, l.Info, l.Offset, l.Childs[0].Offset)
			}
			if l.Offset+l.Length > layout.FileSize {
				t.Fatalf("%s:%s in header %s: %s data extends above end of file with %d bytes", layout.FormatName, path, l.Info, l.Info, layout.FileSize-(l.Offset+l.Length))
			}
			sum := int64(0)
			for _, child := range l.Childs {
				sum += child.Length
				if child.Type == parse.Group {
					t.Fatalf("%s:%s in header %s: child level cant be group %s:%s", layout.FormatName, path, l.Info, l.Info, child.Info)
				}
				if child.Offset+child.Length > layout.FileSize {
					t.Fatalf("%s:%s in header %s: %s:%s (child) data extends above end of file with %d bytes", layout.FormatName, path, l.Info, l.Info, child.Info, layout.FileSize-(l.Offset+l.Length))
				}

				if len(child.Masks) > 0 {
					expectedTot := child.GetBitSize()
					tot := 0
					for _, mask := range child.Masks {
						tot += mask.Length
					}
					if tot != expectedTot {
						t.Fatalf("%s:%s in header %s: %s:%s (child) masks size dont add up. expected %d bits, got %d", layout.FormatName, path, l.Info, l.Info, child.Info, expectedTot, tot)
					}
				}
			}
			if sum != l.Length {
				t.Fatalf("%s:%s in header %s: child sum for %s is %d, but group length is %d", layout.FormatName, path, l.Info, l.Info, sum, l.Length)
			}
		}
		return nil
	})
	assert.Equal(t, nil, err)
}

func TestParseGIFBitFields(t *testing.T) {

	file, err := os.Open("samples/image/gif/gif_89a_004_fish.gif")
	defer file.Close()
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	pl, err := ParseLayout(file)
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	l0 := pl.DecodeBitfieldFromInfo(file, "global color table size")
	l3 := pl.DecodeBitfieldFromInfo(file, "sort flag")
	l4 := pl.DecodeBitfieldFromInfo(file, "color resolution")
	l7 := pl.DecodeBitfieldFromInfo(file, "global color table flag")

	assert.Equal(t, l0, uint32(2))
	assert.Equal(t, l3, uint32(0))
	assert.Equal(t, l4, uint32(7))
	assert.Equal(t, l7, uint32(1))
}

func TestParseGIF87a(t *testing.T) {

	file, err := os.Open("samples/image/gif/gif_87a_001.gif")
	defer file.Close()
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	layout, err := ParseLayout(file)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	assert.Equal(t, `Format: gif (gif_87a_001.gif, 42 bytes)

header (0000), group
  signature (0000), ASCII
  version (0003), ASCII
logical screen descriptor (0006), group
  width (0006), uint16-le
  height (0008), uint16-le
  packed (000a), uint8
  background color (000b), uint8
  aspect ratio (000c), uint8
global color table (000d), group
  color 1 (000d), RGB
  color 2 (0010), RGB
  color 3 (0013), RGB
  color 4 (0016), RGB
image descriptor (0019), group
  image separator (0019), uint8
  image left (001a), uint16-le
  image top (001c), uint16-le
  image width (001e), uint16-le
  image height (0020), uint16-le
  packed #3 (0022), uint8
image data (0023), group
  lzw code size (0023), uint8
  lzw block size (0024), uint8
  lzw block (0025), bytes
  lzw block size (0028), uint8
trailer (0029), group
  trailer (0029), uint8
`, layout.PrettyPrint())
}

func TestParseGIF89a(t *testing.T) {

	file, err := os.Open("samples/image/gif/gif_89a_001.gif")
	defer file.Close()
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	layout, err := ParseLayout(file)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	assert.Equal(t, `Format: gif (gif_89a_001.gif, 69 bytes)

header (0000), group
  signature (0000), ASCII
  version (0003), ASCII
logical screen descriptor (0006), group
  width (0006), uint16-le
  height (0008), uint16-le
  packed (000a), uint8
  background color (000b), uint8
  aspect ratio (000c), uint8
global color table (000d), group
  color 1 (000d), RGB
  color 2 (0010), RGB
  color 3 (0013), RGB
  color 4 (0016), RGB
graphic control extension (0019), group
  block id (extension) (0019), uint8
  type = graphic control (001a), uint8
  byte size (001b), uint8
  packed #2 (001c), uint8
  delay time (001d), uint16-le
  transparent color index (001f), uint8
  block terminator (0020), uint8
image descriptor (0021), group
  image separator (0021), uint8
  image left (0022), uint16-le
  image top (0024), uint16-le
  image width (0026), uint16-le
  image height (0028), uint16-le
  packed #3 (002a), uint8
image data (002b), group
  lzw code size (002b), uint8
  lzw block size (002c), uint8
  lzw block (002d), bytes
  lzw block size (0043), uint8
trailer (0044), group
  trailer (0044), uint8
`, layout.PrettyPrint())

	// perform some bitfield tests on this known file
	assert.Equal(t, uint32(1), layout.DecodeBitfieldFromInfo(file, "global color table flag"))
	assert.Equal(t, uint32(0), layout.DecodeBitfieldFromInfo(file, "local color table flag"))
}

/*
func TestParseARJ(t *testing.T) {

	file, err := os.Open("samples/archive/arj/arj_001_tiny.arj")
	defer file.Close()
	if err != nil {
		t.Fatalf("got error %v", err)
	}


	layout, err := ParseLayout(file)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	assert.Equal(t, `Format: arj (arj_001_tiny.arj, 121 bytes)

main header (0024), group
  magic (0024), uint16-le
  basic header size (0026), uint16-le
  size up to and including 'extra data' (0028), uint8
  archiver version number (0029), uint8
  minimum archiver version to extract (002a), uint8
  host OS (002b), uint8
  arj flags (002c), uint8
  security version (002d), uint8
  file type (002e), uint8
  created time (002f), uint32-le
  modified time (0033), uint32-le
  archive size for secured archive (0037), uint32-le
  security envelope file position (003b), uint32-le
  filespec position in filename (003f), uint32-le
  length in bytes of security envelope data (0043), uint16-le
  encryption version (0045), uint8
  last chapter (0046), uint8
  archive name (0047), ASCIIZ
  comment (0048), ASCIIZ
  crc32 (0049), uint32-le
  ext header size (004d), uint32-le
`, layout.PrettyPrint())
}
*/

func TestParseBMP(t *testing.T) {

	file, err := os.Open("samples/image/bmp/bmp_v3-001.bmp")
	defer file.Close()
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	layout, err := ParseLayout(file)
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	assert.Equal(t, `Format: bmp (bmp_v3-001.bmp, 70 bytes)

file header (0000), group
  magic (0000), ASCII
  file size (0002), uint32-le
  reserved (0006), uint32-le
  image data offset (000a), uint32-le
info header V3 (000e), group
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
image data (0036), group
  image data (0036), bytes
`, layout.PrettyPrint())
}

func runCommandReturnStdout(step string) (string, error) {

	parts := strings.Split(step, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	res := ""

	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return res, err
	}
	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			res += string(stdOutScanner.Bytes()) + "\n"
		}
	}()

	if err := cmd.Start(); err != nil {
		return res, err
	}
	if err := cmd.Wait(); err != nil {
		return res, err
	}
	return res, nil
}
