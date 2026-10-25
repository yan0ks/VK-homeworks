package main

import "fmt"

func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, fmt.Errorf("индекс %d вне диапазона [0, %d)", index, len(nums))
	}
	result := make([]int, 0, len(nums)-1)
	result = append(result, nums[:index]...)
	result = append(result, nums[index+1:]...)
	return result, nil
}

func removalExample(nums []int, index int) {
	result, err := RemoveElement(nums, index)
	fmt.Printf("RemoveElement(%v, %d) -> %v %v, исходный слайс: %v\n",
		nums, index, result, err, nums)
}

func main() {
	removalExample([]int{1, 2, 2, 3, 1}, 3)
	removalExample([]int{1, 2, 3}, 5)
	removalExample([]int{1, 2, 3}, -1)
}

/*
Вопрос: почему решение через append(nums[:index], nums[index+1:]...) испортит исходный слайс?

Ответ: nums[:index] использует тот же массив, что и исходный слайс.
При допустимом индексе его ёмкости хватает для append, поэтому элементы
после index записываются поверх исходных. Чтобы сохранить исходный слайс,
нужно выделить отдельный массив для результата.
*/
