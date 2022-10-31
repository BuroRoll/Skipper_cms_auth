package inputForms

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenReqBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
