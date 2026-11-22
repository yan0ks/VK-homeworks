package repository

import (
	"context"
	"sync"
	"time"

	"clinic/homework/fat/entity"
)

type Memory struct {
	mu          sync.Mutex
	loansByBook map[int]entity.Loan
}

func NewMemory() *Memory {
	return &Memory{loansByBook: make(map[int]entity.Loan)}
}

func (m *Memory) IsTaken(_ context.Context, bookID int) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, taken := m.loansByBook[bookID]
	return taken, nil
}

func (m *Memory) CountByReader(_ context.Context, readerID int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var count int
	for _, loan := range m.loansByBook {
		if loan.ReaderID == readerID {
			count++
		}
	}
	return count, nil
}

func (m *Memory) HasLoanDueBefore(_ context.Context, readerID int, cutoff time.Time) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, loan := range m.loansByBook {
		if loan.ReaderID == readerID && loan.Due.Before(cutoff) {
			return true, nil
		}
	}
	return false, nil
}

func (m *Memory) Save(_ context.Context, loan entity.Loan) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loansByBook[loan.BookID] = loan
	return nil
}
