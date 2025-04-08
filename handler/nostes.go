package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/middleware"
	"notes-api/model"
	"notes-api/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	Store service.INoteService
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	notes, err := h.Store.GetNotesByUserID(int(userID))
	if err != nil {
		http.Error(w, "Ошибка при получении заметок", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	note, err := h.Store.GetNoteByID(id)
	if err != nil {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	created, err := h.Store.CreateNote(userID, note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(middleware.UserIDKey)
	userID, ok := raw.(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var note model.Note
	note, err = h.Store.GetNoteByID(id)
	if err != nil {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	var updated model.Note
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	updatedNote, err := h.Store.UpdateNote(id, updated)
	if err != nil {
		http.Error(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedNote)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	note, err := h.Store.GetNoteByID(id)
	if err != nil {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	if err := h.Store.DeleteNote(id); err != nil {
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Заметка удалена"})
}
