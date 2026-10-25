package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(s) {
		counts[word]++
	}
	return counts
}

func main() {
	fmt.Println(WordCount("go is fun go"))
	fmt.Println(WordCount("Hello hello World"))
	fmt.Println(WordCount(" мир мир труд "))
}
