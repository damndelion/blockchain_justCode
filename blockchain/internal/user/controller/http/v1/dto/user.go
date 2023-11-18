package dto

type UserDetailRequest struct {
	Age      int    `json:"age" binding:"required,min=1,max=120"`
	Phone    string `json:"phone" binding:"required,e164"`
	Address  string `json:"address" binding:"required"`
	Country  string `json:"country" binding:"required"`
	City     string `json:"city" binding:"required"`
	CardNum  string `json:"card_num" binding:"required,len=16"`
	CardType string `json:"card_type" binding:"required"`
	CVV      string `json:"cvv" binding:"required,len=3"`
}

type UserInfoRequest struct {
	UserID  string `json:"userID" binding:"required"`
	Age     int    `json:"age" binding:"required,min=1,max=120"`
	Phone   string `json:"phone" binding:"required,e164"`
	Address string `json:"address" binding:"required"`
	Country string `json:"country" binding:"required"`
	City    string `json:"city" binding:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Wallet   string `json:"wallet" binding:"required"`
	Valid    bool   `json:"valid,omitempty" default:"false"`
}

type UserCredRequest struct {
	UserID   string `json:"userID" binding:"required"`
	CardNum  string `json:"card_num" binding:"required,len=16"`
	CardType string `json:"type" binding:"required"`
	CVV      string `json:"cvv" binding:"required,len=3"`
}

type TopUpRequest struct {
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}
