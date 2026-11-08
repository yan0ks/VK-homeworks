package item

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestUse(t *testing.T) {
	for _, tt := range []struct {
		name      string
		item      Storable
		messages  []string
		lastError error
		record    string
	}{
		{"weapon", NewWeapon("меч", 2, 5, 1), []string{"Используется меч (урон: 5, прочность: 0)", ""}, ErrBroken, "меч|5|0|2\n"},
		{"potion", NewPotion("зелье", 1, "HP", 1), []string{"Используется зелье (эффект: HP, заряды: 0)", ""}, ErrExhausted, "зелье|HP|0|1\n"},
		{"armor", NewArmor("щит", 3, 4), []string{"Броня щит (защита: 4, надета: true)", "Броня щит (защита: 4, надета: false)"}, nil, "щит|4|3|false\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for index, want := range tt.messages {
				var wantError error
				if index == len(tt.messages)-1 {
					wantError = tt.lastError
				}
				got, err := tt.item.Use()
				if got != want || !errors.Is(err, wantError) {
					t.Errorf("Use() = (%q, %v), want (%q, %v)", got, err, want, wantError)
				}
			}
			checkRecord(t, tt.item, tt.record)
		})
	}
}

func TestValidate(t *testing.T) {
	for _, tt := range []struct {
		name  string
		item  Storable
		valid bool
	}{
		{"zero values", NewWeapon("меч", 0, 0, 0), true},
		{"max name length", NewArmor(strings.Repeat("a", 16*1024), 0, 0), true},
		{"name too long", NewArmor(strings.Repeat("a", 16*1024+1), 0, 0), false},
		{"blank name", NewWeapon(" ", 1, 1, 1), false},
		{"separator in name", NewWeapon("a|b", 1, 1, 1), false},
		{"negative weight", NewWeapon("меч", -1, 1, 1), false},
		{"negative damage", NewWeapon("меч", 1, -1, 1), false},
		{"negative durability", NewWeapon("меч", 1, 1, -1), false},
		{"negative defense", NewArmor("щит", 1, -1), false},
		{"negative charges", NewPotion("зелье", 1, "HP", -1), false},
		{"invalid effect", NewPotion("зелье", 1, "HP\n", 0), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Validate(); (err == nil) != tt.valid {
				t.Errorf("Validate() = %v, want valid = %t", err, tt.valid)
			}
		})
	}
}

func TestDeserializeSequential(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("меч|5|1|2\r\nщит|4|3|true\nзелье|HP|0|1"))
	for _, tt := range []struct {
		item   Storable
		record string
	}{
		{&Weapon{}, "меч|5|1|2\n"},
		{&Armor{}, "щит|4|3|true\n"},
		{&Potion{}, "зелье|HP|0|1\n"},
	} {
		if err := tt.item.Deserialize(reader); err != nil {
			t.Fatalf("Deserialize(%T): %v", tt.item, err)
		}
		checkRecord(t, tt.item, tt.record)
	}
	if err := (&Weapon{}).Deserialize(reader); !errors.Is(err, io.EOF) {
		t.Errorf("Deserialize() after last record = %v, want EOF", err)
	}
}

func TestDeserializeInvalidRecord(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  string
		item   Storable
		before string
	}{
		{"wrong weapon field count", "меч|5|1", NewWeapon("лук", 1, 2, 3), "лук|2|3|1\n"},
		{"invalid damage", "меч|x|1|2", NewWeapon("лук", 1, 2, 3), "лук|2|3|1\n"},
		{"invalid durability", "меч|5|x|2", NewWeapon("лук", 1, 2, 3), "лук|2|3|1\n"},
		{"invalid weapon weight", "меч|5|1|x", NewWeapon("лук", 1, 2, 3), "лук|2|3|1\n"},
		{"infinite weight", "меч|5|1|+Inf", NewWeapon("лук", 1, 2, 3), "лук|2|3|1\n"},
		{"wrong armor field count", "щит|4|3", NewArmor("броня", 2, 6), "броня|6|2|false\n"},
		{"invalid defense", "щит|x|3|false", NewArmor("броня", 2, 6), "броня|6|2|false\n"},
		{"invalid armor weight", "щит|4|x|false", NewArmor("броня", 2, 6), "броня|6|2|false\n"},
		{"negative defense", "щит|-1|3|false", NewArmor("броня", 2, 6), "броня|6|2|false\n"},
		{"invalid armor state", "щит|4|3|yes", NewArmor("броня", 2, 6), "броня|6|2|false\n"},
		{"wrong potion field count", "зелье|HP|1", NewPotion("старое", 2, "MP", 3), "старое|MP|3|2\n"},
		{"invalid potion charges", "зелье|HP|x|1", NewPotion("старое", 2, "MP", 3), "старое|MP|3|2\n"},
		{"invalid potion weight", "зелье|HP|1|x", NewPotion("старое", 2, "MP", 3), "старое|MP|3|2\n"},
		{"negative charges", "зелье|HP|-1|1", NewPotion("старое", 2, "MP", 3), "старое|MP|3|2\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Deserialize(bufio.NewReader(strings.NewReader(tt.input))); err == nil {
				t.Errorf("Deserialize(%q) succeeded", tt.input)
			}
			checkRecord(t, tt.item, tt.before)
		})
	}
}

