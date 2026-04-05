package httpapi

type UserInfo struct {
	Birthday  string `json:"birthday"`
	Country   string `json:"country"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Gender    string `json:"gender"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Picture   string `json:"picture"`
	Username  string `json:"username"`
}
