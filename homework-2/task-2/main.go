package main

import (
	"errors"
	"fmt"
)

func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	current := initial
	deposit = func(amount int64) { current += amount }
	withdraw = func(amount int64) error {
		if amount > current {
			return errors.New("недостаточно средств")
		}
		current -= amount
		return nil
	}
	balance = func() int64 { return current }
	return
}

func rubles(kopecks int64) string {
	return fmt.Sprintf("%d.%02d", kopecks/100, kopecks%100)
}

func accountExample() {
	deposit, withdraw, balance := NewAccount(100000)
	_, _, balanceB := NewAccount(0)
	fmt.Printf("Счёт A: %s\n", rubles(balance()))
	deposit(50000)
	fmt.Printf("Пополнение 500.00 -> %s\n", rubles(balance()))
	if err := withdraw(200000); err != nil {
		fmt.Printf("Снятие 2000.00: %v (баланс %s)\n", err, rubles(balance()))
	}
	if err := withdraw(30050); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Снятие 300.50 -> %s\n", rubles(balance()))
	fmt.Printf("Счёт B: %s\n", rubles(balanceB()))
}

func main() {
	accountExample()
}
