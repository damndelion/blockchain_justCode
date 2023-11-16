package entity

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	UserEmail    string `json:"user_email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserCode struct {
	gorm.Model
	Email string `json:"email"`
	Code  string `json:"code"`
}
