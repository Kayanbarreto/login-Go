package repository

import (
	userModels "AutenticaoUsuario/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type IUserRepositoryInterface interface {
	CreateUser(userRegistration userModels.UserRegistration) error
	GetUserByCredentials(email string, senha string) (userModels.UserRegistration, error)
	GetUserByEmail(email string) (userModels.UserRegistration, error)
	UpdateUser(user userModels.UserRegistration) error
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(userRegistration userModels.UserRegistration) error {
	_, err := r.Collection.InsertOne(context.TODO(), userRegistration)
	if err != nil {
		log.Println("Error when entering user in MongoDB:", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByCredentials(email string, senha string) (userModels.UserRegistration, error) {
	var user userModels.UserRegistration
	err := r.Collection.FindOne(context.TODO(), bson.M{"email": email, "password": senha}).Decode(&user)
	if err != nil {
		log.Println("UserLogin not found in MongoDB:", err)
		return userModels.UserRegistration{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (userModels.UserRegistration, error) {
	var user userModels.UserRegistration
	err := r.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		// Se não encontrarmos o usuário, não é um erro, então retornamos um usuário vazio e nil
		if err == mongo.ErrNoDocuments {
			return userModels.UserRegistration{}, nil
		}
		log.Println("Error getting user by email from MongoDB:", err)
		return userModels.UserRegistration{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user userModels.UserRegistration) error {
	filter := bson.M{"id": user.ID}
	update := bson.M{"$set": bson.M{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"password":  user.Password,
		"birthDate": user.BirthDate,
		"active":    user.Active,
	}}

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error updating user in MongoDB:", err)
		return err
	}
	return nil
}
