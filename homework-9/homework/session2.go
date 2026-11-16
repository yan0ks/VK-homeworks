package session2

import "slices"

func HasDuplicate(a []int) bool {
	seen := make(map[int]bool)
	for _, value := range a {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}

func TwoItems(prices []int, k int) (int, int, bool) {
	seen := make(map[int]int)
	for index, price := range prices {
		need := k - price
		if previous, ok := seen[need]; ok {
			return previous, index, true
		}
		seen[price] = index
	}
	return 0, 0, false
}

func BracketsOK(s string) bool {
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	stack := []rune{}
	for _, symbol := range s {
		switch symbol {
		case '(', '[', '{':
			stack = append(stack, symbol)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[symbol] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func IsPalindrome(s string) bool {
	symbols := []rune(s)
	lo := 0
	hi := len(symbols) - 1
	for lo < hi {
		if symbols[lo] != symbols[hi] {
			return false
		}
		lo++
		hi--
	}
	return true
}

func IsAnagramSort(a, b string) bool {
	left := []rune(a)
	right := []rune(b)
	if len(left) != len(right) {
		return false
	}
	slices.Sort(left)
	slices.Sort(right)
	return slices.Equal(left, right)
}

func IsAnagramMap(a, b string) bool {
	counts := make(map[rune]int)
	for _, symbol := range a {
		counts[symbol]++
	}
	for _, symbol := range b {
		if counts[symbol] == 0 {
			return false
		}
		counts[symbol]--
	}
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}

func CanFulfill(stock, order []string) bool {
	left := make(map[string]int)
	for _, product := range stock {
		left[product]++
	}
	for _, product := range order {
		if left[product] == 0 {
			return false
		}
		left[product]--
	}
	return true
}

type MinStack struct {
	values   []int
	minimums []int
}

func (s *MinStack) Push(x int) {
	minimum := x
	if len(s.minimums) > 0 && s.minimums[len(s.minimums)-1] < minimum {
		minimum = s.minimums[len(s.minimums)-1]
	}
	s.values = append(s.values, x)
	s.minimums = append(s.minimums, minimum)
}

func (s *MinStack) Pop() int {
	index := len(s.values) - 1
	value := s.values[index]
	s.values = s.values[:index]
	s.minimums = s.minimums[:index]
	return value
}

func (s *MinStack) Minimum() int {
	return s.minimums[len(s.minimums)-1]
}
