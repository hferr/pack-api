package app

import (
	"context"
)

type Repo interface {
	ListPacks(ctx context.Context) (Packs, error)
	CreatePack(ctx context.Context, pack *Pack) error
}
