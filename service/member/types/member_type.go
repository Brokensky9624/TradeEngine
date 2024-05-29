package types

type Member struct {
	ID       uint   `json:"id"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
