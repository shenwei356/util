package stats

import (
	"math"

	"github.com/twotwotwo/sorts"
)

type valueCount struct {
	Value, Count float64
}

type valueCounts []valueCount

func (c valueCounts) Len() int           { return len(c) }
func (c valueCounts) Less(i, j int) bool { return c[i].Value < c[j].Value }
func (c valueCounts) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type Quantiler struct {
	count map[float64]float64 // value -> count

	n             uint64  // n
	min, max, sum float64 // sum

	// for sorting
	counts    []valueCount // value, count
	accCounts []valueCount // value, accumulative count
	sorted    bool
}

// NewQuantiler initializes a Quantiler
func NewQuantiler() *Quantiler {
	return &Quantiler{count: make(map[float64]float64, 1024), min: math.MaxFloat64}
}

// Add adds a new element
func (stats *Quantiler) Add(value float64) {
	stats.n++
	stats.sum += value
	stats.count[value]++

	if value > stats.max {
		stats.max = value
	}
	if value < stats.min {
		stats.min = value
	}

	stats.sorted = false
}

func (stats *Quantiler) sort() {
	stats.counts = make([]valueCount, 0, len(stats.count))
	for value, count := range stats.count {
		stats.counts = append(stats.counts, valueCount{value, count})
	}
	sorts.Quicksort(valueCounts(stats.counts))

	stats.accCounts = make([]valueCount, len(stats.count))
	for i, data := range stats.counts {
		if i == 0 {
			stats.accCounts[i] = valueCount{data.Value, data.Count}
		} else {
			stats.accCounts[i] = valueCount{data.Value, data.Count + stats.accCounts[i-1].Count}
		}
	}

	stats.sorted = true
}

// Count returns number of elements
func (stats *Quantiler) Count() uint64 {
	return stats.n
}

// Min returns the minimum value
func (stats *Quantiler) Min() float64 {
	if stats.n == 0 {
		return 0
	}
	return stats.min
}

// Max returns the maxinimum length
func (stats *Quantiler) Max() float64 {
	return stats.max
}

// Sum returns the sum
func (stats *Quantiler) Sum() float64 {
	return stats.sum
}

// Mean returns mean
func (stats *Quantiler) Mean() float64 {
	return float64(stats.sum) / float64(stats.n)
}

// Q2 returns Q2
func (stats *Quantiler) Q2() float64 {
	return stats.Median()
}

// Median returns median
func (stats *Quantiler) Median() float64 {
	if !stats.sorted {
		stats.sort()
	}
	if len(stats.counts) == 0 {
		return 0
	}

	if len(stats.counts) == 1 {
		return float64(stats.counts[0].Value)
	}

	even := stats.n&1 == 0        // %2 == 0
	var iMedianL, iMedianR uint64 // 0-based
	if even {
		iMedianL = uint64(stats.n/2) - 1 // 3
		iMedianR = uint64(stats.n / 2)   // 4
	} else {
		iMedianL = uint64(stats.n / 2)
	}

	return stats.getValue(even, iMedianL, iMedianR)
}

// Q1 returns Q1
func (stats *Quantiler) Q1() float64 {
	if !stats.sorted {
		stats.sort()
	}
	if len(stats.counts) == 0 {
		return 0
	}

	if len(stats.counts) == 1 {
		return float64(stats.counts[0].Value) / 2
	}

	even := stats.n&1 == 0        // %2 == 0
	var iMedianL, iMedianR uint64 // 0-based
	var n uint64
	if even {
		n = stats.n / 2
	} else {
		n = (stats.n + 1) / 2
	}

	even = n%2 == 0
	if even {
		iMedianL = uint64(n/2) - 1
		iMedianR = uint64(n / 2)
	} else {
		iMedianL = uint64(n / 2)
	}

	return stats.getValue(even, iMedianL, iMedianR)
}

// Q3 returns Q3
func (stats *Quantiler) Q3() float64 {
	if !stats.sorted {
		stats.sort()
	}
	if len(stats.counts) == 0 {
		return 0
	}

	if len(stats.counts) == 1 {
		return float64(stats.counts[0].Value) / 2
	}

	even := stats.n&1 == 0        // %2 == 0
	var iMedianL, iMedianR uint64 // 0-based
	var mean, n uint64
	if even {
		n = stats.n / 2
		mean = n
	} else {
		n = (stats.n + 1) / 2
		mean = stats.n / 2
	}

	even = n%2 == 0
	if even {
		iMedianL = uint64(n/2) - 1 + mean
		iMedianR = uint64(n/2) + mean
	} else {
		iMedianL = uint64(n/2) + mean
	}

	return stats.getValue(even, iMedianL, iMedianR)
}

func (stats *Quantiler) getValue(even bool, iMedianL uint64, iMedianR uint64) float64 {

	var accCount float64

	var flag bool
	var prev float64

	for _, data := range stats.accCounts {
		accCount = data.Count

		if flag {
			// the middle two having different value.
			// example: 1, 2, 3, 4 or 1, 2
			return (data.Value + prev) / 2
		}

		if accCount >= float64(iMedianL+1) {
			if even {
				if accCount >= float64(iMedianR+1) {
					// having >=2 of same value in the middle.
					// example: 2, 2, 2, 3, 3, 4, 8, 8
					return data.Value
				}
				flag = true
				prev = data.Value
			} else {
				// right here
				return data.Value
			}
		}
	}

	// never happen
	// panic("should never happen")
	return 0
}

func (stats *Quantiler) Percentile(percent float64) float64 {
	if percent <= 0 || percent > 100 {
		panic("invalid percentile")
	}
	if !stats.sorted {
		stats.sort()
	}
	if len(stats.counts) == 0 {
		return 0
	}

	if len(stats.counts) == 1 {
		return float64(stats.counts[0].Value)
	}

	i0 := float64(stats.n) * percent / 100
	i := math.Floor(i0)

	even := math.Abs(i0-i) > 0.001
	var iMedianL, iMedianR uint64 // 0-based
	if even {
		iMedianL = uint64(i) - 1
		iMedianR = uint64(i)
	} else {
		iMedianL = uint64(i - 1)
	}

	return stats.getValue(even, iMedianL, iMedianR)
}
