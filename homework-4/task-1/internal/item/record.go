package item

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Storable задаёт сохранение и восстановление предмета с обработкой ошибок из ДЗ5.
type Storable interface {
	Serialize(io.Writer) error
	Deserialize(io.Reader) error
}

var (
	_ Storable = (*Weapon)(nil)
	_ Storable = (*Armor)(nil)
	_ Storable = (*Potion)(nil)
)

func (w *Weapon) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "Weapon|%s|%d|%d|%g\n", w.name, w.Damage, w.Durability, w.weight)
	return err
}

func (w *Weapon) Deserialize(reader io.Reader) error {
	value, err := readRecord(reader)
	if err != nil {
		return err
	}
	weapon, ok := value.(*Weapon)
	if !ok {
		return errors.New("ожидалось оружие")
	}
	*w = *weapon
	return nil
}

func (a *Armor) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "Armor|%s|%d|%g|%t\n", a.name, a.Defense, a.weight, a.equipped)
	return err
}

func (a *Armor) Deserialize(reader io.Reader) error {
	value, err := readRecord(reader)
	if err != nil {
		return err
	}
	armor, ok := value.(*Armor)
	if !ok {
		return errors.New("ожидалась броня")
	}
	*a = *armor
	return nil
}

func (p *Potion) Serialize(out io.Writer) error {
	_, err := fmt.Fprintf(out, "Potion|%s|%s|%d|%g\n", p.name, p.Effect, p.Charges, p.weight)
	return err
}

func (p *Potion) Deserialize(reader io.Reader) error {
	value, err := readRecord(reader)
	if err != nil {
		return err
	}
	potion, ok := value.(*Potion)
	if !ok {
		return errors.New("ожидалось зелье")
	}
	*p = *potion
	return nil
}

// ValidText проверяет текстовое поле формата с разделителем |.
func ValidText(value string) bool {
	// Даже строка зелья с двумя текстовыми полями должна помещаться в Scanner.
	const maxTextBytes = 16 * 1024
	if len(value) > maxTextBytes || strings.TrimSpace(value) == "" {
		return false
	}
	for _, symbol := range value {
		if symbol == '|' || symbol == '\r' || symbol == '\n' {
			return false
		}
	}
	return true
}

// ValidWeight допускает только конечный неотрицательный вес.
func ValidWeight(weight float64) bool {
	// У NaN и бесконечности разность с собой не равна нулю.
	return weight >= 0 && weight-weight == 0
}

// ParseRecord восстанавливает предмет из записи без ID.
func ParseRecord(line string) (Item, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 4 || !ValidText(fields[1]) {
		return nil, errors.New("неверная запись предмета")
	}
	var value Item
	switch fields[0] {
	case "Weapon":
		if len(fields) != 5 {
			return nil, errors.New("неверный формат оружия")
		}
		damage, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("урон: %w", err)
		}
		durability, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("прочность: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("вес: %w", err)
		}
		if damage < 0 || durability < 0 {
			return nil, errors.New("урон и прочность не могут быть отрицательными")
		}
		value = NewWeapon(fields[1], weight, damage, durability)
	case "Armor":
		if len(fields) != 5 {
			return nil, errors.New("неверный формат брони")
		}
		defense, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("защита: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, fmt.Errorf("вес: %w", err)
		}
		if defense < 0 {
			return nil, errors.New("защита не может быть отрицательной")
		}
		a := NewArmor(fields[1], weight, defense)
		switch fields[4] {
		case "true":
			a.equipped = true
		case "false":
		default:
			return nil, errors.New("неверное состояние брони")
		}
		value = a
	case "Potion":
		if len(fields) != 5 || !ValidText(fields[2]) {
			return nil, errors.New("неверный формат зелья")
		}
		charges, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("заряды: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("вес: %w", err)
		}
		if charges < 0 {
			return nil, errors.New("заряды не могут быть отрицательными")
		}
		value = NewPotion(fields[1], weight, fields[2], charges)
	default:
		return nil, fmt.Errorf("неизвестный тип предмета %q", fields[0])
	}
	if !ValidWeight(value.Weight()) {
		return nil, errors.New("вес должен быть конечным неотрицательным числом")
	}
	return value, nil
}

func readRecord(reader io.Reader) (Item, error) {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ParseRecord(scanner.Text())
}
