package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
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
