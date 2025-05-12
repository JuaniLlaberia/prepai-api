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

type Question struct {
	Id             bson.ObjectID           `json:"id" bson:"_id,omitempty"`
	Question       string                  `json:"question" bson:"question,omitempty"`
	Type           string                  `json:"type" bson:"type,omitempty"`
	Difficulty     string                  `json:"difficulty" bson:"difficulty,omitempty"`
	Explanation    string                  `json:"explanation" bson:"explanation,omitempty"`
	ExpectedLength string                  `json:"expected_length" bson:"expected_length,omitempty"`
	IdealAnswer    internal.QuestionAnswer `json:"ideal_answer" bson:"ideal_answer,omitempty"`
	UserId         bson.ObjectID           `json:"user_id" bson:"user_id"`
	Pinned         bool                    `json:"pinned" bson:"pinned,omitempty"`
}

func GetAllUserQuestions(userId bson.ObjectID) ([]Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("questions")
	projection := bson.M{
		"question":   1,
		"type":       1,
		"difficulty": 1,
		"pinned":     1,
	}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var results []Question
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}

func GetQuestionById(questionId bson.ObjectID) (*Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("questions")

	var question Question
	err := collection.FindOne(ctx, bson.M{"_id": questionId}).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &question, nil
}

func (question *Question) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("questions")
	result, err := collection.InsertOne(ctx, question)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	question.Id = id
	return nil
}

func (question Question) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("questions")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": question.Id})
	if err != nil {
		return err
	}

	return nil
}
