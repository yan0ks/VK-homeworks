package cli

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"testing"

	"homeworks/homework-4/task-1/internal/inventory"
	"homeworks/homework-4/task-1/internal/item"
	"homeworks/homework-4/task-1/internal/storage"
)

func TestInventoryCommands(t *testing.T) {
	memory := &storage.Memory{}
	app := New(memory)
	for _, step := range []struct {
		args []string
		out  string
		size int
	}{
		{[]string{"list"}, "Инвентарь пуст\n", 0},
		{[]string{"add", "-name", "Меч", "-weight", "2", "-damage", "7", "-durability", "1"}, "Добавлен Меч, ID: 1\n", 1},
		{[]string{"add", "-type", "armor", "-name", "Щит", "-defense", "8"}, "Добавлен Щит, ID: 2\n", 2},
		{[]string{"add", "-type", "potion", "-name", "Зелье", "-effect", "HP", "-charges", "2"}, "Добавлен Зелье, ID: 3\n", 3},
		{[]string{"list"}, "ID: 1 — Меч (вес: 2)\nID: 2 — Щит (вес: 0)\nID: 3 — Зелье (вес: 0)\n", 3},
		{[]string{"use", "-id", "1"}, "Используется Меч (урон: 7, прочность: 0)\n", 3},
		{[]string{"use", "-id", "2"}, "Броня Щит (защита: 8, надета: true)\n", 3},
		{[]string{"use", "-id", "3"}, "Используется Зелье (эффект: HP, заряды: 1)\n", 3},
		{[]string{"delete", "-id", "2"}, "Удалён предмет ID: 2\n", 2},
		{[]string{"list"}, "ID: 1 — Меч (вес: 2)\nID: 3 — Зелье (вес: 0)\n", 2},
	} {
		var stdout, stderr bytes.Buffer
		if code := app.Run(step.args, &stdout, &stderr); code != 0 || stderr.Len() != 0 || stdout.String() != step.out {
			t.Fatalf("Run(%v) = (%d, %q, %q), want (0, %q, empty stderr)", step.args, code, stdout.String(), stderr.String(), step.out)
		}
		entries, err := memory.Load()
		if err != nil || len(entries) != step.size {
			t.Fatalf("stored entries after %v = (%v, %v), want %d entries", step.args, entries, err, step.size)
		}
	}
	entries, err := memory.Load()
	if err != nil {
		t.Fatal(err)
	}
	weapon, weaponOK := entries[0].Item.(*item.Weapon)
	potion, potionOK := entries[1].Item.(*item.Potion)
	if !weaponOK || weapon.Durability != 0 || !potionOK || potion.Charges != 1 {
		t.Errorf("stored weapon and potion after use = (%+v, %+v)", weapon, potion)
	}
}

type stubStore struct {
	entries []inventory.Entry
	loadErr error
	saveErr error
	saves   int
}

func (s *stubStore) Load() ([]inventory.Entry, error) {
	return s.entries, s.loadErr
}

func (s *stubStore) Save(entries []inventory.Entry) error {
	s.saves++
	if s.saveErr != nil {
		return s.saveErr
	}
	s.entries = entries
	return nil
}

