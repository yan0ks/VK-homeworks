package main

import (
	"fmt"
	"unicode"
)

func IsPalindrome(s string) bool {
	letters := make([]rune, 0)
	for _, symbol := range s {
		if unicode.IsLetter(symbol) {
			letters = append(letters, unicode.ToLower(symbol))
		}
	}
	for left, right := 0, len(letters)-1; left < right; left, right = left+1, right-1 {
		if letters[left] != letters[right] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(IsPalindrome("шалаш"))
	fmt.Println(IsPalindrome("А роза упала на лапу Азора"))
	fmt.Println(IsPalindrome("Was it a car or a cat I saw?"))
	fmt.Println(IsPalindrome("hello"))
}
