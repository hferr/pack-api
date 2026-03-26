package mock

import (
	"context"

	"github.com/hferr/pack-api/internal/app"
)

type PackService struct {
	ListPacksFn             func(ctx context.Context) (app.Packs, error)
	CreatePackFn            func(ctx context.Context, size int) (*app.Pack, error)
	RebuildPacksFn          func(ctx context.Context, sizes []int) (app.Packs, error)
	CalculateMinPackOrderFn func(ctx context.Context, items int) (map[int]int, error)
}

func (s *PackService) ListPacks(ctx context.Context) (app.Packs, error) {
	return s.ListPacksFn(ctx)
}

func (s *PackService) CreatePack(ctx context.Context, size int) (*app.Pack, error) {
	return s.CreatePackFn(ctx, size)
}

func (s *PackService) RebuildPacks(ctx context.Context, sizes []int) (app.Packs, error) {
	return s.RebuildPacksFn(ctx, sizes)
}

func (s *PackService) CalculateMinPackOrder(ctx context.Context, items int) (map[int]int, error) {
	return s.CalculateMinPackOrderFn(ctx, items)
}
