package collection

import (
	"slices"
	"strconv"
	"testing"
)

func TestFilter(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "keeps order and duplicates", input: []int{4, 1, 2, 4}, want: []int{4, 2, 4}},
		{name: "no matches", input: []int{1, 3}},
		{name: "empty input", input: []int{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, func(value int) bool { return value%2 == 0 })
			if !slices.Equal(got, tt.want) {
				t.Errorf("Filter(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  []string
	}{
		{name: "converts each element in order", input: []int{7, 0, -2}, want: []string{"7", "0", "-2"}},
		{name: "empty input", input: []int{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.input, strconv.Itoa)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Map(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  int
		found bool
	}{
		{name: "returns first match", input: []int{1, 4, 2}, want: 4, found: true},
		{name: "matches zero", input: []int{0, 2}, found: true},
		{name: "no matches", input: []int{1, 3}},
		{name: "empty input", input: []int{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, found := Find(tt.input, func(value int) bool { return value%2 == 0 })
			if got != tt.want || found != tt.found {
				t.Errorf("Find(%v) = (%d, %t), want (%d, %t)", tt.input, got, found, tt.want, tt.found)
			}
		})
	}
}
