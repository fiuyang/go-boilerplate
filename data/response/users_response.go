package response

type UsersResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}
