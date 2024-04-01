package main

import (
	"AutenticaoUsuario/controllers"
	"AutenticaoUsuario/repository"
	"AutenticaoUsuario/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	router := gin.Default()

	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}
	fmt.Println("Conectado ao MongoDB!")

	database := client.Database("seu_banco_de_dados")
	secretKey := []byte("your_secret_key")

	userRepo := repository.NewUserRepository(database)
	tokenService := services.NewTokenService(secretKey)
	mailService := services.NewMailService()
	passwordService := services.NewPasswordService()
	userService := services.NewUserService(userRepo, tokenService, mailService, passwordService)
	loginController := controllers.NewLoginController(tokenService, userRepo, userService, passwordService)
	userController := controllers.NewUserController(userRepo, userService)
	tokenController := controllers.NewTokenController(userService)

	router.POST("/v1/login", loginController.LoginHandler)
	router.GET("/v1/protected", loginController.ProtectedHandler)
	router.POST("/v1/user/register", userController.CreateUserHandler)
	router.POST("/v1/user", userController.GetUserHandler)
	router.POST("/v1/user/generate-token-recover-password", tokenController.GenerateTokenRecoverPasswordHandler)
	router.POST("/v1/user/recover-password", userController.RecoverPasswordHandler)
	router.PUT("/v1/user/updateUser", userController.UpdateUserHandler)

	fmt.Println("Starting the server")
	err = router.Run("localhost:4000")
	if err != nil {
		fmt.Println("Could not start the server", err)
	}
	fmt.Println("Server started. Listening on port 4000")

}
