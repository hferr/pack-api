package app

import "sort"

type Pack struct {
	Size int
}

type Packs []Pack

func (packs Packs) GetSortedSizes() []int {
	sizes := make([]int, len(packs))
	for i, pack := range packs {
		sizes[i] = pack.Size
	}

	sort.Ints(sizes)
	return sizes
}
