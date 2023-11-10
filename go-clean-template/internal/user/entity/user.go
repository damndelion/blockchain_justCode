package userEntity

// todo isActive Role
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
	Valid    bool   `json:"valid"`
}
