package models

type Token struct {
	UserEmail string `json:"userEmail"`
	UserName  string `json:"userName"`
}

type TokenResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
