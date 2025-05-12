package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"prepai.app/configs"
)

type ExamAnswer struct {
	Question    string `json:"question"`
	Answer      int64  `json:"answer"`
	Correct     int64  `json:"correct"`
	Explanation string `json:"explanation"`
}

type ExamAttempt struct {
	Id      bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Time    int64         `json:"time" bson:"time,omitempty"`
	Score   float64       `json:"score" bson:"score,omitempty"`
	Answers []ExamAnswer  `json:"answers" bson:"answers,omitempty"`
	Passed  bool          `json:"passed" bson:"passed,omitempty"`
	UserId  bson.ObjectID `json:"user_id" bson:"user_id"`
	ExamId  bson.ObjectID `json:"exam_id" bson:"exam_id"`
}

func GetAttemptByExamId(examId bson.ObjectID) (*ExamAttempt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("examAttempts")

	var attempt ExamAttempt
	err := collection.FindOne(ctx, bson.M{"exam_id": examId}).Decode(&attempt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &attempt, nil
}

func (attempt *ExamAttempt) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("examAttempts")
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

func (attempt ExamAttempt) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("examAttempts")
	update := bson.M{
		"$set": bson.M{
			"answers": attempt.Answers,
			"passed":  attempt.Passed,
			"score":   attempt.Score,
			"time":    attempt.Time,
		},
	}

	_, err := collection.UpdateByID(ctx, attempt.Id, update)
	if err != nil {
		return err
	}

	return nil
}
