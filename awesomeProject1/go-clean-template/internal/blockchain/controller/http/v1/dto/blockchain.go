package dto

type SendRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type TopUpRequest struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
