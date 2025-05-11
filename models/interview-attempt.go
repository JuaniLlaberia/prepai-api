package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"prepai.app/configs"
)

type InterviewAnswer struct {
	Question     string  `json:"question" bson:"question,omitempty"`
	UserResponse string  `json:"user_response" bson:"user_response,omitempty"`
	Feedback     string  `json:"feedback" bson:"feedback,omitempty"`
	Score        float64 `json:"score" bson:"score,omitempty"`
	Suggestion   string  `json:"suggestion" bson:"suggestion,omitempty"`
}

type InterviewAttempt struct {
	Id             bson.ObjectID     `json:"id" bson:"_id,omitempty"`
	Answers        []InterviewAnswer `json:"answers" bson:"answers,omitempty"`
	Analysis       string            `json:"analysis" bson:"analysis,omitempty"`
	Strengths      []string          `json:"strengths" bson:"strengths,omitempty"`
	AreasToImprove []string          `json:"areas_to_improve" bson:"areas_to_improve,omitempty"`
	Passed         bool              `json:"passed" bson:"passed,omitempty"`
	Score          float64           `json:"score" bson:"score,omitempty"`
	UserId         bson.ObjectID     `json:"user_id" bson:"user_id"`
	InterviewId    bson.ObjectID     `json:"interview_id" bson:"interview_id"`
}

func GetAttemptByInterviewId(interviewId bson.ObjectID) (*InterviewAttempt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviewAttempts")

	var attempt InterviewAttempt
	err := collection.FindOne(ctx, bson.M{"interview_id": interviewId}).Decode(&attempt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &attempt, nil
}

func (attempt *InterviewAttempt) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviewAttempts")
	result, err := collection.InsertOne(ctx, attempt)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	attempt.Id = id
	return nil
}

func (attempt InterviewAttempt) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("interviewAttempts")
	update := bson.M{
		"$set": bson.M{
			"answers":          attempt.Answers,
			"passed":           attempt.Passed,
			"analysis":         attempt.Analysis,
			"areas_to_improve": attempt.AreasToImprove,
			"strengths":        attempt.Strengths,
		},
	}

	_, err := collection.UpdateByID(ctx, attempt.Id, update)
	if err != nil {
		return err
	}

	return nil
}
