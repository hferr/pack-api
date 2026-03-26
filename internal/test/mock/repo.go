package mock

import (
	"context"

	"github.com/hferr/pack-api/internal/app"
)

type Repo struct {
	ListPacksFn  func(ctx context.Context) (app.Packs, error)
	CreatePackFn func(ctx context.Context, pack *app.Pack) error
}

func (r *Repo) ListPacks(ctx context.Context) (app.Packs, error) {
	return r.ListPacksFn(ctx)
}

func (r *Repo) CreatePack(ctx context.Context, pack *app.Pack) error {
	return r.CreatePackFn(ctx, pack)
}
