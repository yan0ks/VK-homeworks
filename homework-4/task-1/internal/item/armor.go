package item

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Armor struct {
	name     string
	weight   float64
	Defense  int
	equipped bool
}

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

func (a *Armor) Validate() error {
	if err := validateBase(a.name, a.weight); err != nil {
		return err
	}
	if a.Defense < 0 {
		return errors.New("защита не может быть отрицательной")
	}
	return nil
}

func (a *Armor) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s|%d|%g|%t\n", a.name, a.Defense, a.weight, a.equipped)
	return err
}

func (a *Armor) Deserialize(reader io.Reader) error {
	fields, err := readFields(reader)
	if err != nil {
		return err
	}
	if len(fields) != 4 {
		return errors.New("неверный формат брони")
	}
	defense, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("защита: %w", err)
	}
	weight, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("вес: %w", err)
	}
	value := NewArmor(fields[0], weight, defense)
	switch fields[3] {
	case "true":
		value.equipped = true
	case "false":
	default:
		return errors.New("неверное состояние брони")
	}
	if err := value.Validate(); err != nil {
		return err
	}
	*a = *value
	return nil
}
