package sort

import (
	"sort"
	"testing"
)

func TestSort_FindAndSearch(t *testing.T) {
	sli := []int{1, 4, 2, 7, 5}
	target := 5

	// sort.Findはバイナリサーチを実施するので、ソート済スライスを渡す必要がある
	sort.SliceStable(sli, func(i, j int) bool {
		return sli[i] < sli[j]
	})

	t.Log(sli)

	// Find:
	i, found := sort.Find(len(sli), func(i int) int {
		if target > sli[i] {
			return 1
		}
		if target < sli[i] {
			return -1
		}
		return 0
	})
	t.Log(i, found)

	// Search:
	i = sort.Search(len(sli), func(i int) bool {
		return sli[i] >= target
	})
	t.Logf("search result: %d", i)
}
