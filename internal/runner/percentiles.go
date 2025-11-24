package runner

import (
	"sort"
	"time"
)

// Percentile calculates the percentile value from a sorted slice of durations
// percentile should be between 0 and 100
func Percentile(durations []time.Duration, percentile float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Calculate index for the percentile
	index := float64(len(sorted)-1) * percentile / 100.0
	lower := int(index)
	upper := lower + 1

	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}

	// Linear interpolation for more accurate percentile calculation
	weight := index - float64(lower)
	lowerValue := float64(sorted[lower])
	upperValue := float64(sorted[upper])

	return time.Duration(lowerValue + weight*(upperValue-lowerValue))
}

