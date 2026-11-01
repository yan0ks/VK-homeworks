package item

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrBroken возвращается при использовании сломанного оружия.
	ErrBroken = errors.New("item is broken")
	// ErrExhausted возвращается при использовании зелья без зарядов.
	ErrExhausted = errors.New("item is exhausted")
)

// Item задаёт общие свойства и использование предмета.
type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

// Weapon теряет единицу прочности при каждом успешном использовании.
type Weapon struct {
	name       string
	weight     float64
	Damage     int
	Durability int
}

// Armor сохраняет состояние: надета броня или снята.
type Armor struct {
	name     string
	weight   float64
	Defense  int
	equipped bool
}

// Potion расходует один заряд при каждом успешном использовании.
type Potion struct {
	name    string
	weight  float64
	Effect  string
	Charges int
}

var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Armor)(nil)
	_ Item = (*Potion)(nil)
)

// NewWeapon создаёт оружие с заданными характеристиками.
func NewWeapon(name string, weight float64, damage, durability int) *Weapon {
	return &Weapon{name: name, weight: weight, Damage: damage, Durability: durability}
}

func (w *Weapon) Name() string {
	return w.name
}

func (w *Weapon) Weight() float64 {
	return w.weight
}

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", ErrBroken
	}
	w.Durability--
	return fmt.Sprintf(
		"Используется %s (урон: %d, прочность: %d)", w.name, w.Damage, w.Durability,
	), nil
}

// NewArmor создаёт снятую броню с заданной защитой.
func NewArmor(name string, weight float64, defense int) *Armor {
	return &Armor{name: name, weight: weight, Defense: defense}
}

func (a *Armor) Name() string {
	return a.name
}

func (a *Armor) Weight() float64 {
	return a.weight
}

func (a *Armor) Use() (string, error) {
	a.equipped = !a.equipped
	return fmt.Sprintf("Броня %s (защита: %d, надета: %t)", a.name, a.Defense, a.equipped), nil
}

// NewPotion создаёт зелье с заданным эффектом и числом зарядов.
func NewPotion(name string, weight float64, effect string, charges int) *Potion {
	return &Potion{name: name, weight: weight, Effect: effect, Charges: charges}
}

func (p *Potion) Name() string {
	return p.name
}

func (p *Potion) Weight() float64 {
	return p.weight
}

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", ErrExhausted
	}
	p.Charges--
	return fmt.Sprintf("Используется %s (эффект: %s, заряды: %d)", p.name, p.Effect, p.Charges), nil
}

// IsNil распознаёт nil-интерфейс и nil-указатель внутри интерфейса Item.
func IsNil(value Item) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Pointer && v.IsNil()
}

// DescribeItem возвращает название и вес или сообщает об отсутствии предмета.
func DescribeItem(value Item) string {
	if IsNil(value) {
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%s (вес: %g)", value.Name(), value.Weight())
}

// UseDangerousItem преобразует учебную панику Ящика Пандоры в ошибку.
func UseDangerousItem(value Item) (message string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			switch cause := recovered.(type) {
			case error:
				err = fmt.Errorf("опасный предмет: %w", cause)
			default:
				err = fmt.Errorf("опасный предмет: %v", cause)
			}
		}
	}()
	if IsNil(value) {
		return "", errors.New("предмет отсутствует")
	}
	if value.Name() == "Ящик Пандоры" {
		panic("ящик Пандоры открыт")
	}
	return value.Use()
}
