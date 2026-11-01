package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"homeworks/homework-4/task-1/internal/inventory"
)

// File сохраняет инвентарь в файл, заменяя его только после успешной записи.
type File struct {
	path string
}

// NewFile задаёт путь; каталог данных создаётся только при сохранении.
func NewFile(path string) *File {
	return &File{path: path}
}

// Load возвращает пустой инвентарь, если файл отсутствует.
func (s *File) Load() (entries []inventory.Entry, err error) {
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("открытие %q: %w", s.path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("закрытие %q: %w", s.path, closeErr))
			entries = nil
		}
	}()
	entries, err = decodeEntries(file)
	if err != nil {
		return nil, fmt.Errorf("загрузка %q: %w", s.path, err)
	}
	return entries, nil
}

// Save проверяет записи и сохраняет их через временный файл.
func (s *File) Save(entries []inventory.Entry) error {
	return s.save(entries, writeData)
}

// Функция записи передаётся отдельно, чтобы проверить отказ посередине записи.
func (s *File) save(entries []inventory.Entry, write func(io.Writer, string) error) (err error) {
	data, err := encodeEntries(entries)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return fmt.Errorf("каталог для %q: %w", s.path, err)
	}
	// Временный файл создаётся рядом с основным: замена не пересекает файловые системы.
	file, err := os.CreateTemp(filepath.Dir(s.path), ".inventory-*.tmp")
	if err != nil {
		return fmt.Errorf("временный файл для %q: %w", s.path, err)
	}
	tempPath := file.Name()
	defer func() {
		if file != nil {
			if closeErr := file.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("закрытие временного файла: %w", closeErr))
			}
		}
		// После успешной замены временного файла уже нет.
		if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			err = errors.Join(err, fmt.Errorf("удаление временного файла: %w", removeErr))
		}
	}()
	if err := write(file, data); err != nil {
		return fmt.Errorf("запись для %q: %w", s.path, err)
	}
	closeErr := file.Close()
	file = nil
	if closeErr != nil {
		return fmt.Errorf("закрытие временного файла: %w", closeErr)
	}
	if err := os.Rename(tempPath, s.path); err != nil {
		return fmt.Errorf("замена %q: %w", s.path, err)
	}
	return nil
}

func writeData(writer io.Writer, data string) error {
	_, err := io.WriteString(writer, data)
	return err
}
