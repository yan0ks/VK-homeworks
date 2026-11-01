package item

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Weapon struct {
	name       string
	weight     float64
	Damage     int
	Durability int
}

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

func (w *Weapon) Validate() error {
	if err := validateBase(w.name, w.weight); err != nil {
		return err
	}
	if w.Damage < 0 || w.Durability < 0 {
		return errors.New("урон и прочность не могут быть отрицательными")
	}
	return nil
}

func (w *Weapon) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s|%d|%d|%g\n", w.name, w.Damage, w.Durability, w.weight)
	return err
}

func (w *Weapon) Deserialize(reader io.Reader) error {
	fields, err := readFields(reader)
	if err != nil {
		return err
	}
	if len(fields) != 4 {
		return errors.New("неверный формат оружия")
	}
	damage, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("урон: %w", err)
	}
	durability, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("прочность: %w", err)
	}
	weight, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("вес: %w", err)
	}
	value := NewWeapon(fields[0], weight, damage, durability)
	if err := value.Validate(); err != nil {
		return err
	}
	*w = *value
	return nil
}
