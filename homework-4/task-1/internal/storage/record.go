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
	if err := value.Deserialize(strings.NewReader(strings.Join(fields[2:], "|"))); err != nil {
		return inventory.Entry{}, err
	}
	return inventory.Entry{ID: id, Item: value}, nil
}

func encodeEntries(entries []inventory.Entry) (string, error) {
	if err := validateEntries(entries); err != nil {
		return "", err
	}
	var data strings.Builder
	for _, entry := range entries {
		kind, err := recordType(entry.Item)
		if err != nil {
			return "", err
		}
		if _, err := fmt.Fprintf(&data, "%s|%d|", kind, entry.ID); err != nil {
			return "", err
		}
		storable, ok := entry.Item.(item.Storable)
		if !ok {
			return "", errors.New("предмет не поддерживает сохранение")
		}
		if err := storable.Serialize(&data); err != nil {
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

func validateEntries(entries []inventory.Entry) error {
	ids := make(map[int]bool)
	for _, entry := range entries {
		if !validID(entry.ID) {
			return errors.New("ID вне допустимого диапазона")
		}
		if ids[entry.ID] {
			return fmt.Errorf("сохранение: повторяющийся ID %d", entry.ID)
		}
		ids[entry.ID] = true
		if item.IsNil(entry.Item) {
			return errors.New("сохранение: предмет отсутствует")
		}
		storable, ok := entry.Item.(item.Storable)
		if !ok {
			return fmt.Errorf("сохранение %q не поддерживается", entry.Item.Name())
		}
		if err := storable.Validate(); err != nil {
			return fmt.Errorf("сохранение: %w", err)
		}
		if _, err := recordType(entry.Item); err != nil {
			return err
		}
	}
	return nil
}
