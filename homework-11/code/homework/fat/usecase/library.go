package usecase

import (
	"context"
	"fmt"
	"time"

	"clinic/homework/fat/entity"
)

type Repository interface {
	IsTaken(ctx context.Context, bookID int) (bool, error)
	CountByReader(ctx context.Context, readerID int) (int, error)
	HasLoanDueBefore(ctx context.Context, readerID int, cutoff time.Time) (bool, error)
	Save(ctx context.Context, loan entity.Loan) error
}

type Library struct {
	repo Repository
}

func New(repo Repository) *Library {
	return &Library{repo: repo}
}

func (u *Library) Lend(ctx context.Context, bookID, readerID int) (time.Time, error) {
	taken, err := u.repo.IsTaken(ctx, bookID)
	if err != nil {
		return time.Time{}, fmt.Errorf("is taken: %w", err)
	}
	if taken {
		return time.Time{}, entity.ErrTaken
	}
	onHands, err := u.repo.CountByReader(ctx, readerID)
	if err != nil {
		return time.Time{}, fmt.Errorf("count by reader: %w", err)
	}
	if onHands >= 3 {
		return time.Time{}, entity.ErrTooMany
	}
	now := time.Now()
	overdue, err := u.repo.HasLoanDueBefore(ctx, readerID, now)
	if err != nil {
		return time.Time{}, fmt.Errorf("check overdue loans: %w", err)
	}
	if overdue {
		return time.Time{}, entity.ErrOverdue
	}
	due := now.AddDate(0, 0, 14)
	loan := entity.Loan{BookID: bookID, ReaderID: readerID, Due: due}
	if err := u.repo.Save(ctx, loan); err != nil {
		return time.Time{}, fmt.Errorf("save: %w", err)
	}
	return due, nil
}
