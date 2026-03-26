package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hferr/pack-api/internal/app"
	"github.com/hferr/pack-api/internal/test/mock"
)

func TestListPacks(t *testing.T) {
	var testCases = map[string]struct {
		wantErr  bool
		expected *int
		repo     *mock.Repo
	}{
		"2 packs": {
			wantErr:  false,
			expected: new(2),
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return app.Packs{
						{Size: 50},
						{Size: 6},
					}, nil
				},
			},
		},
		"error listing packs": {
			wantErr:  true,
			expected: nil,
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return app.Packs{}, fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svs := app.NewPackService(tc.repo)
			got, err := svs.ListPacks(context.Background())

			if !tc.wantErr {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if len(got) != *tc.expected {
					t.Errorf("expected %d pack sizes, got %d", *tc.expected, len(got))
				}
			}

			if err == nil && tc.wantErr {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestCreatePack(t *testing.T) {
	var testCases = map[string]struct {
		wantErr bool
		size    int
		repo    *mock.Repo
	}{
		"success": {
			wantErr: false,
			size:    50,
			repo: &mock.Repo{
				CreatePackFn: func(ctx context.Context, pack *app.Pack) error {
					return nil
				},
			},
		},
		"validation error creating pack": {
			wantErr: true,
			size:    -1,
		},
		"error creating pack": {
			wantErr: true,
			size:    30,
			repo: &mock.Repo{
				CreatePackFn: func(ctx context.Context, pack *app.Pack) error {
					return fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svs := app.NewPackService(tc.repo)
			got, err := svs.CreatePack(context.Background(), tc.size)

			if !tc.wantErr {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if got.Size != tc.size {
					t.Errorf("expected pack to have size %d, got: %d", tc.size, got.Size)
				}
			}

			if err == nil && tc.wantErr {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestCalculateMinPackOrder(t *testing.T) {
	var packSizes = app.Packs{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}

	var testCases = map[string]struct {
		items    int
		expected map[int]int
		repo     *mock.Repo
	}{
		"no packs available": {
			items:    10,
			expected: map[int]int{},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return nil, nil
				},
			},
		},
		"250 items ordered": {
			items: 250,
			expected: map[int]int{
				250: 1,
			},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return packSizes, nil
				},
			},
		},
		"251 items ordered": {
			items: 251,
			expected: map[int]int{
				500: 1,
			},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return packSizes, nil
				},
			},
		},
		"501 items ordered": {
			items: 501,
			expected: map[int]int{
				500: 1,
				250: 1,
			},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return packSizes, nil
				},
			},
		},
		"12001 items ordered": {
			items: 12001,
			expected: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return packSizes, nil
				},
			},
		},
		"edge case pack test": {
			items: 500000,
			expected: map[int]int{
				23: 2,
				31: 7,
				53: 9429,
			},
			repo: &mock.Repo{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return app.Packs{
						{Size: 23},
						{Size: 31},
						{Size: 53},
					}, nil
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svs := app.NewPackService(tc.repo)
			got, err := svs.CalculateMinPackOrder(context.Background(), tc.items)

			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if len(got) != len(tc.expected) {
				t.Errorf("expected map len %d, got %d", len(tc.expected), len(got))
			}

			for k, v := range tc.expected {
				if got[k] != v {
					t.Errorf("for key %d: expected %d, got %d", k, v, got[k])
				}
			}
		})
	}
}
