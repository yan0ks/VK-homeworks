package item

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Potion struct {
	name    string
	weight  float64
	Effect  string
	Charges int
}

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

func (p *Potion) Validate() error {
	if err := validateBase(p.name, p.weight); err != nil {
		return err
	}
	if p.Charges < 0 {
		return errors.New("заряды не могут быть отрицательными")
	}
	if !validText(p.Effect) {
		return errors.New("неверный эффект зелья")
	}
	return nil
}

func (p *Potion) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s|%s|%d|%g\n", p.name, p.Effect, p.Charges, p.weight)
	return err
}

func (p *Potion) Deserialize(reader io.Reader) error {
	fields, err := readFields(reader)
	if err != nil {
		return err
	}
	if len(fields) != 4 {
		return errors.New("неверный формат зелья")
	}
	charges, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("заряды: %w", err)
	}
	weight, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("вес: %w", err)
	}
	value := NewPotion(fields[0], weight, fields[1], charges)
	if err := value.Validate(); err != nil {
		return err
	}
	*p = *value
	return nil
}
