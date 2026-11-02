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
	if !validID(id) {
		return inventory.Entry{}, errors.New("ID вне допустимого диапазона")
	}
	var value item.Storable
	switch fields[0] {
	case "Weapon":
		value = &item.Weapon{}
	case "Armor":
		value = &item.Armor{}
	case "Potion":
		value = &item.Potion{}
	default:
		return inventory.Entry{}, fmt.Errorf("неизвестный тип предмета %q", fields[0])
	}
	reader := bufio.NewReader(strings.NewReader(strings.Join(fields[2:], "|")))
	if err := value.Deserialize(reader); err != nil {
		return inventory.Entry{}, err
	}
	return inventory.Entry{ID: id, Item: value}, nil
}

func encodeEntries(entries []inventory.Entry) (string, error) {
	prepared, err := prepareEntries(entries)
	if err != nil {
		return "", err
	}
	var data strings.Builder
	for _, entry := range prepared {
		if _, err := fmt.Fprintf(&data, "%s|%d|", entry.kind, entry.id); err != nil {
			return "", err
		}
		if err := entry.value.Serialize(&data); err != nil {
			return "", fmt.Errorf("сериализация: %w", err)
		}
	}
	return data.String(), nil
}

func validID(id int) bool { return id > 0 && id != int(^uint(0)>>1) }

func recordType(value item.Item) (string, error) {
	switch value.(type) {
	case *item.Weapon:
		return "Weapon", nil
	case *item.Armor:
		return "Armor", nil
	case *item.Potion:
		return "Potion", nil
	default:
		return "", errors.New("неподдерживаемый тип предмета")
	}
}

type preparedEntry struct {
	id    int
	kind  string
	value item.Storable
}

func prepareEntries(entries []inventory.Entry) ([]preparedEntry, error) {
	prepared := make([]preparedEntry, 0, len(entries))
	ids := make(map[int]bool)
	for _, entry := range entries {
		if !validID(entry.ID) {
			return nil, errors.New("ID вне допустимого диапазона")
		}
		if ids[entry.ID] {
			return nil, fmt.Errorf("сохранение: повторяющийся ID %d", entry.ID)
		}
		ids[entry.ID] = true
		if item.IsNil(entry.Item) {
			return nil, errors.New("сохранение: предмет отсутствует")
		}
		storable, ok := entry.Item.(item.Storable)
		if !ok {
			return nil, fmt.Errorf("сохранение %q не поддерживается", entry.Item.Name())
		}
		if err := storable.Validate(); err != nil {
			return nil, fmt.Errorf("сохранение: %w", err)
		}
		kind, err := recordType(entry.Item)
		if err != nil {
			return nil, err
		}
		prepared = append(prepared, preparedEntry{id: entry.ID, kind: kind, value: storable})
	}
	return prepared, nil
}
