package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type streamHandler struct {
	store StreamStore
}

func newStreamHandler(store StreamStore) *streamHandler {
	return &streamHandler{store: store}
}

func (s *streamHandler) GetStream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing buff id in path", http.StatusBadRequest)
		return
	}

	parsedid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id, %v", err), http.StatusBadRequest)
		return
	}

	stream, err := s.store.GetStream(parsedid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, fmt.Sprintf("not found"), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to get stream, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stream); err != nil {
		http.Error(w, fmt.Sprintf("failed encode stream, %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *streamHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
	var req CreateStream
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal request, %v", err), http.StatusBadRequest)
		return
	}
	id, err := s.store.SetStream(Stream{
		ID:   0,
		Name: req.Name,
		BuffIDs: req.BuffIDs,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create buff, %v", err), http.StatusInternalServerError)
		return
	}

	stream, err := s.store.GetStream(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get buff, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stream); err != nil {
		http.Error(w, fmt.Sprintf("failed encode buff, %v", err), http.StatusInternalServerError)
		return
	}
}
