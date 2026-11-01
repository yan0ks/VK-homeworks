package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
	"homeworks/homework-4/task-1/internal/storage"
)

var errUsage = errors.New("неверные аргументы")

type store interface {
	Load() ([]inventory.Entry, error)
	Save([]inventory.Entry) error
}

type App struct {
	store store
}

func New(s store) *App {
	return &App{store: s}
}

func (app *App) Run(args []string, stdout, stderr io.Writer) int {
	err := app.run(args, stdout, stderr)
	if err == nil || errors.Is(err, flag.ErrHelp) {
		return 0
	}
	var parseErr *storage.ParseItemError
	switch {
	case errors.As(err, &parseErr):
		fmt.Fprintf(stderr, "Ошибка загрузки: строка %d, содержимое %q: %v\n",
			parseErr.Line, parseErr.Value, parseErr.Err)
	case errors.Is(err, item.ErrBroken):
		fmt.Fprintln(stderr, "Ошибка: оружие сломано")
	case errors.Is(err, item.ErrExhausted):
		fmt.Fprintln(stderr, "Ошибка: зелье закончилось")
	default:
		fmt.Fprintln(stderr, "Ошибка:", err)
	}
	if errors.Is(err, errUsage) {
		return 2
	}
	return 1
}

func (app *App) run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return demo(stdout)
	}
	switch args[0] {
	case "list":
		return app.list(args[1:], stdout, stderr)
	case "add":
		return app.add(args[1:], stdout, stderr)
	case "use", "delete":
		return app.change(args[0], args[1:], stdout, stderr)
	case "demo":
		flags := newFlags("demo", stderr)
		if err := parseFlags(flags, args[1:]); err != nil {
			return err
		}
		return demo(stdout)
	case "help", "-h", "--help":
		flags := newFlags("help", stderr)
		if err := parseFlags(flags, args[1:]); err != nil {
			return err
		}
		_, err := fmt.Fprintln(stdout, "Команды: list, add, use, delete, demo. Справка: <команда> -h")
		return err
	default:
		return fmt.Errorf("%w: неизвестная команда %q", errUsage, args[0])
	}
}

func (app *App) load() (*inventory.Inventory, error) {
	entries, err := app.store.Load()
	if err != nil {
		return nil, err
	}
	return inventory.New(entries), nil
}

func (app *App) list(args []string, stdout, stderr io.Writer) error {
	flags := newFlags("list", stderr)
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	inv, err := app.load()
	if err != nil {
		return err
	}
	if len(inv.Entries()) == 0 {
		_, err := fmt.Fprintln(stdout, "Инвентарь пуст")
		return err
	}
	for _, entry := range inv.Entries() {
		description := item.DescribeItem(entry.Item)
		if _, err := fmt.Fprintf(stdout, "ID: %d — %s\n", entry.ID, description); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) add(args []string, stdout, stderr io.Writer) error {
	flags := newFlags("add", stderr)
	name := flags.String("name", "", "название предмета")
	weight := flags.Float64("weight", 0, "вес предмета")
	kind := flags.String("type", "weapon", "тип: weapon, armor, potion")
	damage := flags.Int("damage", 10, "урон оружия")
	durability := flags.Int("durability", 5, "прочность оружия")
	defense := flags.Int("defense", 5, "защита брони")
	effect := flags.String("effect", "+50 HP", "эффект зелья")
	charges := flags.Int("charges", 3, "заряды зелья")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	var value item.Storable
	switch *kind {
	case "weapon":
		value = item.NewWeapon(*name, *weight, *damage, *durability)
	case "armor":
		value = item.NewArmor(*name, *weight, *defense)
	case "potion":
		value = item.NewPotion(*name, *weight, *effect, *charges)
	default:
		return fmt.Errorf("%w: неизвестный тип %q", errUsage, *kind)
	}
	if err := value.Validate(); err != nil {
		return fmt.Errorf("%w: %w", errUsage, err)
	}
	inv, err := app.load()
	if err != nil {
		return err
	}
	id := inv.AddItem(value)
	if err := app.store.Save(inv.Entries()); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Добавлен %s, ID: %d\n", value.Name(), id)
	return err
}

func (app *App) change(command string, args []string, stdout, stderr io.Writer) error {
	flags := newFlags(command, stderr)
	id := flags.Int("id", 0, "ID предмета")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if *id <= 0 {
		return fmt.Errorf("%w: -id должен быть положительным", errUsage)
	}
	inv, err := app.load()
	if err != nil {
		return err
	}
	var message string
	if command == "delete" {
		if err := inv.Delete(*id); err != nil {
			return err
		}
		message = fmt.Sprintf("Удалён предмет ID: %d", *id)
	} else {
		value, err := inv.FindByID(*id)
		if err != nil {
			return err
		}
		message, err = item.UseDangerousItem(value)
		if err != nil {
			return err
		}
	}
	if err := app.store.Save(inv.Entries()); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, message)
	return err
}

func newFlags(name string, stderr io.Writer) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	return flags
}

func parseFlags(flags *flag.FlagSet, args []string) error {
	output := flags.Output()
	var details strings.Builder
	flags.SetOutput(&details)
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			if _, writeErr := fmt.Fprintf(output, "%s", details.String()); writeErr != nil {
				return writeErr
			}
		}
		return fmt.Errorf("%w: %w", errUsage, err)
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("%w: лишние аргументы %v", errUsage, flags.Args())
	}
	return nil
}
