package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
)

// ParseItemError содержит номер, текст и причину ошибки повреждённой записи.
type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("строка %d %q: %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error {
	return e.Err
}

// Останавливаемся на первой ошибке и не возвращаем частичный инвентарь.
func decodeEntries(reader io.Reader) (entries []inventory.Entry, err error) {
	scanner := bufio.NewScanner(reader)
	ids := make(map[int]bool)
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		entry, parseErr := parseEntry(text)
		if parseErr == nil && ids[entry.ID] {
			parseErr = fmt.Errorf("повторяющийся ID %d", entry.ID)
		}
		if parseErr != nil {
			return nil, &ParseItemError{Line: line, Value: text, Err: parseErr}
		}
		ids[entry.ID] = true
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("чтение записей: %w", err)
	}
	return entries, nil
}

func parseEntry(text string) (inventory.Entry, error) {
	fields := strings.Split(text, "|")
	if len(fields) < 3 {
		return inventory.Entry{}, errors.New("ожидается тип, ID и поля предмета")
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return inventory.Entry{}, fmt.Errorf("ID: %w", err)
	}
	if id <= 0 || id == int(^uint(0)>>1) {
		return inventory.Entry{}, errors.New("ID вне допустимого диапазона")
	}
	value, err := item.ParseRecord(fields[0] + "|" + strings.Join(fields[2:], "|"))
	if err != nil {
		return inventory.Entry{}, err
	}
	return inventory.Entry{ID: id, Item: value}, nil
}

func encodeEntries(entries []inventory.Entry) (string, error) {
	// Сначала проверяем записи, чтобы ошибка формата не обнулила старый файл.
	var data strings.Builder
	ids := make(map[int]bool)
	for _, entry := range entries {
		if item.IsNil(entry.Item) {
			return "", errors.New("сохранение: предмет отсутствует")
		}
		storable, ok := entry.Item.(item.Storable)
		if !ok {
			return "", fmt.Errorf("сохранение %q не поддерживается", entry.Item.Name())
		}
		var record strings.Builder
		if err := storable.Serialize(&record); err != nil {
			return "", fmt.Errorf("сериализация: %w", err)
		}
		fields := strings.Split(record.String(), "|")
		line := fmt.Sprintf("%s|%d|%s", fields[0], entry.ID, strings.Join(fields[1:], "|"))
		if _, err := parseEntry(line[:len(line)-1]); err != nil {
			return "", fmt.Errorf("сохранение: %w", err)
		}
		if ids[entry.ID] {
			return "", fmt.Errorf("сохранение: повторяющийся ID %d", entry.ID)
		}
		ids[entry.ID] = true
		if _, err := data.WriteString(line); err != nil {
			return "", fmt.Errorf("подготовка записи: %w", err)
		}
	}
	return data.String(), nil
}
