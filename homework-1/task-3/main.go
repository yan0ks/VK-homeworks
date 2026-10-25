package main

import "fmt"

func main() {
	const (
		initialDeposit = 100000.0
		annualRate     = 12.0
		monthsInYear   = 12
		depositMonths  = 12
	)

	monthlyRate := annualRate / monthsInYear / 100
	balance := initialDeposit

	for month := 1; month <= depositMonths; month++ {
		balance += balance * monthlyRate
		fmt.Printf("месяц %2d: %.2f\n", month, balance)
	}

	balance = initialDeposit
	months := 0
	for balance < initialDeposit*2 {
		balance += balance * monthlyRate
		months++
	}

	fmt.Printf("Вклад удвоится через %d мес. Баланс: %.2f\n", months, balance)
}

/*
Вопрос: почему в реальных системах деньги не хранят во float64?

Ответ: float64 хранит многие десятичные дроби приблизительно, поэтому
при расчётах могут накапливаться ошибки округления.
*/
