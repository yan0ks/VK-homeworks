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
	buffered, ok := reader.(*bufio.Reader)
	if !ok {
		return nil, errors.New("Передайте bufio.NewReader(...). Для последовательного чтения предметов используйте один и тот же reader")
	}
	line, err := buffered.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(line) == 0 && err == io.EOF {
		return nil, io.EOF
	}
	if len(line) >= maxRecordBytes {
		return nil, errors.New("слишком длинная запись")
	}
	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	for _, symbol := range line {
		if symbol == '\n' || symbol == '\r' {
			return nil, errors.New("ожидается одна запись")
		}
	}
	return strings.Split(line, "|"), nil
}
