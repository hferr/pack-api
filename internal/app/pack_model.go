package app

import (
	"sort"

	"github.com/google/uuid"
)

type Pack struct {
	ID   uuid.UUID
	Size int
}

func NewPack(size int) *Pack {
	return &Pack{
		ID:   uuid.New(),
		Size: size,
	}
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
