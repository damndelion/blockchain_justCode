package userEntity

import "gorm.io/gorm"

type UserCredentials struct {
	gorm.Model
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	CardNum string `json:"card_num"`
	Type    string `json:"type"`
	CVV     string `json:"cvv"`
}
