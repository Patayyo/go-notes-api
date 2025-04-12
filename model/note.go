package model

// Note представляет заметку пользователя
// @Description Модель заметки
type Note struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
