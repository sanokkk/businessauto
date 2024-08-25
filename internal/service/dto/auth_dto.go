package dto

type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type GetUserResponse struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}
