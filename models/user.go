package models

import (
	"context"
	"errors"
	"go-rest-api/config"
	"go-rest-api/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email", binding:"required`
	Password string             `json:"password", binding:"required`
}

func (u User) Save() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(UsersCollection)

	var existingUser User
	if err := collection.FindOne(ctx, bson.M{"email": u.Email}).Decode(&existingUser); err == nil {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return err
	}

	u.Password = string(hashedPassword)

	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (u *User) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(UsersCollection)
	var existingUser User

	err := collection.FindOne(ctx, bson.M{"email": u.Email}).Decode(&existingUser)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(u.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return utils.GenerateToken(existingUser.Email, existingUser.ID.Hex())
}
