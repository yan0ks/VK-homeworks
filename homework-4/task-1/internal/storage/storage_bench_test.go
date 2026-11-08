package storage

import (
	"path/filepath"
	"testing"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
)

func BenchmarkStorageLoad(b *testing.B) {
	const groups = 100
	entries := make([]inventory.Entry, 0, groups*3)
	for group := 0; group < groups; group++ {
		entries = append(entries,
			inventory.Entry{ID: group*3 + 1, Item: item.NewWeapon("Меч", 3, 10, 5)},
			inventory.Entry{ID: group*3 + 2, Item: item.NewArmor("Щит", 4, 7)},
			inventory.Entry{ID: group*3 + 3, Item: item.NewPotion("Зелье", 0.5, "HP", 2)},
		)
	}
	want, err := encodeEntries(entries)
	if err != nil {
		b.Fatal(err)
	}

	for _, tt := range []struct {
		name  string
		store interface {
			Save([]inventory.Entry) error
			Load() ([]inventory.Entry, error)
		}
	}{
		{"memory", &Memory{}},
		{"file", NewFile(filepath.Join(b.TempDir(), "inventory.txt"))},
	} {
		b.Run(tt.name, func(b *testing.B) {
			if err := tt.store.Save(entries); err != nil {
				b.Fatal(err)
			}
			loaded, err := tt.store.Load()
			if err != nil {
				b.Fatal(err)
			}
			got, err := encodeEntries(loaded)
			if err != nil {
				b.Fatal(err)
			}
			if got != want {
				b.Fatal("loaded records differ from the benchmark input")
			}

			for b.Loop() {
				if _, err := tt.store.Load(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
