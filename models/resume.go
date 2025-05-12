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

type Resume struct {
	Id                     bson.ObjectID    `json:"id" bson:"_id,omitempty"`
	FileUrl                string           `json:"file_url" bson:"file_url,empty"`
	Title                  string           `json:"title" bson:"title,omitempty"`
	OverallScore           int64            `json:"overall_score" bson:"overall_score,omitempty"`
	AnalysisSummary        string           `json:"analysis_summary" bson:"analysis_summary,omitempty"`
	ImprovementSuggestions string           `json:"improvement_suggestions" bson:"improvement_suggestions,omitempty"`
	Metrics                internal.Metrics `json:"metrics" bson:"metrics,omitempty"`
	UserId                 bson.ObjectID    `json:"user_id" bson:"user_id"`
}

func GetAllUserResumes(userId bson.ObjectID) ([]Resume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("resumes")
	projection := bson.M{
		"file_url":      1,
		"title":         1,
		"overall_score": 1,
	}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var results []Resume
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetResumeById(resumeId bson.ObjectID) (*Resume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("resumes")

	var resume Resume
	err := collection.FindOne(ctx, bson.M{"_id": resumeId}).Decode(&resume)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &resume, nil
}

func (resume *Resume) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("resumes")
	result, err := collection.InsertOne(ctx, resume)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.New("failed to get document id")
	}

	resume.Id = id
	return nil
}

func (resume Resume) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.GetCollection("resumes")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": resume.Id})
	if err != nil {
		return err
	}

	return nil
}
