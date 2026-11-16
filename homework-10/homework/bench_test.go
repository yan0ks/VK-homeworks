package session3

import (
	"math/rand"
	"testing"
)

const benchN, benchQ = 100_000, 1_000

var benchmarkTotal int

func benchData() ([]int, [][2]int) {
	r := rand.New(rand.NewSource(1))
	nums := make([]int, benchN)
	for i := range nums {
		nums[i] = r.Intn(1000)
	}
	queries := make([][2]int, benchQ)
	for i := range queries {
		l, rr := r.Intn(benchN), r.Intn(benchN)
		queries[i] = [2]int{min(l, rr), max(l, rr)}
	}
	return nums, queries
}

func BenchmarkSumRangeNaive(b *testing.B) {
	nums, queries := benchData()
	var total int
	b.ResetTimer()
	for range b.N {
		for _, q := range queries {
			total += SumRangeNaive(nums, q[0], q[1])
		}
	}
	benchmarkTotal = total
}

func BenchmarkSumRangePrefix(b *testing.B) {
	nums, queries := benchData()
	var total int
	b.ResetTimer()
	for range b.N {
		array := Constructor(nums)
		for _, q := range queries {
			total += array.SumRange(q[0], q[1])
		}
	}
	benchmarkTotal = total
}
