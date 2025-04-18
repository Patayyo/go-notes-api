package middleware

import (
	"context"
	"errors"
	"net/http"
	"notes-api/logger"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Warn("Отсутствует или некорректный заголовок Authorization")
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default_secret"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Log.Error("Недопустимый метод подписи")
				return nil, errors.New("invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.Log.WithError(err).Warn("Недопустимый или просроченный токен")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			logger.Log.Warn("Отсутствует user_id в токене")
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		userID := uint(claims["user_id"].(float64))
		logger.Log.WithField("user_id", userID).Info("Аутентификация прошла успешно")
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
