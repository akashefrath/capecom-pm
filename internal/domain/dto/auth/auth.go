package authdto

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" default:"Bearer"`
	// ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}
