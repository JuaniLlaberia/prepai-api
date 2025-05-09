package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
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
