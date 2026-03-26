package app_test

import (
	"testing"

	"github.com/hferr/pack-api/internal/app"
)

func TestCalculateMinPackOrder(t *testing.T) {
	var packSizes = app.Packs{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}

	var testCases = map[string]struct {
		packs    app.Packs
		items    int
		expected map[int]int
	}{
		"no packs available": {
			packs:    nil,
			items:    10,
			expected: map[int]int{},
		},
		"250 items ordered": {
			packs: packSizes,
			items: 250,
			expected: map[int]int{
				250: 1,
			},
		},
		"251 items ordered": {
			packs: packSizes,
			items: 251,
			expected: map[int]int{
				500: 1,
			},
		},
		"501 items ordered": {
			packs: packSizes,
			items: 501,
			expected: map[int]int{
				500: 1,
				250: 1,
			},
		},
		"12001 items ordered": {
			packs: packSizes,
			items: 12001,
			expected: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
		"edge case pack test": {
			packs: app.Packs{
				{Size: 23},
				{Size: 31},
				{Size: 53},
			},
			items: 500000,
			expected: map[int]int{
				23: 2,
				31: 7,
				53: 9429,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svs := app.NewPackService()
			got := svs.CalculateMinPackOrder(tc.packs, tc.items)

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
