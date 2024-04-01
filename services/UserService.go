package services

import (
	"AutenticaoUsuario/models"
	"AutenticaoUsuario/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
)

type IUserService interface {
	ValidateAndCreateUser(userRegistration models.UserRegistration) (models.UserRegistration, error)
	GenerateTokenRecoverPassword(email string) (string, error)
	RecoverPassword(token string, novaSenha string) error
	UpdateUser(user models.UserRegistration) error
	GetUserByCredentials(email string, senha string) (models.UserRegistration, error)
}

type UserService struct {
	UserRepository  repository.IUserRepositoryInterface
	TokenService    ITokenService
	MailService     IMailService
	PasswordService IPasswordService
}

func NewUserService(userRepo repository.IUserRepositoryInterface,
	tokenService ITokenService,
	mailService IMailService, passwordService IPasswordService) *UserService {
	return &UserService{
		UserRepository:  userRepo,
		TokenService:    tokenService,
		MailService:     mailService,
		PasswordService: passwordService,
	}
}

func validatePassword(password string) bool {
	// Regex para verificar se a senha contém pelo menos uma letra e um número
	passwordRegex := regexp.MustCompile(`[A-Za-z].*[0-9]|[0-9].*[A-Za-z]`)
	if !passwordRegex.MatchString(password) {
		return false
	}
	return true
}

func validateEmail(email string) bool {
	// Regex para verificar se o email é válido
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (us UserService) ValidateAndCreateUser(userRegistration models.UserRegistration) (models.UserRegistration, error) {
	var err error
	if userRegistration.Name == "" {
		return models.UserRegistration{}, errors.New("username cannot be empty")
	}

	mailValid := validateEmail(userRegistration.Email)
	if !mailValid {
		return models.UserRegistration{}, errors.New("email inválido")
	}

	if userRegistration.Password == "" {
		return models.UserRegistration{}, errors.New("password cannot be empty")
	}

	if len(userRegistration.Password) < 3 {
		return models.UserRegistration{}, errors.New("password must be at least 3 characters long")
	}

	passwordValid := validatePassword(userRegistration.Password)
	if !passwordValid {
		return models.UserRegistration{}, errors.New("password must contain at least one letter and one number")
	}

	userRegistration.Password, err = us.PasswordService.HashPassword(userRegistration.Password)
	if err != nil {
		return models.UserRegistration{}, fmt.Errorf("Falhar ao criptografar senha.")
	}

	existingUser, err := us.UserRepository.GetUserByEmail(userRegistration.Email)
	if err != nil {
		return models.UserRegistration{}, err
	}
	if existingUser.Email != "" {
		return models.UserRegistration{}, errors.New("email already exists")
	}

	userRegistration.ID = uuid.New()

	err = us.UserRepository.CreateUser(userRegistration)
	if err != nil {
		return models.UserRegistration{}, err
	}

	return userRegistration, nil
}

func (us UserService) RecoverPassword(token string, novaSenha string) error {
	err := us.TokenService.VerifyToken(token)
	if err != nil {
		return fmt.Errorf("Token inválido: %v", err)
	}

	// Obter o email associado ao token (você pode ter que modificar isso dependendo de como seu token é estruturado)
	claims, err := us.TokenService.ParseToken(token)
	if err != nil {
		return fmt.Errorf("Erro ao obter claims do token: %v", err)
	}
	emailClaim, ok := claims["email"].(string)
	if !ok {
		return errors.New("Não foi possível obter o email do token")
	}

	user, err := us.UserRepository.GetUserByEmail(emailClaim)
	if err != nil {
		return fmt.Errorf("Erro ao buscar usuário pelo email: %v", err)
	}
	if user.Email == "" {
		return errors.New("Usuário não encontrado")
	}

	passwordValid := validatePassword(user.Password)
	if !passwordValid {
		return fmt.Errorf("password must contain at least one letter and one number")
	}
	user.Password, err = us.PasswordService.HashPassword(novaSenha)
	if err != nil {
		return fmt.Errorf("Falhar ao criptografar senha.")
	}
	err = us.UserRepository.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("Erro ao atualizar a senha do usuário: %v", err)
	}

	return nil
}

func (us UserService) GenerateTokenRecoverPassword(email string) (string, error) {
	// Verificar se o email existe
	user, err := us.UserRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user.Email == "" {
		return "", errors.New("email not found")
	}

	//Cria token com validade de 24 horas
	tokenString, err := us.TokenService.CreateToken(email)
	if err != nil {
		return "", err
	}

	err = us.MailService.SendEmail(user.Email, "Recuperação de Senha", "Seu token de recuperação de senha com validade de 24 horas: "+tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us UserService) UpdateUser(user models.UserRegistration) error {
	passwordValid := validatePassword(user.Password)
	if !passwordValid {
		return fmt.Errorf("password must contain at least one letter and one number")
	}
	var err error
	user.Password, err = us.PasswordService.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("Falhar ao criptografar senha.")
	}
	mailValid := validateEmail(user.Email)
	if !mailValid {
		return errors.New("email inválido")
	}
	err = us.UserRepository.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("Error updating user: %v", err)
	}
	return nil
}

func (us UserService) GetUserByCredentials(email string, senha string) (models.UserRegistration, error) {
	var err error

	user, err := us.UserRepository.GetUserByEmail(email)
	if !user.Active {
		return user, errors.New("user disabled")
	}
	senhaValida := us.PasswordService.CheckPasswordHash(senha, user.Password)
	if !senhaValida {
		return models.UserRegistration{}, fmt.Errorf("Senha INválida.")
	}
	return user, err
}

func (us UserService) GetUserByMail(email string) (models.UserRegistration, error) {
	user, err := us.UserRepository.GetUserByEmail(email)
	return user, err
}
