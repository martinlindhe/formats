package formats

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getFileSize(file *os.File) int64 {

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}

// return file extension, without leading dot
func fileExt(file *os.File) string {

	ext := filepath.Ext(file.Name())
	if len(ext) > 0 {
		ext = ext[1:]
	}
	return ext
}

func byteSliceEquals(a []byte, b []byte) bool {

	if len(a) != len(b) {
		fmt.Println("error: a has len", len(a), " and b has len ", len(b))
		return false
	}

	for i, c1 := range a {
		c2 := b[i]
		if c1 != c2 {
			return false
		}
	}
	return true
}
