package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ITokenService interface {
	CreateToken(email string) (string, error)
	VerifyToken(tokenString string) error
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type TokenService struct {
	secretKey []byte
}

func NewTokenService(secretKey []byte) *TokenService {
	return &TokenService{secretKey}
}

func (tm *TokenService) CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (tm *TokenService) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tm.secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}

func (tm *TokenService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tm.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Erro ao fazer parse do token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token inválido")
	}

	return claims, nil
}
