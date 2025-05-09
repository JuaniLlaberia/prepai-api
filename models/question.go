package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
)

type Question struct {
	Id             bson.ObjectID           `json:"id" bson:"_id,omitempty"`
	Question       string                  `json:"question" bson:"question,omitempty"`
	Type           string                  `json:"type" bson:"type,omitempty"`
	Difficulty     string                  `json:"difficulty" bson:"difficulty,omitempty"`
	Exmplanation   string                  `json:"explanation" bson:"explanation,omitempty"`
	ExpectedLength string                  `json:"expected_length" bson:"expected_length,omitempty"`
	IdealAnswer    internal.QuestionAnswer `json:"ideal_answer" bson:"ideal_answer,omitempty"`
	UserId         bson.ObjectID           `json:"user_id" bson:"user_id"`
	Pinned         bool                    `json:"pinned" bson:"pinned,omitempty"`
}