func TestCommandErrors(t *testing.T) {
	broken := item.NewWeapon("Меч", 2, 7, 0)
	exhausted := item.NewPotion("Зелье", 1, "HP", 0)
	failure := errors.New("хранилище недоступно")
	parseErr := &storage.ParseItemError{Line: 2, Value: "invalid", Err: errors.New("неверная запись")}
	for _, tt := range []struct {
		name    string
		args    []string
		loadErr error
		saveErr error
		code    int
		saves   int
		message string
	}{
		{"unknown command", []string{"unknown"}, nil, nil, 2, 0, "неизвестная команда"},
		{"unknown flag", []string{"list", "-unknown"}, nil, nil, 2, 0, "flag provided but not defined"},
		{"extra argument", []string{"list", "extra"}, nil, nil, 2, 0, "лишние аргументы"},
		{"extra demo argument", []string{"demo", "extra"}, nil, nil, 2, 0, "лишние аргументы"},
		{"extra help argument", []string{"help", "extra"}, nil, nil, 2, 0, "лишние аргументы"},
		{"unknown type", []string{"add", "-type", "unknown"}, nil, nil, 2, 0, "неизвестный тип"},
		{"missing name", []string{"add"}, nil, nil, 2, 0, "неверное имя"},
		{"invalid ID", []string{"use", "-id", "0"}, nil, nil, 2, 0, "-id должен быть положительным"},
		{"non-numeric ID", []string{"use", "-id", "abc"}, nil, nil, 2, 0, "invalid value"},
		{"use missing item", []string{"use", "-id", "3"}, nil, nil, 1, 0, "предмет не найден"},
		{"delete missing item", []string{"delete", "-id", "3"}, nil, nil, 1, 0, "предмет не найден"},
		{"broken weapon", []string{"use", "-id", "1"}, nil, nil, 1, 0, "оружие сломано"},
		{"exhausted potion", []string{"use", "-id", "2"}, nil, nil, 1, 0, "зелье закончилось"},
		{"corrupt records", []string{"list"}, parseErr, nil, 1, 0, "строка 2, содержимое \"invalid\": неверная запись"},
		{"load before add", []string{"add", "-name", "Щит"}, failure, nil, 1, 0, failure.Error()},
		{"load before delete", []string{"delete", "-id", "1"}, failure, nil, 1, 0, failure.Error()},
		{"save after add", []string{"add", "-name", "Щит"}, nil, failure, 1, 1, failure.Error()},
		{"save after delete", []string{"delete", "-id", "1"}, nil, failure, 1, 1, failure.Error()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			store := &stubStore{
				entries: []inventory.Entry{{ID: 1, Item: broken}, {ID: 2, Item: exhausted}},
				loadErr: tt.loadErr,
				saveErr: tt.saveErr,
			}
			var stdout, stderr bytes.Buffer
			code := New(store).Run(tt.args, &stdout, &stderr)
			if code != tt.code || stdout.Len() != 0 || !strings.Contains(stderr.String(), tt.message) {
				t.Errorf("Run(%v) = (%d, %q, %q), want (%d, empty stdout, stderr containing %q)", tt.args, code, stdout.String(), stderr.String(), tt.code, tt.message)
			}
			if store.saves != tt.saves {
				t.Errorf("Save calls = %d, want %d", store.saves, tt.saves)
			}
			if len(store.entries) != 2 || store.entries[0].Item != broken || store.entries[1].Item != exhausted || broken.Durability != 0 || exhausted.Charges != 0 {
				t.Errorf("stored entries changed: %v", store.entries)
			}
		})
	}
}

func TestDemoAndHelp(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		stdout string
		stderr string
	}{
		{"default demo", nil, "Поиск Щит: Щит (вес: 4.5) true", ""},
		{"explicit demo", []string{"demo"}, "Сломанный лук (вес: 2)", ""},
		{"help command", []string{"help"}, "Команды: list, add, use, delete, demo", ""},
		{"command help flag", []string{"add", "-h"}, "", "Usage of add:"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := New(nil).Run(tt.args, &stdout, &stderr)
			if code != 0 || !containsOutput(stdout.String(), tt.stdout) || !containsOutput(stderr.String(), tt.stderr) {
				t.Errorf("Run(%v) = (%d, %q, %q), want exit 0 and output matching (%q, %q)", tt.args, code, stdout.String(), stderr.String(), tt.stdout, tt.stderr)
			}
		})
	}
}

func containsOutput(got, want string) bool {
	if want == "" {
		return got == ""
	}
	return strings.Contains(got, want)
}

type failingWriter struct {
	err    error
	writes int
	failAt int
}

