package controllers

import (
	"AutenticaoUsuario/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenController struct {
	UserService services.IUserService
}

func NewTokenController(userService services.IUserService) *TokenController {
	return &TokenController{
		UserService: userService,
	}
}

func (tm *TokenController) GenerateTokenRecoverPasswordHandler(c *gin.Context) {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
		return
	}

	if requestData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email cannot be empty"})
		return
	}

	token, err := tm.UserService.GenerateTokenRecoverPassword(requestData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token de recuperação enviado para o email", "token": token})
}
