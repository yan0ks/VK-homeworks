package item

// Filter возвращает новый слайс с элементами, прошедшими проверку.
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, value := range items {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

// Map возвращает результаты преобразования каждого элемента.
func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, len(items))
	for index, value := range items {
		result[index] = transform(value)
	}
	return result
}

// Find возвращает первый подходящий элемент или нулевое значение и false.
func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, value := range items {
		if condition(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}
