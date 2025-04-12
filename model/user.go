package model

// User представляет пользователя
// @Description Модель пользователя
type User struct {
	ID           uint   `json:"id" gorm:"primarykey"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"password"`
	Hash         string `json:"hash"`
	RefreshToken string `json:"-" gorm:"column:refresh_token"`
}
