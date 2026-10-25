package main

import "fmt"

func main() {
	type CustomerStatus int

	const (
		Regular CustomerStatus = iota
		Silver
		Gold
	)

	const (
		freeDeliveryThreshold = 3000.0
		deliveryPrice         = 299.0
	)

	var amount float64 = 4500
	status := Gold

	var discountRate float64
	var statusName string

	switch status {
	case Regular:
		statusName = "Regular"
	case Silver:
		statusName = "Silver"
		discountRate = 0.05
	case Gold:
		statusName = "Gold"
		discountRate = 0.10
	default:
		statusName = "неизвестный статус"
	}

	discount := amount * discountRate
	amountAfterDiscount := amount - discount

	delivery := deliveryPrice
	if amountAfterDiscount >= freeDeliveryThreshold {
		delivery = 0
	}

	total := amountAfterDiscount + delivery

	fmt.Printf("Сумма: %.2f, статус: %s\n", amount, statusName)
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f\n", total)
}
