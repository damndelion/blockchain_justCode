package userentity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
	Valid    bool   `json:"valid"`
	Role     string `json:"role"`
}
