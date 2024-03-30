package models

import (
	"context"
	"go-rest-api/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Registration struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	EventId primitive.ObjectID `json:"eventId", binding:"required`
	UserId  primitive.ObjectID `json:"userId", binding:"required`
}

func (e Event) Register(userId string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(RegistrationsCollection)

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	registration := Registration{
		EventId: e.ID,
		UserId:  userID,
	}

	_, err = collection.InsertOne(ctx, registration)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (e Event) Unregister(userId string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(RegistrationsCollection)

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var existingRegistration Registration
	err = collection.FindOne(ctx, bson.M{"eventid": e.ID, "userid": userID}).Decode(&existingRegistration)
	if err != nil {
		return err
	}

	// err = collection.FindOne(ctx, bson.M{"userid": userID,}).Decode(&existingRegistration)
	// if err != nil {
	// 	return err
	// }

	_, err = collection.DeleteOne(ctx, existingRegistration)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
