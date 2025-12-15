package param


type UserInfo struct{
	ID uint `json:"id"`
	Name string `json:"name"`
	PhoneNumber string  `json:"phone_number"`
}

type RegisterRequest struct {
	Name string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}


type RegisterResponse struct {
	User UserInfo `json:"user"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}



type LoginResponse struct {
	User UserInfo
	Tokens Tokens
}