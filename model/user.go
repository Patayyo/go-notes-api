package model

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Hash     string `json:"hash"`
}
