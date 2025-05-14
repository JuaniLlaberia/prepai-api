package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
	"prepai.app/models"
)

func GetExams(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	exams, err := models.GetExams(userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch exams."})
		return
	}

	context.JSON(http.StatusOK, exams)
}

func GetExam(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid exam ID format",
		})
		return
	}

	exam, err := models.GetExamById(examId, false)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch exam."})
		return
	}

	if exam.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Exam does not belong to you",
		})
		return
	}

	context.JSON(http.StatusOK, exam)
}

func CreateExam(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	var exam models.Exam
	err = context.ShouldBindJSON(&exam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
		})
		return
	}

	if exam.Subject == "" || exam.Difficulty == "" || exam.Type == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing either subject, difficulty or exam type",
		})
		return
	}

	result, err := internal.GenerateExam(exam.Subject, exam.Difficulty, exam.Type)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	exam.Title = result.Title
	exam.Questions = result.Questions
	exam.UserId = userId

	err = exam.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Exam created successfully",
		"data":    exam,
	})
}

func UpdateExam(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid exam ID format",
		})
		return
	}

	exam, err := models.GetExamById(examId, false)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch exan",
		})
		return
	}

	if exam.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Exam does not belong to you",
		})
		return
	}

	err = context.ShouldBindJSON(&exam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
		})
		return
	}

	err = exam.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Exam updated successfully",
		"data":    exam,
	})
}

func DeleteExam(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid exam ID format",
		})
		return
	}

	exam, err := models.GetExamById(examId, false)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not fetch exan",
		})
		return
	}

	if exam.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Exam does not belong to you",
		})
		return
	}

	err = exam.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Exam deleted successfully",
	})
}

func RegenerateExam(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	examId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid exam ID format",
		})
		return
	}

	exam, err := models.GetExamById(examId, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch exam"})
		return
	}

	if exam.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Exam does not belong to you",
		})
		return
	}

	results, err := internal.GenerateExam(exam.Subject, exam.Difficulty, exam.Type)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	exam.Title = results.Title
	exam.Questions = results.Questions

	err = exam.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Exam questions regenerated successfully",
		"data":    exam,
	})
}
