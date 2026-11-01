package item

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Storable interface {
	Item
	Validate() error
	Serialize(io.Writer) error
	Deserialize(io.Reader) error
}

var (
	_ Storable = (*Weapon)(nil)
	_ Storable = (*Armor)(nil)
	_ Storable = (*Potion)(nil)
)

func validText(value string) bool {
	// Даже строка зелья с двумя текстовыми полями должна помещаться в Scanner.
	const maxTextBytes = 16 * 1024
	if len(value) > maxTextBytes || strings.TrimSpace(value) == "" {
		return false
	}
	for _, symbol := range value {
		if symbol == '|' || symbol == '\r' || symbol == '\n' {
			return false
		}
	}
	return true
}

func validWeight(weight float64) bool {
	// У NaN и бесконечности разность с собой не равна нулю.
	return weight >= 0 && weight-weight == 0
}

func validateBase(name string, weight float64) error {
	if !validText(name) {
		return errors.New("неверное имя предмета")
	}
	if !validWeight(weight) {
		return errors.New("вес должен быть конечным неотрицательным числом")
	}
	return nil
}

func readFields(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return strings.Split(scanner.Text(), "|"), nil
}