type failingReader struct {
	reader io.Reader
	err    error
}

func (r failingReader) Read(data []byte) (int, error) {
	n, err := r.reader.Read(data)
	if err == io.EOF {
		return n, r.err
	}
	return n, err
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestItemStreamFailures(t *testing.T) {
	for _, tt := range []struct {
		name   string
		value  Storable
		record string
	}{
		{"weapon", NewWeapon("меч", 2, 5, 1), "меч|5|1|2\n"},
		{"armor", NewArmor("щит", 3, 4), "щит|4|3|false\n"},
		{"potion", NewPotion("зелье", 1, "HP", 2), "зелье|HP|2|1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			failure := errors.New("stream failed")
			reader := bufio.NewReader(failingReader{reader: strings.NewReader("partial record"), err: failure})
			if err := tt.value.Deserialize(reader); !errors.Is(err, failure) {
				t.Errorf("Deserialize() = %v, want %v", err, failure)
			}
			checkRecord(t, tt.value, tt.record)
			if err := tt.value.Serialize(failingWriter{err: failure}); !errors.Is(err, failure) {
				t.Errorf("Serialize() = %v, want %v", err, failure)
			}
			checkRecord(t, tt.value, tt.record)
		})
	}
}

func TestReadFieldsBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name   string
		reader io.Reader
		want   string
	}{
		{"unbuffered reader", strings.NewReader("name|1|2|3"), "Передайте bufio.NewReader"},
		{"record at size limit", bufio.NewReader(strings.NewReader(strings.Repeat("a", 64*1024))), "слишком длинная запись"},
		{"carriage return inside record", bufio.NewReader(strings.NewReader("na\rme|1|2|3\n")), "ожидается одна запись"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := readFields(tt.reader)
			if fields != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("readFields() = (%v, %v), want no fields and error containing %q", fields, err, tt.want)
			}
		})
	}
	const size = 64*1024 - 1
	fields, err := readFields(bufio.NewReader(strings.NewReader(strings.Repeat("a", size))))
	if err != nil || len(fields) != 1 || len(fields[0]) != size {
		t.Errorf("readFields() below size limit returned %d fields and error %v", len(fields), err)
	}
}

func TestDescribeItem(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value Item
		want  string
	}{
		{"weapon", NewWeapon("меч", 2, 5, 1), "меч (вес: 2)"},
		{"armor", NewArmor("щит", 3, 4), "щит (вес: 3)"},
		{"potion", NewPotion("зелье", 1, "HP", 1), "зелье (вес: 1)"},
		{"nil interface", nil, "Предмет отсутствует"},
		{"nil pointer", (*Weapon)(nil), "Предмет отсутствует"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescribeItem(tt.value); got != tt.want {
				t.Errorf("DescribeItem() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUseDangerousItem(t *testing.T) {
	box := NewWeapon("Ящик Пандоры", 1, 0, 1)
	weapon := NewWeapon("меч", 1, 5, 1)
	message, err := UseDangerousItem(weapon)
	if err != nil || message != "Используется меч (урон: 5, прочность: 0)" || weapon.Durability != 0 {
		t.Fatalf("UseDangerousItem(weapon) = (%q, %v), durability = %d", message, err, weapon.Durability)
	}
	if _, err := UseDangerousItem(weapon); !errors.Is(err, ErrBroken) {
		t.Errorf("UseDangerousItem(broken weapon) = %v, want %v", err, ErrBroken)
	}
	if _, err := UseDangerousItem(box); err == nil || err.Error() != "опасный предмет: ящик Пандоры открыт" || box.Durability != 1 {
		t.Errorf("UseDangerousItem(box) error = %v, durability = %d", err, box.Durability)
	}
	if _, err := UseDangerousItem(nil); err == nil {
		t.Error("UseDangerousItem(nil) succeeded")
	}
}

func checkRecord(t *testing.T, value Storable, want string) {
	t.Helper()
	var out bytes.Buffer
	if err := value.Serialize(&out); err != nil {
		t.Fatalf("Serialize(%T): %v", value, err)
	}
	if got := out.String(); got != want {
		t.Errorf("Serialize(%T) = %q, want %q", value, got, want)
	}
}
