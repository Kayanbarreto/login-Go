package controllers

import (
	"AutenticaoUsuario/models"
	"AutenticaoUsuario/repository"
	"AutenticaoUsuario/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	UserRepo    repository.IUserRepositoryInterface
	UserService services.IUserService
}

func NewUserController(userRepo repository.IUserRepositoryInterface, userService services.IUserService) *UserController {
	return &UserController{
		UserRepo:    userRepo,
		UserService: userService,
	}
}

func (uc *UserController) CreateUserHandler(c *gin.Context) {
	var user models.UserRegistration
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
		return
	}

	userValidate, err := uc.UserService.ValidateAndCreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.UserRepo.CreateUser(userValidate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userValidate)
}

func (uc *UserController) GetUserHandler(c *gin.Context) {
	var credentials models.UserLogin
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.UserRepo.GetUserByCredentials(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) RecoverPasswordHandler(c *gin.Context) {
	var requestData struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
		return
	}

	if requestData.Token == "" || requestData.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token or new password cannot be empty"})
		return
	}

	err := uc.UserService.RecoverPassword(requestData.Token, requestData.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha recuperada com sucesso"})
}

func (uc *UserController) UpdateUserHandler(c *gin.Context) {
	var updateUser models.UserRegistration
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar JSON"})
		return
	}

	existingUser, err := uc.UserRepo.GetUserByEmail(updateUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar usuário"})
		return
	}
	if existingUser.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não existente com o e-mail informado."})
		return
	}
	updateUser.ID = existingUser.ID

	err = uc.UserService.UpdateUser(updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}
