package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
)

type Interview struct {
	Id         bson.ObjectID                `json:"id" bson:"_id,omitempty"`
	Title      string                       `json:"full_name" bson:"full_name,omitempty"`
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
