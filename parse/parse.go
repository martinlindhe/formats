package parse

import (
	"os"
)

type ParseChecker struct {
	Header       [0xffff]byte
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
