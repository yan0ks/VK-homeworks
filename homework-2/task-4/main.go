package main

import "fmt"

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if p == nil {
		return
	}
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Age != nil {
		p.Age = *patch.Age
	}
	if patch.Email != nil {
		p.Email = *patch.Email
	}
	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}

func profileExample() {
	p := Profile{Name: "Анна", Age: 30, Email: "anna@example.com", Newsletter: true}
	age, newsletter := 31, false
	ApplyPatch(&p, ProfilePatch{Age: &age, Newsletter: &newsletter})
	fmt.Printf("%+v\n", p)
	ApplyPatch(&p, ProfilePatch{})
	fmt.Printf("После пустого патча: %+v\n", p)
	ApplyPatch(nil, ProfilePatch{Age: &age})
}

func main() {
	profileExample()
}

/*
Вопрос: почему в патче нельзя просто объявить Newsletter bool?

Ответ: bool не позволяет отличить «не менять поле» от «записать false».
У *bool есть три состояния: nil — не менять, указатель на true — включить
рассылку, указатель на false — выключить.
*/
