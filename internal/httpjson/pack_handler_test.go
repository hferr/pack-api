package httpjson_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hferr/pack-api/internal/app"
	"github.com/hferr/pack-api/internal/httpjson"
	"github.com/hferr/pack-api/internal/test"
	"github.com/hferr/pack-api/internal/test/mock"
)

func TestListPackSizes(t *testing.T) {
	var testCases = map[string]struct {
		wantCode int
		svs      *mock.PackService
	}{
		"sucess": {
			wantCode: 200,
			svs: &mock.PackService{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					var packs app.Packs
					for i := range 2 {
						packs = append(packs, *app.NewPack(i))
					}
					return packs, nil
				},
			},
		},
		"no packs": {
			wantCode: 200,
			svs: &mock.PackService{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return nil, nil
				},
			},
		},
		"internal error": {
			wantCode: 500,
			svs: &mock.PackService{
				ListPacksFn: func(ctx context.Context) (app.Packs, error) {
					return nil, fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := httpjson.NewHandler(tc.svs)
			resp := test.DoHttpRequest(
				handler,
				http.MethodGet,
				"/packs/sizes",
				nil,
			)

			gotCode := resp.StatusCode

			if tc.wantCode != gotCode {
				t.Errorf("expected status code %d, got: %d", tc.wantCode, gotCode)
			}
		})
	}
}

func TestCreatePack(t *testing.T) {
	var testCases = map[string]struct {
		wantCode int
		packSize int
		svs      *mock.PackService
	}{
		"sucess": {
			wantCode: 201,
			packSize: 300,
			svs: &mock.PackService{
				CreatePackFn: func(ctx context.Context, size int) (*app.Pack, error) {
					return &app.Pack{
						ID:   uuid.New(),
						Size: size,
					}, nil
				},
			},
		},
		"validation error throws bad request": {
			wantCode: 400,
			packSize: -1,
			svs: &mock.PackService{
				CreatePackFn: func(ctx context.Context, size int) (*app.Pack, error) {
					return nil, &app.ValidationError{}
				},
			},
		},
		"internal error": {
			wantCode: 500,
			packSize: 300,
			svs: &mock.PackService{
				CreatePackFn: func(ctx context.Context, size int) (*app.Pack, error) {
					return nil, fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := httpjson.NewHandler(tc.svs)

			req := httpjson.CreatePackRequest{
				Size: tc.packSize,
			}
			reqJson, err := json.Marshal(req)
			if err != nil {
				t.Error(err)
			}

			resp := test.DoHttpRequest(
				handler,
				http.MethodPost,
				"/packs",
				bytes.NewReader(reqJson),
			)

			gotCode := resp.StatusCode

			if tc.wantCode != gotCode {
				t.Errorf("expected status code %d, got: %d", tc.wantCode, gotCode)
			}
		})
	}
}

func TestRebuildPacks(t *testing.T) {
	var testCases = map[string]struct {
		wantCode int
		sizes    []int
		svs      *mock.PackService
	}{
		"sucess": {
			wantCode: 200,
			sizes:    []int{200, 300, 400},
			svs: &mock.PackService{
				RebuildPacksFn: func(ctx context.Context, sizes []int) (app.Packs, error) {
					var packs app.Packs
					for i := range sizes {
						packs = append(packs, *app.NewPack(i))
					}
					return packs, nil
				},
			},
		},
		"validation error throws bad request": {
			wantCode: 400,
			sizes:    []int{200, 300, 400, 0},
			svs: &mock.PackService{
				RebuildPacksFn: func(ctx context.Context, sizes []int) (app.Packs, error) {
					return nil, &app.ValidationError{}
				},
			},
		},
		"internal error": {
			wantCode: 500,
			sizes:    []int{200, 300, 400},
			svs: &mock.PackService{
				RebuildPacksFn: func(ctx context.Context, sizes []int) (app.Packs, error) {
					return nil, fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := httpjson.NewHandler(tc.svs)

			req := httpjson.RebuildPacksRequest{
				Sizes: tc.sizes,
			}
			reqJson, err := json.Marshal(req)
			if err != nil {
				t.Error(err)
			}

			resp := test.DoHttpRequest(
				handler,
				http.MethodPost,
				"/packs/rebuild",
				bytes.NewReader(reqJson),
			)

			gotCode := resp.StatusCode

			if tc.wantCode != gotCode {
				t.Errorf("expected status code %d, got: %d", tc.wantCode, gotCode)
			}
		})
	}
}

func TestCalculateMinPackOrder(t *testing.T) {
	var testCases = map[string]struct {
		wantCode int
		items    int
		svs      *mock.PackService
	}{
		"sucess": {
			wantCode: 200,
			items:    300,
			svs: &mock.PackService{
				CalculateMinPackOrderFn: func(ctx context.Context, items int) (map[int]int, error) {
					return map[int]int{
						500: 1,
					}, nil
				},
			},
		},
		"internal error": {
			wantCode: 500,
			items:    300,
			svs: &mock.PackService{
				CalculateMinPackOrderFn: func(ctx context.Context, items int) (map[int]int, error) {
					return nil, fmt.Errorf("error")
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := httpjson.NewHandler(tc.svs)

			req := httpjson.CalculatePacksForItemsRequest{
				Items: tc.items,
			}
			reqJson, err := json.Marshal(req)
			if err != nil {
				t.Error(err)
			}

			resp := test.DoHttpRequest(
				handler,
				http.MethodPost,
				"/packs/calculate-order",
				bytes.NewReader(reqJson),
			)

			gotCode := resp.StatusCode

			if tc.wantCode != gotCode {
				t.Errorf("expected status code %d, got: %d", tc.wantCode, gotCode)
			}
		})
	}
}
