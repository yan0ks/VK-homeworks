package storage

import (
	"strings"

	"homeworks/homework-4/task-1/internal/inventory"
)

// Memory хранит снимок данных в том же формате, что и File.
// Возвращаемые предметы не разделяют указатели с сохранёнными данными.
// Нулевое значение Memory готово к работе.
type Memory struct {
	data string
}

func (s *Memory) Load() ([]inventory.Entry, error) {
	return decodeEntries(strings.NewReader(s.data))
}

func (s *Memory) Save(entries []inventory.Entry) error {
	data, err := encodeEntries(entries)
	if err != nil {
		return err
	}
	s.data = data
	return nil
}
