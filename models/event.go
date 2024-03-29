package models

import (
	"context"
	"log"
	"time"

	"go-rest-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name", binding:"required`
	Description string             `json:"description", binding:"required`
	Location    string             `json:"location", binding:"required`
	DateTime    time.Time          `json:"dateTime", binding:"required`
	UserID      primitive.ObjectID
}

func (e *Event) Save() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(EventsCollection)
	_, err := collection.InsertOne(ctx, e)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetAllEvents() ([]Event, error) {

	var events []Event

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(EventsCollection)

	options := options.Find()

	cursor, err := collection.Find(ctx, bson.D{{}}, options)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event Event
		err := cursor.Decode(&event)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		events = append(events, event)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return events, nil
}

func GetEvent(id string) (Event, error) {

	var event Event

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(EventsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return event, err
	}

	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func (e *Event) Update(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(EventsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: objectID}}, bson.D{{Key: "$set", Value: e}})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func DeleteEvent(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database(config.DatabaseName).Collection(EventsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectID}})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
