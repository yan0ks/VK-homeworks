package inventory

import (
	"errors"
	"fmt"

	"homeworks/homework-4/task-1/internal/collection"
	"homeworks/homework-4/task-1/internal/item"
)

var ErrNotFound = errors.New("предмет не найден")

type Entry struct {
	ID   int
	Item item.Item
}

type Inventory struct {
	entries []Entry
	nextID  int
}

func New(entries []Entry) *Inventory {
	inv := &Inventory{nextID: 1}
	inv.entries = append(inv.entries, entries...)
	for _, entry := range entries {
		if entry.ID >= inv.nextID {
			inv.nextID = entry.ID + 1
		}
	}
	return inv
}

func (inv *Inventory) AddItem(value item.Item) int {
	id := inv.nextID
	inv.nextID++
	inv.entries = append(inv.entries, Entry{ID: id, Item: value})
	return id
}

func (inv *Inventory) Entries() []Entry {
	return append([]Entry(nil), inv.entries...)
}

func (inv *Inventory) Weapons() []*item.Weapon {
	weapons := collection.Filter(inv.values(), func(value item.Item) bool {
		weapon, ok := value.(*item.Weapon)
		return ok && weapon != nil
	})
	return collection.Map(weapons, func(value item.Item) *item.Weapon {
		weapon, _ := value.(*item.Weapon)
		return weapon
	})
}

func (inv *Inventory) BrokenItems() []item.Item {
	return collection.Filter(inv.values(), func(value item.Item) bool {
		switch v := value.(type) {
		case *item.Weapon:
			return v != nil && v.Durability <= 0
		case *item.Potion:
			return v != nil && v.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) ItemNames() []string {
	return collection.Map(inv.values(), func(value item.Item) string {
		if item.IsNil(value) {
			return "Предмет отсутствует"
		}
		return value.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (item.Item, bool) {
	return collection.Find(inv.values(), func(value item.Item) bool {
		return !item.IsNil(value) && value.Name() == name
	})
}

func (inv *Inventory) FindByID(id int) (item.Item, error) {
	entry, ok := collection.Find(inv.entries, func(entry Entry) bool { return entry.ID == id })
	if !ok {
		return nil, fmt.Errorf("ID %d: %w", id, ErrNotFound)
	}
	return entry.Item, nil
}

func (inv *Inventory) Delete(id int) error {
	for index, entry := range inv.entries {
		if entry.ID == id {
			inv.entries = append(inv.entries[:index], inv.entries[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("ID %d: %w", id, ErrNotFound)
}

func (inv *Inventory) values() []item.Item {
	return collection.Map(inv.entries, func(entry Entry) item.Item {
		return entry.Item
	})
}
