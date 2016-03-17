package formats

import (
	"fmt"
	"os"
)

func psxBiosProbe(f *os.File) bool {
	fmt.Printf("psx-bios probe\n")

	return false
}
