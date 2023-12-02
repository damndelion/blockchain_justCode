package userentity

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
	Valid    bool   `json:"valid"`
	Role     string `json:"role"`
}
