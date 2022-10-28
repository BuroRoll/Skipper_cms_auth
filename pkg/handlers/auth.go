package handlers

import (
	"Skipper_cms_auth/pkg/models/forms/inputForms"
	"Skipper_cms_auth/pkg/models/forms/outputForms"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Description Регистрация нового пользователя
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param 		request 	body 		inputForms.SignUpUserForm 	true 	"query params"
// @Success 	200 		{object} 	outputForms.AuthResponse
// @Failure     400         {object}  	outputForms.ErrorResponse
// @Failure     500         {object}  	outputForms.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input inputForms.SignUpUserForm
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, outputForms.ErrorResponse{
			Error: "Неверная форма регистрации",
		})
		return
	}
	userId, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputForms.ErrorResponse{
			Error: "Ошибка создание профиля",
		})
		return
	}
	token, refreshToken, err := h.services.Authorization.GenerateToken(input.Phone, input.Password)
	user, err := h.services.GetUserData(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputForms.ErrorResponse{
			Error: "Ошибка генерации токена",
		})
		return
	}
	c.JSON(http.StatusOK, outputForms.AuthResponse{
		RefreshToken: refreshToken,
		Token:        token,
		User:         user,
	},
	)
}

// @Description 	Вход
// @Tags 			Auth
// @Accept 			json
// @Produce 		json
// @Param 			request 	body 		inputForms.SignInInput 	true 	"query params"
// @Success 		200 		{object} 	outputForms.AuthResponse
// @Failure     	400         {object}  	outputForms.ErrorResponse
// @Failure     	500         {object}  	outputForms.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input inputForms.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, outputForms.ErrorResponse{
			Error: "Неверная форма авторизации",
		})
		return
	}
	token, refreshToken, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
	userId, _, err := h.services.GetUser(input.Login, input.Password)
	user, err := h.services.GetUserData(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputForms.ErrorResponse{
			Error: "Неверный логин или пароль",
		})
		return
	}
	c.JSON(http.StatusOK, outputForms.AuthResponse{
		RefreshToken: refreshToken,
		Token:        token,
		User:         user,
	},
	)
}

// @Description Обновление токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param 			request 	body 		inputForms.TokenReqBody 	true 	"query params"
// @Success 		200 		{object} 	outputForms.RefreshTokenResponse
// @Failure     	400         {object}  	outputForms.ErrorResponse
// @Failure     	500         {object}  	outputForms.ErrorResponse
// @Router /auth/refresh-token [post]
func (h *Handler) refreshToken(c *gin.Context) {
	var input inputForms.TokenReqBody
	err := c.Bind(&input)
	userId, err := h.services.ParseRefreshToken(input.RefreshToken)
	user, _ := h.services.GetUserData(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, outputForms.ErrorResponse{Error: "Ошибка чтения токена"})
		return
	}
	token, _, err := h.services.Authorization.GenerateTokenByID(userId, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, outputForms.ErrorResponse{Error: "Ошибка регенерации токена"})
		return
	}
	c.JSON(http.StatusOK, outputForms.RefreshTokenResponse{Token: token})
}
