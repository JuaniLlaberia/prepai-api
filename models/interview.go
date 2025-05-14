package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"prepai.app/configs"
	"prepai.app/internal"
)

type Interview struct {
	Id         bson.ObjectID                `json:"id" bson:"_id,omitempty"`
	Title      string                       `json:"title" bson:"title,omitempty"`
	JobRole    string                       `json:"job_role" bson:"job_role,omitempty" validate:"required"`
	JobLevel   string                       `json:"job_level" bson:"job_level,omitempty" validate:"required"`
	Topics     []string                     `json:"topics" bson:"topics,omitempty" validate:"required"`
	Taken      bool                         `json:"taken" bson:"taken,omitempty"`
	Pinned     bool                         `json:"pinned" bson:"pinned,omitempty"`
	Passed     bool                         `json:"passed" bson:"passed,omitempty"`
	Questions  []internal.InterviewQuestion `json:"questions" bson:"questions,omitempty"`
	UserId     bson.ObjectID                `json:"user_id" bson:"user_id"`
	ActividyId bson.ObjectID                `json:"activity_id" bson:"activity_id"`
}

func GetAllUserInterviews(userId bson.ObjectID) ([]Interview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviews")
	projection := bson.M{
		"questions": 0,
	}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var results []Interview
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetInterviewById(interviewId bson.ObjectID) (*Interview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviews")

	var interview Interview
	err := collection.FindOne(ctx, bson.M{"_id": interviewId}).Decode(&interview)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &interview, nil
}

func (interview *Interview) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviews")
	result, err := collection.InsertOne(ctx, interview)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	interview.Id = id
	return nil
}

func (interview Interview) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviews")
	update := bson.M{
		"$set": bson.M{
			"title":     interview.Title,
			"taken":     interview.Taken,
			"passed":    interview.Passed,
			"pinned":    interview.Pinned,
			"questions": interview.Questions,
		},
	}

	_, err := collection.UpdateByID(ctx, interview.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (interview Interview) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviews")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": interview.Id})
	if err != nil {
		return err
	}

	return nil
}
