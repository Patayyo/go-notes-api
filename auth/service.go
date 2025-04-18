package auth

import (
	"errors"
	"notes-api/logger"
	"notes-api/model"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(email, password string) error {
	logger.Log.Infof("Попытка регистрации: %s", email)
	if err := validateCredentials(email, password); err != nil {
		return err
	}

	var existing model.User
	if err := s.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("пользователь с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Email: email,
		Hash:  string(hash),
	}

	return s.DB.Create(&user).Error
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	logger.Log.Infof("Попытка входа: %s", email)
	var user model.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return "", "", errors.New("неверный email или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	user.RefreshTokenHash = refreshTokenString
	if err := s.DB.Save(&user).Error; err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неподдерживаемый алгоритм подписи")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("некорректный или просроченный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return "", errors.New("невалидный payload токена")
	}

	userID := uint(claims["user_id"].(float64))

	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return "", errors.New("пользователь не найден")
	}

	if user.RefreshTokenHash != refreshTokenString {
		return "", errors.New("refresh токен недействителен")
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	})

	newTokenString, err := newToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

func (s *AuthService) Logout(userID uint) error {
	return s.DB.Model(&model.User{}).Where("id = ?", userID).Update("refresh_token_hash", "").Error
}

func validateCredentials(email, password string) error {
	if email == "" || password == "" {
		return errors.New("email и пароль обязательны")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		return errors.New("неверный формат email")
	}

	if len(password) < 6 {
		return errors.New("пароль должен содержать не менее 6 символов")
	}

	if len([]byte(password)) > 72 {
		return errors.New("пароль не должен превышать 72 байта")
	}

	return nil
}
