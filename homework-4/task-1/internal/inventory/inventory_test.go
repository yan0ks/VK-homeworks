package inventory

import (
	"errors"
	"slices"
	"testing"

	"homeworks/homework-4/task-1/internal/item"
)

func TestAddItem(t *testing.T) {
	weapon := item.NewWeapon("меч", 2, 5, 1)
	armor := item.NewArmor("щит", 3, 4)
	potion := item.NewPotion("зелье", 1, "HP", 1)
	for _, tt := range []struct {
		name    string
		entries []Entry
		wantID  int
	}{
		{"empty inventory", nil, 1},
		{"unsorted IDs", []Entry{{ID: 8, Item: weapon}, {ID: 3, Item: armor}}, 9},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inv := New(tt.entries)
			if got := inv.AddItem(potion); got != tt.wantID {
				t.Errorf("AddItem() ID = %d, want %d", got, tt.wantID)
			}
			want := append(tt.entries, Entry{ID: tt.wantID, Item: potion})
			if got := inv.Entries(); !slices.Equal(got, want) {
				t.Errorf("Entries() = %v, want %v", got, want)
			}
		})
	}
}

func TestEntriesCopy(t *testing.T) {
	entry := Entry{ID: 4, Item: item.NewArmor("щит", 3, 4)}
	input := []Entry{entry}
	inv := New(input)
	input[0] = Entry{}
	output := inv.Entries()
	output[0] = Entry{}
	if got := inv.Entries(); !slices.Equal(got, []Entry{entry}) {
		t.Errorf("Entries() after changing input and output = %v, want [%v]", got, entry)
	}
}

func TestFindByID(t *testing.T) {
	armor := item.NewArmor("щит", 3, 4)
	inv := New([]Entry{{ID: 4, Item: armor}})
	for _, tt := range []struct {
		name    string
		id      int
		want    item.Item
		wantErr error
	}{
		{"existing ID", 4, armor, nil},
		{"missing ID", 5, nil, ErrNotFound},
		{"zero ID", 0, nil, ErrNotFound},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inv.FindByID(tt.id)
			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Errorf("FindByID(%d) = (%v, %v), want (%v, %v)", tt.id, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	first := Entry{ID: 1, Item: item.NewWeapon("меч", 2, 5, 1)}
	second := Entry{ID: 2, Item: item.NewArmor("щит", 3, 4)}
	last := Entry{ID: 3, Item: item.NewPotion("зелье", 1, "HP", 1)}
	entries := []Entry{first, second, last}
	for _, tt := range []struct {
		name    string
		id      int
		want    []Entry
		wantErr error
	}{
		{"first", 1, []Entry{second, last}, nil},
		{"middle", 2, []Entry{first, last}, nil},
		{"last", 3, []Entry{first, second}, nil},
		{"missing", 4, entries, ErrNotFound},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inv := New(entries)
			if err := inv.Delete(tt.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete(%d) = %v, want %v", tt.id, err, tt.wantErr)
			}
			if got := inv.Entries(); !slices.Equal(got, tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
			}
			if id := inv.AddItem(last.Item); id != 4 {
				t.Errorf("AddItem() after Delete() ID = %d, want 4", id)
			}
		})
	}
}

func TestItemQueries(t *testing.T) {
	weapon := item.NewWeapon("меч", 2, 5, 1)
	broken := item.NewWeapon("меч", 3, 7, 0)
	armor := item.NewArmor("щит", 3, 4)
	potion := item.NewPotion("зелье", 1, "HP", 1)
	exhausted := item.NewPotion("пустое зелье", 1, "HP", 0)
	inv := New(nil)
	for _, value := range []item.Item{nil, weapon, armor, (*item.Weapon)(nil), broken, potion, exhausted, (*item.Potion)(nil)} {
		inv.AddItem(value)
	}
	if got, want := inv.Weapons(), []*item.Weapon{weapon, broken}; !slices.Equal(got, want) {
		t.Errorf("Weapons() = %v, want %v", got, want)
	}
	if got, want := inv.BrokenItems(), []item.Item{broken, exhausted}; !slices.Equal(got, want) {
		t.Errorf("BrokenItems() = %v, want %v", got, want)
	}
	wantNames := []string{"Предмет отсутствует", "меч", "щит", "Предмет отсутствует", "меч", "зелье", "пустое зелье", "Предмет отсутствует"}
	if got := inv.ItemNames(); !slices.Equal(got, wantNames) {
		t.Errorf("ItemNames() = %q, want %q", got, wantNames)
	}
	for _, tt := range []struct {
		name  string
		want  item.Item
		found bool
	}{
		{"меч", weapon, true},
		{"неизвестный предмет", nil, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, found := inv.FindItemByName(tt.name)
			if got != tt.want || found != tt.found {
				t.Errorf("FindItemByName(%q) = (%v, %t), want (%v, %t)", tt.name, got, found, tt.want, tt.found)
			}
		})
	}
}
