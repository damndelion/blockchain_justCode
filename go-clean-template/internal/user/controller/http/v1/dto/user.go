package dto

type UserDetailRequest struct {
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Country  string `json:"country"`
	City     string `json:"city"`
	CardNum  string `json:"card_num"`
	CardType string `json:"card_type"`
	CVV      string `json:"cvv"`
}

type TopUpRequest struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
