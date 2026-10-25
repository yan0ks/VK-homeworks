package main

import (
	"fmt"
	"strings"
)

func AreAnagrams(s1, s2 string) bool {
	counts := make(map[rune]int)
	for _, letter := range strings.ToLower(s1) {
		if letter != ' ' {
			counts[letter]++
		}
	}
	for _, letter := range strings.ToLower(s2) {
		if letter != ' ' {
			counts[letter]--
		}
	}
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(AreAnagrams("listen", "silent"))
	fmt.Println(AreAnagrams("Eleven plus two", "Twelve plus one"))
	fmt.Println(AreAnagrams("Апельсин", "спаниель"))
	fmt.Println(AreAnagrams("Hello", "World"))
}
