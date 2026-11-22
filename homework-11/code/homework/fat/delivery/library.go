package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"clinic/homework/fat/entity"
)

type Usecase interface {
	Lend(ctx context.Context, bookID, readerID int) (time.Time, error)
}

type Handler struct {
	uc Usecase
}

func New(uc Usecase) *Handler {
	return &Handler{uc: uc}
}

type lendRequest struct {
	BookID   int `json:"book_id"`
	ReaderID int `json:"reader_id"`
}

type lendResponse struct {
	Due string `json:"due"`
}

func (h *Handler) Lend(w http.ResponseWriter, r *http.Request) {
	var req lendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookID == 0 || req.ReaderID == 0 {
		http.Error(w, "плохой запрос", http.StatusBadRequest)
		return
	}
	due, err := h.uc.Lend(r.Context(), req.BookID, req.ReaderID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrTaken):
			http.Error(w, entity.ErrTaken.Error(), http.StatusConflict)
		case errors.Is(err, entity.ErrTooMany):
			http.Error(w, entity.ErrTooMany.Error(), http.StatusConflict)
		case errors.Is(err, entity.ErrOverdue):
			http.Error(w, entity.ErrOverdue.Error(), http.StatusConflict)
		default:
			log.Printf("lend: %v", err)
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(lendResponse{Due: due.Format("2006-01-02")}); err != nil {
		log.Printf("write response: %v", err)
	}
}
