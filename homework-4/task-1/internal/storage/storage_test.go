package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
)

type unsupportedItem struct {
	item.Item
}

type unsupportedStorable struct {
	item.Storable
}

func TestStorageRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name  string
		store interface {
			Load() ([]inventory.Entry, error)
			Save([]inventory.Entry) error
		}
	}{
		{"memory", &Memory{}},
		{"file", NewFile(filepath.Join(t.TempDir(), "nested", "inventory.txt"))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.store.Load(); err != nil || len(got) != 0 {
				t.Fatalf("initial Load() = (%v, %v), want empty storage", got, err)
			}
			armor := item.NewArmor("Щит", 4, 7)
			if _, err := armor.Use(); err != nil {
				t.Fatal(err)
			}
			entries := []inventory.Entry{
				{ID: 7, Item: item.NewWeapon("Меч", 3, 10, 5)},
				{ID: 2, Item: armor},
				{ID: 4, Item: item.NewPotion("Зелье", 0.5, "HP", 0)},
			}
			if err := tt.store.Save(entries); err != nil {
				t.Fatal(err)
			}
			const want = "Weapon|7|Меч|10|5|3\nArmor|2|Щит|7|4|true\nPotion|4|Зелье|HP|0|0.5\n"
			for _, invalid := range [][]inventory.Entry{
				{{ID: 0, Item: armor}},
				{entries[0], entries[0]},
				{{ID: 1, Item: (*item.Weapon)(nil)}},
				{{ID: 1, Item: item.NewWeapon("Меч", 3, 10, -1)}},
				{{ID: 1, Item: unsupportedItem{Item: armor}}},
				{{ID: 1, Item: unsupportedStorable{Storable: armor}}},
			} {
				if err := tt.store.Save(invalid); err == nil {
					t.Errorf("Save(%v) succeeded", invalid)
				}
				got, err := tt.store.Load()
				if err != nil {
					t.Fatal(err)
				}
				checkEntries(t, got, want)
			}
			if err := tt.store.Save(nil); err != nil {
				t.Fatal(err)
			}
			if got, err := tt.store.Load(); err != nil || len(got) != 0 {
				t.Errorf("Load() after clearing = (%v, %v), want empty storage", got, err)
			}
		})
	}
}

func TestMemoryIsolation(t *testing.T) {
	original := item.NewWeapon("Меч", 3, 10, 5)
	store := &Memory{}
	if err := store.Save([]inventory.Entry{{ID: 7, Item: original}}); err != nil {
		t.Fatal(err)
	}
	if _, err := original.Use(); err != nil {
		t.Fatal(err)
	}
	if original.Durability != 4 {
		t.Fatalf("original durability = %d, want 4", original.Durability)
	}

	entries, err := store.Load()
	if err != nil || len(entries) != 1 {
		t.Fatalf("Load() after changing original = (%v, %v), want 1 entry", entries, err)
	}
	loaded, ok := entries[0].Item.(*item.Weapon)
	if !ok {
		t.Fatalf("loaded item type = %T, want *item.Weapon", entries[0].Item)
	}
	if entries[0].ID != 7 || loaded.Durability != 5 {
		t.Errorf("loaded ID and durability = (%d, %d), want (7, 5)", entries[0].ID, loaded.Durability)
	}
	if _, err := loaded.Use(); err != nil {
		t.Fatal(err)
	}
	if loaded.Durability != 4 {
		t.Fatalf("loaded durability after Use() = %d, want 4", loaded.Durability)
	}

	entries, err = store.Load()
	if err != nil || len(entries) != 1 {
		t.Fatalf("Load() after changing loaded item = (%v, %v), want 1 entry", entries, err)
	}
	reloaded, ok := entries[0].Item.(*item.Weapon)
	if !ok {
		t.Fatalf("reloaded item type = %T, want *item.Weapon", entries[0].Item)
	}
	if entries[0].ID != 7 || reloaded.Durability != 5 {
		t.Errorf("reloaded ID and durability = (%d, %d), want (7, 5)", entries[0].ID, reloaded.Durability)
	}
}

