package entity

import (
	"errors"
	"time"
)

var (
	ErrTaken   = errors.New("книга уже выдана")
	ErrTooMany = errors.New("на руках уже 3 книги")
	ErrOverdue = errors.New("есть просроченная книга")
)

type Loan struct {
	BookID   int
	ReaderID int
	Due      time.Time
}
