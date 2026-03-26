package app

import (
	"context"
	"fmt"
	"math"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for field: %s, message: %s", e.Field, e.Message)
}

type PackService interface {
	ListPacks(ctx context.Context) (Packs, error)
	CreatePack(ctx context.Context, size int) (*Pack, error)
	RebuildPacks(ctx context.Context, sizes []int) (Packs, error)
	CalculateMinPackOrder(ctx context.Context, items int) (map[int]int, error)
}

type packService struct {
	repo Repo
}

func NewPackService(r Repo) PackService {
	return &packService{
		repo: r,
	}
}

func (s *packService) ListPacks(ctx context.Context) (Packs, error) {
	return s.repo.ListPacks(ctx)
}

func (s *packService) CreatePack(ctx context.Context, size int) (*Pack, error) {
	pack := NewPack(size)
	if err := pack.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreatePack(ctx, pack); err != nil {
		return nil, fmt.Errorf("Error while creating pack: %v", err)
	}

	return pack, nil
}

func (s *packService) RebuildPacks(ctx context.Context, sizes []int) (Packs, error) {
	// For simplicity, RebuildPackages takes in an array of pack sizes, deletes all existing
	// entries in the packs table and rebuilds it with the new values
	packs := make(Packs, len(sizes))
	for i, size := range sizes {
		pack := NewPack(size)
		if err := pack.Validate(); err != nil {
			return nil, err
		}

		packs[i] = *pack
	}

	if err := s.repo.RebuildPacks(ctx, packs); err != nil {
		return nil, fmt.Errorf("Error occurred while rebuilding packs: %v", err)
	}

	return s.ListPacks(ctx)
}

// CalculateMinPack calculates the optimal combination of pack sizes needed to fulfill an order according
// to the following constraints:
// 1. Only whole packs can be used
// 2. Least amount of items are sent (Priority)
// 3. Least amount of packs are sent
//
// Returns
// - map[int]int: map of pack size and quantity
func (s *packService) CalculateMinPackOrder(ctx context.Context, items int) (map[int]int, error) {
	packs, err := s.repo.ListPacks(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error occurred when listing packs: %v", err)
	}

	return minPacks(packs, items), nil
}

type State struct {
	totalItems int
	packs      int
	prev       int
	lastPack   int
}

func minPacks(packs Packs, items int) map[int]int {
	// Don't calculate if there's no packs
	if len(packs) == 0 {
		return map[int]int{}
	}

	packSizes := packs.GetSortedSizes()

	// Early return, if the items ordered are less than the smallest size then that
	// pack is enough
	if items < packSizes[0] {
		return map[int]int{
			packSizes[0]: 1,
		}
	}

	maxPack := packSizes[len(packSizes)-1]
	// Anything beyond items + maxPack is worse
	limit := items + maxPack

	dp := make([]State, limit+1)

	// Initialize all states with an "unreachable" number
	for i := range dp {
		dp[i] = State{totalItems: math.MaxInt32}
	}

	dp[0] = State{
		totalItems: 0,
		packs:      0,
		prev:       -1, // no previous state
	}

	for i := 0; i <= limit; i++ {
		if dp[i].totalItems == math.MaxInt32 {
			continue
		}

		for _, pack := range packSizes {
			next := i + pack
			if next > limit {
				continue
			}

			newItems := dp[i].totalItems + pack
			newPacks := dp[i].packs + 1

			// Prioritize according to the given constraints:
			// 1 - minimize total items
			// 2 - minimize number of packs
			if newItems < dp[next].totalItems ||
				(newItems == dp[next].totalItems && newPacks < dp[next].packs) {

				dp[next] = State{
					totalItems: newItems,
					packs:      newPacks,
					prev:       i,
					lastPack:   pack,
				}
			}
		}
	}

	// Find best solution and reconstruct the combination

	bestIdx := -1
	bestItems := math.MaxInt32
	bestPacks := math.MaxInt32

	for i := items; i <= limit; i++ {
		if dp[i].totalItems < bestItems ||
			(dp[i].totalItems == bestItems && dp[i].packs < bestPacks) {
			bestItems = dp[i].totalItems
			bestPacks = dp[i].packs
			bestIdx = i
		}
	}

	result := make(map[int]int)
	for bestIdx > 0 {
		pack := dp[bestIdx].lastPack
		bestIdx = dp[bestIdx].prev
		result[pack]++
	}

	return result
}
