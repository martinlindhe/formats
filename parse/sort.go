package parse

import (
	"sort"
)

type ByLayout []Layout

func (slice ByLayout) Len() int {
	return len(slice)
}

func (slice ByLayout) Less(i, j int) bool {
	return slice[i].Offset < slice[j].Offset
}

func (slice ByLayout) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (pl *ParsedLayout) Sort() {
	sort.Sort(ByLayout(pl.Layout))
}
