package userentity

type UserInfo struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Age     int    `json:"age"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Country string `json:"country"`
	City    string `json:"city"`
}
