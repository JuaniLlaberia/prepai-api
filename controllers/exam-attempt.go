package controllers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/models"
)

func GetExamAttempt(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	examAttempt, err := models.GetAttemptByExamId(examId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch exam attempt",
		})
		return
	}

	if examAttempt.UserId != userId {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Exam attempt does not belong to you",
		})
		return
	}

	context.JSON(http.StatusOK, examAttempt)
}

func CreateExamAttempt(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	var examAttempt models.ExamAttempt
	examAttempt.UserId = userId
	examAttempt.ExamId = examId

	err = examAttempt.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Exam attempt created successfully",
	})
}

type ExamSubmission struct {
	Responses []int64
	Time      int64
}

func SubmitExamAttempt(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	var userResponse ExamSubmission
	err = context.ShouldBindJSON(&userResponse)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	if len(userResponse.Responses) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Responses array cannot be empty",
		})
		return
	}

	exam, err := models.GetExamById(examId, true)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch exam",
		})
		return
	}

	if exam.UserId != userId {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Exam attempt does not belong to you",
		})
		return
	}

	answers := make([]models.ExamAnswer, len(userResponse.Responses))
	totalScore := 0.0

	for i, response := range userResponse.Responses {
		if response == exam.Questions[i].Correct {
			totalScore++
		}

		answers[i] = models.ExamAnswer{
			Question:    exam.Questions[i].Question,
			Answer:      response,
			Correct:     exam.Questions[i].Correct,
			Explanation: exam.Questions[i].Explanation,
		}
	}

	examAttempt, err := models.GetAttemptByExamId(examId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch exam attempt",
		})
		return
	}

	totalQuestions := float64(len(answers))

	scoreOutOf10 := (totalScore / totalQuestions) * 10.0
	scoreOutOf10Rounded := math.Round(scoreOutOf10*10) / 10
	percentageScore := (totalScore / totalQuestions) * 100.0
	passed := percentageScore >= 70.0

	examAttempt.Answers = answers
	examAttempt.Time = userResponse.Time
	examAttempt.Score = scoreOutOf10Rounded
	examAttempt.Passed = passed

	err = examAttempt.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update exam attempt: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Exam submitted successfully",
		"data":    examAttempt,
	})
}
