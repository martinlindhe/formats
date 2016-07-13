package parse

import (
	"os"
)

// Checker (parse checker)
type Checker struct {
	Header       []byte
	File         *os.File
	ParsedLayout ParsedLayout
}

// PercentMapped returns the % of total file size mapped to known structures
func (pl *ParsedLayout) PercentMapped(totalSize int64) float64 {

	mapped := 0
	for _, l := range pl.Layout {
		mapped += int(l.Length)
	}
	pct := (float64(mapped) / float64(totalSize)) * 100
	return pct
}

// PercentUnmapped returns the % of total file size not mapped to known structures
func (pl *ParsedLayout) PercentUnmapped(totalSize int64) float64 {
	return 100 - pl.PercentMapped(totalSize)
}

// HasSignatureInHeader returns true if `sig` is found in `hdr`
func HasSignatureInHeader(hdr []byte, pos int64, sig []byte) bool {

	for i := int64(0); i < int64(len(sig)); i++ {

		ofs := pos + i
		if ofs >= int64(len(hdr)) {
			return false
		}
		if hdr[ofs] != sig[i] {
			return false
		}
	}
	return true
}
