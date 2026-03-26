package app

import "math"

type PackService interface {
	CalculateMinPackOrder(packs Packs, items int) map[int]int
}

type packService struct{}

func NewPackService() PackService {
	return &packService{}
}

// CalculateMinPack calculates the optimal combination of pack sizes needed to fulfill an order according
// to the following constraints:
// 1. Only whole packs can be used
// 2. Least amount of items are sent (Priority)
// 3. Least amount of packs are sent
//
// Returns
// - map[int]int: map of pack size and quantity
func (s *packService) CalculateMinPackOrder(packs Packs, items int) map[int]int {
	return minPacks(packs, items)
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
