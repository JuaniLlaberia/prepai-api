package configs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var DB *mongo.Client
var Ctx context.Context
var Cancel context.CancelFunc

func ConnectDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(ProcessEnv("MONGOURI")).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	DB, err = mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := DB.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}

func DisconnectDB() {
	if DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := DB.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from database: %v", err)
		} else {
			fmt.Println("Database disconnected successfully")
		}
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Database("PrepAi").Collection(collectionName)
}

func InitDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := SetupUserCollection(ctx)
	if err != nil {
		log.Fatalf("Failed to set up user collection: %v", err)
	}

}

func SetupUserCollection(ctx context.Context) error {
	collection := GetCollection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create email index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"email", "password"},
		"properties": bson.M{
			"full_name": bson.M{
				"bsonType":    "string",
				"description": "Full name of the user",
			},
			"email": bson.M{
				"bsonType":    "string",
				"pattern":     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				"description": "Email must be a valid email address",
			},
			"password": bson.M{
				"bsonType":    "string",
				"minLength":   8,
				"description": "Password must be at least 8 characters",
			},
			"image_url": bson.M{
				"bsonType":    "string",
				"description": "URL to user's profile image",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "users"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "users", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create users collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}
