package main

import (
	"os"

	"homeworks/homework-4/task-1/internal/cli"
	"homeworks/homework-4/task-1/internal/storage"
)

func main() {
	// Запускать из homework-4/task-1. Служебные данные остаются в task-99.
	const dataDir = "../../task-99/homework-4/task-1"
	store := storage.NewFile(dataDir + "/inventory.txt")
	app := cli.New(store)
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))
}

/*
Вопрос: Почему ошибка сломанного предмета не должна быть паникой?

Ответ: Это ожидаемая ситуация, которую вызывающий код может обработать.
Метод возвращает ErrBroken, а вызывающий код проверяет её через errors.Is.
*/

/*
Вопрос: Какая строка после panic не будет выполнена?

Ответ: В UseDangerousItem после panic не выполнится return value.Use().
Паника прерывает обычное выполнение функции, но отложенные вызовы выполняются.
*/

/*
Вопрос: Почему recover должен находиться в отложенной функции?

Ответ: Во время паники обычное выполнение уже прервано.
Отложенная функция запускается при выходе из функции и может остановить панику,
если вызывает recover непосредственно внутри себя.
*/
