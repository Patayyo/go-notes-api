package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/model"
	"notes-api/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	Store service.INoteService
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Store.GetAllNotes()
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
	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	created, err := h.Store.CreateNote(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	updated, err := h.Store.UpdateNote(id, note)
	if err != nil {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	if err := h.Store.DeleteNote(id); err != nil {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Заметка удалена"})
}
