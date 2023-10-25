package entity

type Token struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}