func TestDecodeEntries(t *testing.T) {
	const first = "Weapon|7|Меч|10|5|3"
	const last = "Potion|4|Зелье|HP|0|0.5"
	entries, err := decodeEntries(strings.NewReader(first + "\r\n" + last))
	if err != nil {
		t.Fatal(err)
	}
	checkEntries(t, entries, first+"\n"+last+"\n")

	for _, tt := range []struct {
		name   string
		record string
		cause  string
	}{
		{"missing fields", "Weapon|4", "ожидается тип"},
		{"invalid ID", "Weapon|x|Меч|10|5|3", "ID:"},
		{"zero ID", "Weapon|0|Меч|10|5|3", "ID вне допустимого диапазона"},
		{"duplicate ID", first, "повторяющийся ID"},
		{"unknown type", "Other|4|Меч|10|5|3", "неизвестный тип"},
		{"invalid item", "Armor|4|Щит|7|4|yes", "неверное состояние брони"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeEntries(strings.NewReader(first + "\n" + tt.record + "\n" + last))
			var parseErr *ParseItemError
			if !errors.As(err, &parseErr) {
				t.Fatalf("decodeEntries() error = %v, want *ParseItemError", err)
			}
			if got != nil || parseErr.Line != 2 || parseErr.Value != tt.record {
				t.Errorf("decodeEntries() = (%v, %v), want no entries and error on line 2 for %q", got, parseErr, tt.record)
			}
			if parseErr.Err == nil || !errors.Is(parseErr, parseErr.Err) || !strings.Contains(parseErr.Err.Error(), tt.cause) {
				t.Errorf("underlying error = %v, want %q", parseErr.Err, tt.cause)
			}
		})
	}
}

type failingReader struct {
	reader io.Reader
	err    error
}

func (r failingReader) Read(data []byte) (int, error) {
	n, err := r.reader.Read(data)
	if err == io.EOF {
		return n, r.err
	}
	return n, err
}

func TestDecodeReadFailures(t *testing.T) {
	failure := errors.New("read failed")
	const first = "Weapon|1|Меч|10|5|3\n"
	for _, tt := range []struct {
		name   string
		reader io.Reader
		want   error
	}{
		{"before first record", failingReader{reader: strings.NewReader(""), err: failure}, failure},
		{"after first record", failingReader{reader: strings.NewReader(first), err: failure}, failure},
		{"oversized record", strings.NewReader(first + strings.Repeat("a", 64*1024)), nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := decodeEntries(tt.reader)
			if entries != nil || err == nil || (tt.want != nil && !errors.Is(err, tt.want)) {
				t.Errorf("decodeEntries() = (%v, %v), want no entries and %v", entries, err, tt.want)
			}
		})
	}
}

func TestFileLoadCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.txt")
	const invalid = "Armor|2|Щит|7|4|yes"
	data := "Weapon|1|Меч|10|5|3\n" + invalid
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := NewFile(path).Load()
	var parseErr *ParseItemError
	if got != nil || !errors.As(err, &parseErr) {
		t.Fatalf("Load() = (%v, %v), want no entries and a wrapped *ParseItemError", got, err)
	}
	if parseErr.Line != 2 || parseErr.Value != invalid {
		t.Errorf("Load() error = %v, want line 2 and record %q", parseErr, invalid)
	}
}

func TestFileWriteFailurePreservesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inventory.txt")
	const original = "Weapon|1|Меч|10|5|3\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("write failed")
	entries := []inventory.Entry{{ID: 2, Item: item.NewPotion("Зелье", 0.5, "HP", 3)}}
	err := NewFile(path).save(entries, func(writer io.Writer, data string) error {
		if _, err := io.WriteString(writer, data[:len(data)/2]); err != nil {
			return err
		}
		return failure
	})
	if !errors.Is(err, failure) {
		t.Errorf("save() = %v, want %v", err, failure)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != original {
		t.Errorf("original file = (%q, %v), want %q", got, err, original)
	}
	files, err := os.ReadDir(dir)
	if err != nil || len(files) != 1 || files[0].Name() != "inventory.txt" {
		t.Errorf("directory after failed write = (%v, %v)", files, err)
	}
}

func TestFileReplaceFailureCleansTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inventory.txt")
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := NewFile(path).Save(nil); err == nil {
		t.Error("Save() succeeded with a directory as the destination")
	}
	files, err := os.ReadDir(dir)
	if err != nil || len(files) != 1 || files[0].Name() != "inventory.txt" || !files[0].IsDir() {
		t.Errorf("directory after failed replacement = (%v, %v)", files, err)
	}
}

func checkEntries(t *testing.T, entries []inventory.Entry, want string) {
	t.Helper()
	got, err := encodeEntries(entries)
	if err != nil || got != want {
		t.Errorf("loaded entries encode as (%q, %v), want %q", got, err, want)
	}
}
