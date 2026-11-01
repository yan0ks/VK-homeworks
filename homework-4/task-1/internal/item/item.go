package item

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
)

type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Armor)(nil)
	_ Item = (*Potion)(nil)
)

func IsNil(value Item) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Pointer && v.IsNil()
}

func DescribeItem(value Item) string {
	if IsNil(value) {
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%s (вес: %g)", value.Name(), value.Weight())
}

func UseDangerousItem(value Item) (message string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			switch cause := recovered.(type) {
			case error:
				err = fmt.Errorf("опасный предмет: %w", cause)
			default:
				err = fmt.Errorf("опасный предмет: %v", cause)
			}
		}
	}()
	if IsNil(value) {
		return "", errors.New("предмет отсутствует")
	}
	if value.Name() == "Ящик Пандоры" {
		panic("ящик Пандоры открыт")
	}
	return value.Use()
}
