package model

// User представляет пользователя
// @Description Модель пользователя
type User struct {
	ID               uint   `json:"id" gorm:"primarykey"`
	Email            string `json:"email" gorm:"unique"`
	Password         string `json:"password"`
	Hash             string `json:"hash"`
	RefreshToken     string `json:"refresh_token" gorm:"column:refresh_token"`
	RefreshTokenHash string `json:"refresh_token_hash" gorm:"column:refresh_token_hash"`
}
