package controllers

import (
	"AutenticaoUsuario/models"
	"AutenticaoUsuario/repository"
	"AutenticaoUsuario/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type LoginController struct {
	TokenService    services.ITokenService
	UserRepo        repository.IUserRepositoryInterface
	UserService     services.IUserService
	PasswordService services.IPasswordService
}

func NewLoginController(tokenService services.ITokenService,
	userRepo repository.IUserRepositoryInterface,
	userService services.IUserService,
	passwordService services.IPasswordService) *LoginController {
	return &LoginController{
		TokenService:    tokenService,
		UserRepo:        userRepo,
		UserService:     userService,
		PasswordService: passwordService,
	}
}

func (lc *LoginController) LoginHandler(c *gin.Context) {
	var user models.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar o JSON"})
		return
	}

	userBanco, err := lc.UserService.GetUserByCredentials(user.Email, user.Password)
	if err != nil {
		if !userBanco.Active && userBanco.Email != "" {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Usuário desativado."})
			return
		}
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario nao encontrado"})
		return
	}

	if !userBanco.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "Usuário desativado"})
		return
	}

	if lc.PasswordService.CheckPasswordHash(user.Password, userBanco.Password) {
		tokenString, err := lc.TokenService.CreateToken(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
	}
}

func (lc *LoginController) ProtectedHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta o cabeçalho de autorização"})
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := lc.TokenService.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bem-vindo à área protegida"})
}
