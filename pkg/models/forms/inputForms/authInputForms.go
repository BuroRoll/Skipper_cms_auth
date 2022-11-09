package inputForms

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenReqBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password string `json:"password" binding:"required"`
}
