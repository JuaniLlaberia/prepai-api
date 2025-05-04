package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"prepai.app/configs"
	"prepai.app/utils"
)

type User struct {
	Id       bson.ObjectID `bson:"_id,omitempty"`
	FullName string        `bson:"full_name,omitempty"`
	Email    string        `bson:"email" validate:"required"`
	Password string        `bson:"password" validate:"required"`
	ImageUrl string        `bson:"image_url,omitempty"`
}

func (user User) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("users")

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key") {
			return errors.New("email address is already taken")
		}

		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	user.Id = id
	return nil
}

func (user *User) ValidateCredentials() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("users")

	projection := bson.M{
		"_id":      1,
		"password": 1,
	}
	opts := options.FindOne().SetProjection(projection)

	var result User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}, opts)
	if err != nil {
		if err.Err() == mongo.ErrNoDocuments {
			return errors.New("creadentials invalid")
		}
		return err.Err()
	}

	isPasswordValid := utils.CheckPasswordHash(user.Password, result.Password)

	if !isPasswordValid {
		return errors.New("creadentials invalid")
	}

	return nil
}
