package main

import "fmt"

func main() {
	const N = 8

	for row := 0; row < N; row++ {
		for column := 0; column < N; column++ {
			if (row+column)%2 == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
