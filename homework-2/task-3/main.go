package main

import (
	"fmt"
	"time"
)

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n <= 0 {
		return
	}
	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
}

func printMoves(n int) {
	number := 0
	Hanoi(n, "A", "C", "B", func(from, to string) {
		number++
		fmt.Printf("%d: %s -> %s\n", number, from, to)
	})
}

func countMoves(n int) {
	started := time.Now()
	defer func() {
		fmt.Printf("время: %s\n", time.Since(started))
	}()
	count := 0
	Hanoi(n, "A", "C", "B", func(from, to string) { count++ })
	fmt.Printf("n=%d: %d ходов\n", n, count)
}

func main() {
	printMoves(3)
	countMoves(20)
}

/*
Вопрос: почему ходов ровно 2ⁿ − 1?

Ответ: сначала переносим n−1 дисков на вспомогательный стержень,
затем самый большой диск на целевой, затем n−1 дисков на него.
Поэтому T(n) = 2*T(n−1) + 1, а T(0) = 0.
Раскрывая рекурсию, получаем T(n) = 2ⁿ − 1.
*/
