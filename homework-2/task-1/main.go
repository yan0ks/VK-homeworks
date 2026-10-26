package main

import (
	"errors"
	"fmt"
)

type Operation func(a, b float64) (float64, error)

var operations = map[string]Operation{
	"+": func(a, b float64) (float64, error) { return a + b, nil },
	"-": func(a, b float64) (float64, error) { return a - b, nil },
	"*": func(a, b float64) (float64, error) { return a * b, nil },
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	},
}

func Calculate(a float64, op string, b float64) (float64, error) {
	operation, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("неизвестная операция %q", op)
	}
	return operation(a, b)
}

func printCalculation(a float64, op string, b float64) {
	result, err := Calculate(a, op, b)
	if err != nil {
		fmt.Printf("%g %s %g: %v\n", a, op, b, err)
		return
	}
	fmt.Printf("%g %s %g = %.2f\n", a, op, b, result)
}

func main() {
	printCalculation(10, "+", 5)
	printCalculation(10, "/", 4)
	printCalculation(7, "*", 6)
	printCalculation(1, "/", 0)
	printCalculation(2, "%", 3)
}