func (w *failingWriter) Write(data []byte) (int, error) {
	w.writes++
	if w.writes == w.failAt {
		return 0, w.err
	}
	return len(data), nil
}

func TestOutputWriteFailure(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{"list with items", []string{"list"}},
		{"help", []string{"help"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			failure := errors.New("ошибка записи вывода")
			stdout := &failingWriter{err: failure, failAt: 1}
			var stderr bytes.Buffer
			store := &stubStore{
				entries: []inventory.Entry{
					{ID: 1, Item: item.NewWeapon("Меч", 2, 7, 1)},
					{ID: 2, Item: item.NewArmor("Щит", 3, 4)},
				},
			}
			code := New(store).Run(tt.args, stdout, &stderr)
			if want := "Ошибка: " + failure.Error() + "\n"; code != 1 || stderr.String() != want {
				t.Errorf("Run(%v) = (%d, %q), want (1, %q)", tt.args, code, stderr.String(), want)
			}
			if stdout.writes != 1 || store.saves != 0 {
				t.Errorf("Run(%v) attempted %d writes and %d saves, want 1 write and 0 saves", tt.args, stdout.writes, store.saves)
			}
		})
	}
}

func TestPrintItemUsesStopsAtBrokenWeapon(t *testing.T) {
	first := item.NewPotion("Зелье", 1, "HP", 1)
	broken := item.NewWeapon("Меч", 2, 7, 0)
	last := item.NewWeapon("Лук", 2, 5, 1)
	var stdout bytes.Buffer
	if err := printItemUses(&stdout, []item.Item{first, broken, last}); !errors.Is(err, item.ErrBroken) {
		t.Errorf("printItemUses() = %v, want %v", err, item.ErrBroken)
	}
	if want := "Используется Зелье (эффект: HP, заряды: 0)\n"; stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
	if first.Charges != 0 || broken.Durability != 0 || last.Durability != 1 {
		t.Errorf("charges and durability = (%d, %d, %d), want (0, 0, 1)", first.Charges, broken.Durability, last.Durability)
	}
}

func TestDemoWriteFailures(t *testing.T) {
	baseline := &failingWriter{}
	if err := demo(baseline); err != nil || baseline.writes == 0 {
		t.Fatalf("demo() with working writer = %v after %d writes", err, baseline.writes)
	}
	failure := errors.New("ошибка записи демо")
	for attempt := 1; attempt <= baseline.writes; attempt++ {
		t.Run("write "+strconv.Itoa(attempt), func(t *testing.T) {
			stdout := &failingWriter{err: failure, failAt: attempt}
			var stderr bytes.Buffer
			code := New(nil).Run([]string{"demo"}, stdout, &stderr)
			if want := "Ошибка: " + failure.Error() + "\n"; code != 1 || stderr.String() != want {
				t.Errorf("Run(demo) = (%d, %q), want (1, %q)", code, stderr.String(), want)
			}
			if stdout.writes != attempt {
				t.Errorf("write attempts = %d, want %d", stdout.writes, attempt)
			}
		})
	}
}

func TestHelpFlagWriteFailure(t *testing.T) {
	failure := errors.New("ошибка записи справки")
	stderr := &failingWriter{err: failure, failAt: 1}
	flags := newFlags("add", stderr)
	if err := parseFlags(flags, []string{"-h"}); !errors.Is(err, failure) || stderr.writes != 1 {
		t.Errorf("parseFlags(-h) = %v after %d writes, want %v after 1 write", err, stderr.writes, failure)
	}

	stderr = &failingWriter{err: failure, failAt: 1}
	var stdout bytes.Buffer
	if code := New(nil).Run([]string{"add", "-h"}, &stdout, stderr); code != 1 || stdout.Len() != 0 {
		t.Errorf("Run(add -h) = (%d, %q), want (1, empty stdout)", code, stdout.String())
	}
}
