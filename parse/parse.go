package parse

import (
	"os"
)

type ParseChecker struct {
	Header       []byte
	File         *os.File
	ParsedLayout ParsedLayout
}

func (pl *ParsedLayout) PercentMapped(totalSize int64) float64 {

	mapped := 0
	for _, l := range pl.Layout {
		mapped += int(l.Length)
	}
	pct := (float64(mapped) / float64(totalSize)) * 100
	return pct
}

func (pl *ParsedLayout) PercentUnmapped(totalSize int64) float64 {
	return 100 - pl.PercentMapped(totalSize)
}

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
