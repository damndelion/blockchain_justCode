package entity

type Token struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	UserEmail    string `json:"user_email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserVerifications struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
