package outputForms

import "Skipper_cms_auth/pkg/models"

type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
