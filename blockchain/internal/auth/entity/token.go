package entity

type Token struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type UserCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
