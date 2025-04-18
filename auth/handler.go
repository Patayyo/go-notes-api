package auth

import (
	"encoding/json"
	"net/http"
	"notes-api/logger"
	"notes-api/middleware"
)

type AuthHandler struct {
	Service *AuthService
}

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth.AuthInput true "Данные пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Неверный запрос"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input AuthInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.WithError(err).Warn("Ошибка декодирования запроса на регистрацию")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.Service.Register(input.Email, input.Password); err != nil {
		logger.Log.WithFields(logger.Fields{
			"email": input.Email,
		}).WithError(err).Warn("Ошибка регистрации")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Log.WithFields(logger.Fields{
		"email": input.Email,
	}).Info("Пользователь успешно зарегистрирован")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	})
}

// Login godoc
// @Summary Авторизация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth.AuthInput true "Данные пользователя"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Неверные данные"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input AuthInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.WithError(err).Warn("Ошибка декодирования запроса на вход")
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"email": input.Email,
		}).WithError(err).Warn("Ошибка входа")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logger.Log.WithFields(logger.Fields{
		"email": input.Email,
	}).Info("Успешный вход пользователя")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh godoc
// @Summary Обновить access token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Refresh токен"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Ошибка токена"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.WithError(err).Warn("Ошибка декодирования refresh токена")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	newAccessToken, err := h.Service.RefreshToken(input.RefreshToken)
	if err != nil {
		logger.Log.WithError(err).Warn("Ошибка при обновлении токена")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Log.Info("Access token успешно обновлён")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := rawUserID.(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	if err := h.Service.Logout(userID); err != nil {
		http.Error(w, "Ошибка при выходе", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Вы успешно вышли",
	})
}
