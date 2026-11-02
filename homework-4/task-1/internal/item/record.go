package item

import (
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
	const maxRecordBytes = 64 * 1024
	data, err := io.ReadAll(io.LimitReader(reader, maxRecordBytes))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, io.EOF
	}
	if len(data) >= maxRecordBytes {
		return nil, errors.New("слишком длинная запись")
	}
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	if len(data) > 0 && data[len(data)-1] == '\r' {
		data = data[:len(data)-1]
	}
	for _, symbol := range data {
		if symbol == '\n' || symbol == '\r' {
			return nil, errors.New("ожидается одна запись")
		}
	}
	return strings.Split(string(data), "|"), nil
}
