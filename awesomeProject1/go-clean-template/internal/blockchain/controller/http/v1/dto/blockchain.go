package dto

type SendRequest struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount"`
}

type AddressRequest struct {
	Address string `json:"address"`
}
