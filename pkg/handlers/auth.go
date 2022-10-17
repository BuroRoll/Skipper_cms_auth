package handlers

import (
	"Skipper_cms_auth/pkg/models/forms"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /

// @Schemes
// @Description Регистрация нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param   phone      		query     string     true  "Телефон для авторизации"
// @Param   first_name      query     string     true  "Имя пользователя"
// @Param   second_name     query     string     true  "Фамилия пользователя"
// @Param   password      	query     string     true  "Пароль"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input forms.SignUpUserForm
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная форма регистрации"})
		return
	}
	_, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания профиля"})
		return
	}
	token, refreshToken, err := h.services.Authorization.GenerateToken(input.Phone, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токенов"})
		return
	}
	response := map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
	}
	c.JSON(http.StatusOK, response)
}

// @Schemes
// @Description Вход
// @Tags Auth
// @Accept json
// @Produce json
// @Param   login      		query     string     true  "Логин для авторизации"
// @Param   password      	query     string     true  "Пароль"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input forms.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная форма авторизации"})
		return
	}
	token, refreshToken, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
	})
}

// @Schemes
// @Description Обновление токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param   refresh_token      		query     string     true  "Refresh Token"
// @Router /auth/refresh-token [post]
func (h *Handler) refreshToken(c *gin.Context) {
	var input forms.TokenReqBody
	err := c.Bind(&input)
	userId, err := h.services.ParseRefreshToken(input.RefreshToken)
	user, _ := h.services.GetUserData(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка чтения токена"})
		return
	}
	token, _, err := h.services.Authorization.GenerateTokenByID(userId, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регенерации токена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
