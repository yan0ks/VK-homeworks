package session2

import "testing"

// Тесты уже готовы. Задача - сделать их зелёными: go test ./...
// Запустить одну функцию: go test -run TestHasDuplicate ./...
// В каждую таблицу можно и нужно дописывать свои случаи.

func TestHasDuplicate(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want bool
	}{
		{"есть повтор", []int{3, 1, 3}, true},
		{"нет повтора", []int{1, 2, 3}, false},
		{"пустой вход", nil, false},
		{"один элемент", []int{5}, false},
		{"повтор в самом конце", []int{1, 2, 3, 4, 1}, true},
		{"все одинаковые", []int{7, 7, 7}, true},
	}
	for _, c := range cases {
		if got := HasDuplicate(c.in); got != c.want {
			t.Errorf("%s: HasDuplicate(%v) = %v, ждали %v", c.name, c.in, got, c.want)
		}
	}
}

func TestTwoItems(t *testing.T) {
	cases := []struct {
		name   string
		prices []int
		k      int
		i, j   int
		wantOK bool
	}{
		{"пара в середине", []int{700, 1500, 300, 1200}, 1500, 2, 3, true},
		{"ответ на позиции 0", []int{700, 800}, 1500, 0, 1, true},
		{"один товар дважды брать нельзя", []int{700, 500}, 1400, 0, 0, false},
		{"пустой вход", nil, 5, 0, 0, false},
		{"один товар", []int{500}, 1000, 0, 0, false},
		{"одинаковые цены", []int{750, 750}, 1500, 0, 1, true},
	}
	for _, c := range cases {
		i, j, ok := TwoItems(c.prices, c.k)
		if ok != c.wantOK || (ok && (i != c.i || j != c.j)) {
			t.Errorf("%s: TwoItems(%v, %d) = %d, %d, %v, ждали %d, %d, %v",
				c.name, c.prices, c.k, i, j, ok, c.i, c.j, c.wantOK)
		}
	}
}

func TestBracketsOK(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"([{}])", true},
		{"([)]", false}, // счётчик бы пропустил, стек - нет
		{"(]", false},
		{"((", false},
		{"", true},
		{")", false}, // снимать со стека нечего
		{"a(b)c", true},
		{"()[]{}", true},
	}
	for _, c := range cases {
		if got := BracketsOK(c.in); got != c.want {
			t.Errorf("BracketsOK(%q) = %v, ждали %v", c.in, got, c.want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"шалаш", true},
		{"шалаша", false},
		{"", true},
		{"я", true},
		{"abba", true},
		{"abca", false},
		{"топот", true},
	}
	for _, c := range cases {
		if got := IsPalindrome(c.in); got != c.want {
			t.Errorf("IsPalindrome(%q) = %v, ждали %v", c.in, got, c.want)
		}
	}
}

func TestIsAnagram(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"кот", "ток", true},
		{"кот", "кит", false},
		{"ааб", "абб", false}, // буквы те же, а количество разное
		{"", "", true},
		{"кот", "ко", false},
		{"listen", "silent", true},
	}
	funcs := []struct {
		name string
		f    func(a, b string) bool
	}{
		{"IsAnagramSort", IsAnagramSort},
		{"IsAnagramMap", IsAnagramMap},
	}
	for _, fn := range funcs {
		for _, c := range cases {
			if got := fn.f(c.a, c.b); got != c.want {
				t.Errorf("%s(%q, %q) = %v, ждали %v", fn.name, c.a, c.b, got, c.want)
			}
		}
	}
}

// По желанию.
func TestCanFulfill(t *testing.T) {
	stock := []string{"мыло", "мыло", "щётка"}
	cases := []struct {
		name  string
		stock []string
		order []string
		want  bool
	}{
		{"всё есть", stock, []string{"мыло", "щётка"}, true},
		{"щётка одна, а просят две", stock, []string{"щётка", "щётка"}, false},
		{"пустой заказ", stock, nil, true},
		{"пустой склад", nil, []string{"мыло"}, false},
		{"забрали всё", stock, []string{"мыло", "щётка", "мыло"}, true},
	}
	for _, c := range cases {
		if got := CanFulfill(c.stock, c.order); got != c.want {
			t.Errorf("%s: CanFulfill(%v, %v) = %v, ждали %v", c.name, c.stock, c.order, got, c.want)
		}
	}
}

// Со звёздочкой.
func TestMinStack(t *testing.T) {
	var s MinStack
	steps := []struct {
		op      string // "push" или "pop"
		x       int    // что кладём (для push) или что ждём от Pop (для pop)
		wantMin int    // Minimum() после шага
	}{
		{"push", 5, 5},
		{"push", 3, 3},
		{"push", 7, 3},
		{"push", 3, 3},
		{"pop", 3, 3},
		{"pop", 7, 3},
		{"pop", 3, 5},
		{"push", 1, 1},
		{"pop", 1, 5},
	}
	for n, st := range steps {
		if st.op == "push" {
			s.Push(st.x)
		} else if got := s.Pop(); got != st.x {
			t.Fatalf("шаг %d: Pop() = %d, ждали %d", n, got, st.x)
		}
		if got := s.Minimum(); got != st.wantMin {
			t.Fatalf("шаг %d (%s %d): Minimum() = %d, ждали %d", n, st.op, st.x, got, st.wantMin)
		}
	}
}
