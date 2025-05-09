package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
)

type Exam struct {
	Id         bson.ObjectID           `json:"id" bson:"_id,omitempty"`
	Title      string                  `json:"full_name" bson:"full_name,omitempty"`
	Subject    string                  `json:"subject" bson:"subject,omitempty"`
	Difficulty string                  `json:"difficulty" bson:"difficulty,omitempty"`
	ExamType   string                  `json:"exam_type" bson:"exam_type,omitempty"`
	Taken      bool                    `json:"taken" bson:"taken,omitempty"`
	Pinned     bool                    `json:"pinned" bson:"pinned,omitempty"`
	Passed     bool                    `json:"passed" bson:"passed,omitempty"`
	Questions  []internal.ExamQuestion `json:"questions" bson:"questions,omitempty"`
	UserId     bson.ObjectID           `json:"user_id" bson:"user_id"`
	ActividyId bson.ObjectID           `json:"activity_id" bson:"activity_id"`
}
