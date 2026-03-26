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

func (p *Pack) Validate() error {
	if p.Size < 1 {
		return &ValidationError{
			Field:   "Size",
			Message: "has to be greater than 0",
		}
	}
	return nil
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
