package stats

import (
	"math"
	"math/rand"
	"testing"

	stats2 "github.com/montanaflynn/stats"
)

type testCase struct {
	data []float64

	q1, median, q3 float64
	min, max       float64
}

var cases = []testCase{
	{
		data:   []float64{},
		median: 0,
		q1:     0,
		q3:     0,
		min:    0,
		max:    0,
	},
	{
		data:   []float64{2},
		median: 2,
		q1:     1,
		q3:     1,
		min:    2,
		max:    2,
	},
	{
		data:   []float64{1, 2},
		median: 1.5,
		q1:     1,
		q3:     2,
		min:    1,
		max:    2,
	},
	{
		data:   []float64{1, 2, 3},
		median: 2,
		q1:     1.5,
		q3:     2.5,
		min:    1,
		max:    3,
	},
	{
		data:   []float64{1, 2, 3, 4},
		median: 2.5,
		q1:     1.5,
		q3:     3.5,
		min:    1,
		max:    4,
	},
	{
		data:   []float64{2, 3, 4, 5, 6, 7, 8, 9},
		median: 5.5,
		q1:     3.5,
		q3:     7.5,
		min:    2,
		max:    9,
	},
	{
		data:   []float64{0.5, 0.6, 0.7, 0.8, 0.8, 0.85, 0.9},
		median: 0.8,
		q1:     0.65,
		q3:     0.825,
		min:    0.5,
		max:    0.9,
	},
	{
		data:   []float64{1, 0.8, 0.8, 0.85, 0.9},
		median: 0.85,
		q1:     0.8,
		q3:     0.9,
		min:    0.8,
		max:    1,
	},
}

func Test(t *testing.T) {
	for i, _case := range cases {
		rand.Shuffle(len(_case.data), func(i, j int) {
			_case.data[i], _case.data[j] = _case.data[j], _case.data[i]
		})

		stats := NewQuantiler()
		for _, l := range _case.data {
			stats.Add(l)
		}
		if stats.Count() != uint64(len(_case.data)) {
			t.Errorf("case %d: count mismatch", i)
		}

		min := stats.Min()
		if min != _case.min {
			t.Errorf("case %d: min mismatch: %f != %f", i, min, _case.min)
		}

		max := stats.Max()
		if max != _case.max {
			t.Errorf("case %d: max mismatch: %f != %f", i, max, _case.max)
		}

		median := stats.Median()
		if math.Abs(median-_case.median) > 0.001 {
			t.Errorf("case %d: median mismatch: %f != %f", i, median, _case.median)
		}

		q1 := stats.Q1()
		if math.Abs(q1-_case.q1) > 0.001 {
			t.Errorf("case %d: q1 mismatch: %f != %f", i, q1, _case.q1)
		}

		q3 := stats.Q3()
		if math.Abs(q3-_case.q3) > 0.001 {
			t.Errorf("case %d: q3 mismatch: %f != %f", i, q3, _case.q3)
		}

	}
}

var cases2 = []testCase{
	{
		data: []float64{},
	},
	{
		data: []float64{0.8},
	},
	{
		data: []float64{1, 2},
	},
	{
		data: []float64{1, 2, 3},
	},
	{
		data: []float64{1, 2, 3, 3, 4, 5, 6, 7, 5},
	},
	{
		data: []float64{0, 1, 1, 2, 3, 3, 3, 4, 5, 6, 7, 4, 2, 1, 4, 5, 6, 6, 4, 2, 2, 4, 10},
	},
}

func Test2(t *testing.T) {
	for i, _case := range cases2 {
		rand.Shuffle(len(_case.data), func(i, j int) {
			_case.data[i], _case.data[j] = _case.data[j], _case.data[i]
		})

		stats := NewQuantiler()
		for _, l := range _case.data {
			stats.Add(l)
		}

		p90 := stats.Percentile(90)
		pp90, _ := stats2.Percentile(_case.data, 90)
		if math.Abs(p90-pp90) > 0.001 {
			t.Errorf("case %d: p90 mismatch: %f != %f", i, p90, pp90)
		}
	}
}
