package userentity

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Age     int    `json:"age"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Country string `json:"country"`
	City    string `json:"city"`
}
