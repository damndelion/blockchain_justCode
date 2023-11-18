package dto

type SendRequest struct {
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type AddressRequest struct {
	Address string `json:"address" binding:"required"`
}
