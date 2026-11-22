package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"clinic/homework/fat/entity"
)

type fakeRepo struct {
	taken   bool
	onHands int
	saved   int
}

func (f *fakeRepo) IsTaken(_ context.Context, _ int) (bool, error) {
	return f.taken, nil
}

func (f *fakeRepo) CountByReader(_ context.Context, _ int) (int, error) {
	return f.onHands, nil
}

func (f *fakeRepo) HasLoanDueBefore(_ context.Context, _ int, _ time.Time) (bool, error) {
	return false, nil
}

func (f *fakeRepo) Save(_ context.Context, _ entity.Loan) error {
	f.saved++
	return nil
}

func TestLend(t *testing.T) {
	tests := []struct {
		name      string
		taken     bool
		onHands   int
		wantErr   error
		wantSaved int
	}{
		{"выдача", false, 0, nil, 1},
		{"книга выдана", true, 0, entity.ErrTaken, 0},
		{"лимит читателя", false, 3, entity.ErrTooMany, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{taken: tt.taken, onHands: tt.onHands}
			library := New(repo)
			before := time.Now().AddDate(0, 0, 14)

			due, err := library.Lend(context.Background(), 1, 2)

			after := time.Now().AddDate(0, 0, 14)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if repo.saved != tt.wantSaved {
				t.Fatalf("saved = %d, want %d", repo.saved, tt.wantSaved)
			}
			if tt.wantErr == nil && (due.Before(before) || due.After(after)) {
				t.Fatalf("due = %v, want between %v and %v", due, before, after)
			}
		})
	}
}
