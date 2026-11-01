package cli

import (
	"fmt"
	"io"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
)

// Примеры из ДЗ4 выполняются отдельно от постоянного инвентаря CLI.
func demo(stdout io.Writer) error {
	sword := item.NewWeapon("Меч", 3, 10, 5)
	shield := item.NewArmor("Щит", 4.5, 5)
	potion := item.NewPotion("Лечебное", 0.5, "+50 HP", 3)
	bow := item.NewWeapon("Сломанный лук", 2, 5, 0)
	inv := inventory.New(nil)
	for _, value := range []item.Item{sword, shield, potion, bow} {
		inv.AddItem(value)
	}
	for _, value := range []item.Item{sword, shield, potion} {
		message, err := value.Use()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, message); err != nil {
			return err
		}
	}
	var weapon *item.Weapon
	for _, value := range []item.Item{sword, nil, weapon} {
		if _, err := fmt.Fprintln(stdout, item.DescribeItem(value)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(stdout, "Оружие:"); err != nil {
		return err
	}
	for _, value := range inv.Weapons() {
		if _, err := fmt.Fprintln(stdout, item.DescribeItem(value)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(stdout, "Сломанные предметы:"); err != nil {
		return err
	}
	for _, value := range inv.BrokenItems() {
		if _, err := fmt.Fprintln(stdout, item.DescribeItem(value)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(stdout, "Названия:", inv.ItemNames()); err != nil {
		return err
	}
	value, found := inv.FindItemByName("Щит")
	_, err := fmt.Fprintln(stdout, "Поиск Щит:", item.DescribeItem(value), found)
	return err
}
