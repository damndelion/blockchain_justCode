package entity

type UserCredentials struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	CardNum string `json:"card_num"`
	Type    string `json:"type"`
	CVV     string `json:"cvv"`
}
