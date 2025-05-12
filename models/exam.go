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

type Exam struct {
	Id         bson.ObjectID           `json:"id" bson:"_id,omitempty"`
	Title      string                  `json:"title" bson:"title,omitempty"`
	Subject    string                  `json:"subject" bson:"subject,omitempty"`
	Difficulty string                  `json:"difficulty" bson:"difficulty,omitempty"`
	Type       string                  `json:"type" bson:"type,omitempty"`
	Taken      bool                    `json:"taken" bson:"taken,omitempty"`
	Pinned     bool                    `json:"pinned" bson:"pinned,omitempty"`
	Passed     bool                    `json:"passed" bson:"passed,omitempty"`
	Questions  []internal.ExamQuestion `json:"questions" bson:"questions,omitempty"`
	UserId     bson.ObjectID           `json:"user_id" bson:"user_id"`
	ActividyId bson.ObjectID           `json:"activity_id" bson:"activity_id"`
}

func GetExams(userId bson.ObjectID) ([]Exam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("exams")
	projection := bson.M{
		"questions": 0,
	}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var results []Exam
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetExamById(examId bson.ObjectID, showAnswers bool) (*Exam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("exams")

	var projection bson.M

	if !showAnswers {
		projection = bson.M{
			"questions": bson.M{
				"correct":     0,
				"explanation": 0,
			},
		}
	}

	opts := options.FindOne().SetProjection(projection)

	var exam Exam
	err := collection.FindOne(ctx, bson.M{"_id": examId}, opts).Decode(&exam)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &exam, nil
}

func (exam *Exam) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("exams")
	result, err := collection.InsertOne(ctx, exam)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	exam.Id = id
	return nil
}

func (exam Exam) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("exams")
	update := bson.M{
		"$set": bson.M{
			"title":     exam.Title,
			"taken":     exam.Taken,
			"passed":    exam.Passed,
			"pinned":    exam.Pinned,
			"questions": exam.Questions,
		},
	}

	_, err := collection.UpdateByID(ctx, exam.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (exam Exam) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("exams")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": exam.Id})
	if err != nil {
		return err
	}

	return nil
}
