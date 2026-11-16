package session2

import "testing"

func TestMyHasDuplicate(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []int
		want bool
	}{
		{name: "negative duplicate", in: []int{-3, 2, -3}, want: true},
		{name: "duplicate zero", in: []int{0, 1, 0}, want: true},
		{name: "distinct negative values", in: []int{-3, -2, -1}},
		{name: "single zero", in: []int{0}},
		{name: "empty slice", in: []int{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDuplicate(tt.in); got != tt.want {
				t.Errorf("HasDuplicate(%v) = %t, want %t", tt.in, got, tt.want)
			}
		})
	}
}

func TestMyTwoItems(t *testing.T) {
	for _, tt := range []struct {
		name   string
		prices []int
		k      int
		wantOK bool
	}{
		{name: "position zero", prices: []int{4, 9}, k: 13, wantOK: true},
		{name: "two zero prices", prices: []int{0, 0}, wantOK: true},
		{name: "negative price", prices: []int{-3, 8, 2}, k: 5, wantOK: true},
		{name: "several valid pairs", prices: []int{1, 4, 2, 3}, k: 5, wantOK: true},
		{name: "no matching pair", prices: []int{1, 2, 3}, k: 20},
		{name: "cannot reuse a single item", prices: []int{5}, k: 10},
		{name: "empty input", k: 5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			i, j, ok := TwoItems(tt.prices, tt.k)
			if ok != tt.wantOK {
				t.Fatalf("TwoItems(%v, %d) = (%d, %d, %t), want found = %t", tt.prices, tt.k, i, j, ok, tt.wantOK)
			}
			if !ok {
				if i != 0 || j != 0 {
					t.Errorf("missing pair indices = (%d, %d), want (0, 0)", i, j)
				}
				return
			}
			if i < 0 || i >= j || j >= len(tt.prices) {
				t.Fatalf("invalid pair indices (%d, %d) for %v", i, j, tt.prices)
			}
			if sum := tt.prices[i] + tt.prices[j]; sum != tt.k {
				t.Errorf("pair sum = %d, want %d", sum, tt.k)
			}
		})
	}
}

func TestMyBracketsOK(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		want bool
	}{
		{name: "text without brackets", in: "Привет, мир!", want: true},
		{name: "text inside brackets", in: "а{б[в(г)]д}е", want: true},
		{name: "crossed pairs", in: "([)]"},
		{name: "closing bracket first", in: "текст]"},
		{name: "opening bracket left over", in: "([]"},
		{name: "empty input", want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := BracketsOK(tt.in); got != tt.want {
				t.Errorf("BracketsOK(%q) = %t, want %t", tt.in, got, tt.want)
			}
		})
	}
}

func TestMyIsPalindrome(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		want bool
	}{
		{name: "cyrillic palindrome", in: "казак", want: true},
		{name: "one rune", in: "ё", want: true},
		{name: "different case", in: "Аа"},
		{name: "spaces are included", in: "а а", want: true},
		{name: "asymmetric space", in: "а а "},
		{name: "punctuation is included", in: "а!"},
		{name: "symmetric punctuation", in: "!а!", want: true},
		{name: "empty input", want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.in); got != tt.want {
				t.Errorf("IsPalindrome(%q) = %t, want %t", tt.in, got, tt.want)
			}
		})
	}
}

func TestMyIsAnagram(t *testing.T) {
	for _, fn := range []struct {
		name string
		run  func(string, string) bool
	}{
		{name: "sort", run: IsAnagramSort},
		{name: "map", run: IsAnagramMap},
	} {
		t.Run(fn.name, func(t *testing.T) {
			for _, tt := range []struct {
				name string
				a, b string
				want bool
			}{
				{name: "cyrillic", a: "липа", b: "пила", want: true},
				{name: "different frequencies", a: "ааё", b: "аёё"},
				{name: "different case", a: "Кот", b: "ток"},
				{name: "space is missing", a: "а б", b: "ба"},
				{name: "spaces and punctuation", a: "а! б", b: "б а!", want: true},
				{name: "different rune lengths", a: "я", b: "aa"},
				{name: "both empty", want: true},
				{name: "first empty", b: "а"},
				{name: "second empty", a: "а"},
			} {
				t.Run(tt.name, func(t *testing.T) {
					if got := fn.run(tt.a, tt.b); got != tt.want {
						t.Errorf("%s(%q, %q) = %t, want %t", fn.name, tt.a, tt.b, got, tt.want)
					}
				})
			}
		})
	}
}

func TestMyCanFulfill(t *testing.T) {
	for _, tt := range []struct {
		name         string
		stock, order []string
		want         bool
	}{
		{name: "both empty", want: true},
		{name: "repeated product", stock: []string{"чай", "чай"}, order: []string{"чай", "чай"}, want: true},
		{name: "not enough copies", stock: []string{"чай", "кофе"}, order: []string{"чай", "чай"}},
		{name: "missing product", stock: []string{"чай"}, order: []string{"сахар"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanFulfill(tt.stock, tt.order); got != tt.want {
				t.Errorf("CanFulfill(%v, %v) = %t, want %t", tt.stock, tt.order, got, tt.want)
			}
		})
	}
}

func TestMyMinStack(t *testing.T) {
	for _, tt := range []struct {
		name    string
		values  []int
		wantPop []int
		wantMin []int
	}{
		{name: "repeated minimum", values: []int{5, 3, 3, 7}, wantPop: []int{7, 3, 3, 5}, wantMin: []int{3, 3, 3, 5}},
		{name: "negative minimum", values: []int{0, -5, -2}, wantPop: []int{-2, -5, 0}, wantMin: []int{-5, -5, 0}},
		{name: "single value", values: []int{8}, wantPop: []int{8}, wantMin: []int{8}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stack MinStack
			for _, value := range tt.values {
				stack.Push(value)
			}
			for index, want := range tt.wantPop {
				if got := stack.Minimum(); got != tt.wantMin[index] {
					t.Errorf("Minimum() before pop %d = %d, want %d", index, got, tt.wantMin[index])
				}
				if got := stack.Pop(); got != want {
					t.Fatalf("Pop() = %d, want %d", got, want)
				}
			}
			stack.Push(9)
			if got := stack.Minimum(); got != 9 {
				t.Errorf("Minimum() after refilling = %d, want 9", got)
			}
			if got := stack.Pop(); got != 9 {
				t.Errorf("Pop() after refilling = %d, want 9", got)
			}
		})
	}
}
