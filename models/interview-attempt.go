package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
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
	UserId         bson.ObjectID     `json:"user_id" bson:"user_id"`
	InterviewId    bson.ObjectID     `json:"interview_id" bson:"interview_id"`
}
