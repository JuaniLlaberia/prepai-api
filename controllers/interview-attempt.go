package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
	"prepai.app/models"
)

func GetInterviewAttempt(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	interviewId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	interviewAttempt, err := models.GetAttemptByInterviewId(interviewId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch interview attempt",
		})
		return
	}

	if interviewAttempt.UserId != userId {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "This interview attempt does not belong to you",
		})
		return

	}

	context.JSON(http.StatusOK, interviewAttempt)

}

func CreateInterviewAttempt(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	interviewId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	var interviewAttempt models.InterviewAttempt
	interviewAttempt.UserId = userId
	interviewAttempt.InterviewId = interviewId

	err = interviewAttempt.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Interview attempt created successfully",
	})
}

func CreateInterviewAttemptFeedback(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	interviewId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	var userResponses []internal.UserInterviewResponse
	err = context.ShouldBindJSON(&userResponses)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	interviewAttempt, err := models.GetAttemptByInterviewId(interviewId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch interview attempt",
		})
		return
	}

	if interviewAttempt.UserId != userId {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "This interview attempt does not belong to you",
		})
		return

	}

	results, err := internal.GenerateInterviewFeedback(userResponses)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	answers := make([]models.InterviewAnswer, len(userResponses))
	totalScore := 0.0

	for i, userResponse := range userResponses {
		feedback := results.Feedbacks[i]
		answers[i] = models.InterviewAnswer{
			Question:     userResponse.Question,
			UserResponse: userResponse.Answer,
			Feedback:     feedback.Feedback,
			Score:        feedback.Score,
			Suggestion:   feedback.Suggestion,
		}
		totalScore += feedback.Score
	}

	averageScore := totalScore / float64(len(answers))
	passed := averageScore >= 70.0

	interviewAttempt.Answers = answers
	interviewAttempt.Passed = passed
	interviewAttempt.Score = averageScore
	interviewAttempt.Analysis = results.Analysis
	interviewAttempt.AreasToImprove = results.AreasToImprove
	interviewAttempt.Strengths = results.Strengths

	err = interviewAttempt.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update interview attempt: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Interview feedback generated successfully",
		"data":    interviewAttempt,
	})
}
