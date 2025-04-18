package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/logger"
	"notes-api/middleware"
	"notes-api/model"
	"notes-api/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	Store service.INoteService
}

// GetAll godoc
// @Summary Получить все заметки текущего пользователя
// @Tags notes
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} model.Note "Список заметок"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 500 {string} string "Ошибка при получении заметок"
// @Router /notes [get]
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		logger.Log.Warn("Пользователь не аутентифицирован")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	notes, err := h.Store.GetNotesByUserID(int(userID))
	if err != nil {
		logger.Log.WithError(err).WithField("user_id", userID).Error("Ошибка при получении заметок")
		http.Error(w, "Ошибка при получении заметок", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notes)

	logger.Log.WithField("user_id", userID).Info("Получение всех заметок пользователя")
}

// GetByID godoc
// @Summary Получить заметку по ID
// @Tags notes
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID заметки"
// @Success 200 {object} model.Note "Заметка"
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Заметка не найдена"
// @Router /notes/{id} [get]
func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithError(err).Warn("Неверный ID")
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	note, err := h.Store.GetNoteByID(id)
	if err != nil {
		logger.Log.WithError(err).Warn("Заметка не найдена")
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)

	logger.Log.WithFields(logger.Fields{
		"note_id": id,
	}).Info("Заметка найдена и возвращена")
}

// Create godoc
// @Summary Создать новую заметку
// @Tags notes
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param note body model.Note true "Данные заметки"
// @Success 201 {object} model.Note "Созданная заметка"
// @Failure 400 {string} string "Неверный запрос или ошибка валидации"
// @Router /notes [post]
func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		logger.Log.Warn("Пользователь не аутентифицирован")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		logger.Log.Warn("Неверный запрос")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	created, err := h.Store.CreateNote(userID, note)
	if err != nil {
		logger.Log.WithError(err).WithFields(logger.Fields{
			"user_id": userID,
			"note":    note,
		}).Warn("Ошибка при создании заметки")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)

	logger.Log.WithFields(logger.Fields{
		"user_id": userID,
		"note_id": created.ID,
	}).Info("Заметка успешно создана")
}

// Update godoc
// @Summary Обновить заметку по ID (только владелец может обновить)
// @Tags notes
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Param note body model.Note true "Обновлённые данные заметки"
// @Success 200 {object} model.Note "Обновленная заметка"
// @Failure 400 {string} string "Неверный запрос или ID"
// @Failure 403 {string} string "Доступ запрещён"
// @Failure 404 {string} string "Заметка не найдена"
// @Router /notes/{id} [put]
func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(middleware.UserIDKey)
	userID, ok := raw.(uint)
	if !ok {
		logger.Log.Warn("Пользователь не авторизован")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithError(err).Warn("Неверный ID")
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var note model.Note
	note, err = h.Store.GetNoteByID(id)
	if err != nil {
		logger.Log.WithError(err).Warn("Заметка не найдена")
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		logger.Log.Warn("Доступ запрещён")
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	var updated model.Note
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		logger.Log.WithError(err).Warn("Неверный запрос")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	updatedNote, err := h.Store.UpdateNote(id, updated)
	if err != nil {
		logger.Log.WithError(err).WithFields(logger.Fields{
			"id":   id,
			"note": updated,
		}).Warn("Ошибка при обновлении")
		http.Error(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedNote)

	logger.Log.WithFields(logger.Fields{
		"user_id": userID,
		"note_id": id,
	}).Info("Заметка успешно обновлена")
}

// Delete godoc
// @Summary Удалить заметку по ID (только владелец может удалить)
// @Tags notes
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID заметки"
// @Success 200 {object} map[string]string "Сообщение об удалении"
// @Failure 400 {string} string "Неверный ID"
// @Failure 403 {string} string "Доступ запрещён"
// @Failure 404 {string} string "Заметка не найдена"
// @Router /notes/{id} [delete]
func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		logger.Log.Warn("Пользователь не аутентифицирован")
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithError(err).Warn("Неверный ID")
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	note, err := h.Store.GetNoteByID(id)
	if err != nil {
		logger.Log.WithError(err).Warn("Заметка не найдена")
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		logger.Log.Warn("Доступ запрещён")
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	if err := h.Store.DeleteNote(id); err != nil {
		logger.Log.WithError(err).WithField("id", id).Warn("Ошибка при удалении")
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Заметка удалена"})

	logger.Log.WithFields(logger.Fields{
		"user_id": userID,
		"note_id": id,
	}).Info("Заметка успешно удалена")
}
