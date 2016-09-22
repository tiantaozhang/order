package orderIntegral

import (
	"sort"
)

type kvpair struct {
	k, v string
}
type kvpairs []kvpair

func (k kvpairs) Less(i, j int) bool {
	return k[i].k < k[j].k
}

func (k kvpairs) Swap(i, j int) {
	k[j], k[i] = k[i], k[j]
}

func (k kvpairs) Sort() {
	sort.Sort(k)
}

func (k kvpairs) Len() int {
	return len(k)
}
