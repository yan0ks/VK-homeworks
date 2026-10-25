package main

import "fmt"

func FirstUnique(s string) rune {
	counts := make(map[rune]int)
	for _, symbol := range s {
		counts[symbol]++
	}
	for _, symbol := range s {
		if counts[symbol] == 1 {
			return symbol
		}
	}
	return 0
}

func main() {
	fmt.Printf("%q\n", FirstUnique("abacabad"))
	fmt.Printf("%q\n", FirstUnique("abcabcd"))
	fmt.Printf("%q\n", FirstUnique("абракадабра"))
	fmt.Printf("%q\n", FirstUnique("aabbcc"))
}
