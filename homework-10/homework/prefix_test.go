package session3

import "testing"

func TestSumRange(t *testing.T) {
	cases := []struct {
		name        string
		nums        []int
		left, right int
		want        int
	}{
		{"с начала", []int{-2, 0, 3, -5, 2, -1}, 0, 2, 1},
		{"до конца", []int{-2, 0, 3, -5, 2, -1}, 2, 5, -1},
		{"весь массив", []int{-2, 0, 3, -5, 2, -1}, 0, 5, -3},
		{"как на доске", []int{3, 1, 4, 1, 5, 9}, 1, 3, 6},
		{"один элемент", []int{7}, 0, 0, 7},
		{"left == right в середине", []int{3, 1, 4}, 1, 1, 1},
	}
	for _, c := range cases {
		if got := SumRangeNaive(c.nums, c.left, c.right); got != c.want {
			t.Errorf("%s: SumRangeNaive(%d, %d) = %d, ждали %d", c.name, c.left, c.right, got, c.want)
		}
		na := Constructor(c.nums)
		if got := na.SumRange(c.left, c.right); got != c.want {
			t.Errorf("%s: SumRange(%d, %d) = %d, ждали %d", c.name, c.left, c.right, got, c.want)
		}
	}
}
