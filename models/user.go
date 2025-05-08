package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"prepai.app/configs"
	"prepai.app/utils"
)

type User struct {
	Id       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName string        `json:"full_name" bson:"full_name,omitempty"`
	Email    string        `json:"email" bson:"email" validate:"required,email"`
	Password string        `json:"password" bson:"password" validate:"required"`
	ImageUrl string        `json:"image_url" bson:"image_url,omitempty"`
}

func GetUser(userId bson.ObjectID) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("users")

	projection := bson.M{
		"password": 0,
	}
	opts := options.FindOne().SetProjection(projection)

	var user User
	err := collection.FindOne(ctx, bson.M{"_id": userId}, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func GetOrCreateUser(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("users")

	projection := bson.M{
		"password": 0,
	}
	opts := options.FindOne().SetProjection(projection)

	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.Email = email
			user.Save()

			return &user, nil
		}
		return nil, err
	}

	return &user, nil
}

func (user *User) Save() error {
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
	err := collection.FindOne(ctx, bson.M{"email": user.Email}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("credentials invalid")
		}
		return err
	}

	isPasswordValid := utils.CheckPasswordHash(user.Password, result.Password)
	if !isPasswordValid {
		return errors.New("credentials invalid")
	}

	user.Id = result.Id
	return nil
}

func (user User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("users")
	update := bson.M{
		"$set": bson.M{
			"full_name": user.FullName,
			"image_url": user.ImageUrl,
		},
	}

	_, err := collection.UpdateByID(ctx, user.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (user User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println(user)

	collection := configs.GetCollection("users")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": user.Id})
	if err != nil {
		return err
	}

	return nil
}
