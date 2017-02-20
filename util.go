package formats

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// return file extension, without leading dot
func fileExt(file *os.File) string {

	ext := filepath.Ext(file.Name())
	if len(ext) > 0 {
		ext = ext[1:]
	}
	return ext
}

func fileSize(file *os.File) (int64, error) {

	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
