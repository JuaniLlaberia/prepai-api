package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
	"prepai.app/models"
)

func GetQuestions(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	questions, err := models.GetAllUserQuestions(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch questions"})
		return
	}

	context.JSON(http.StatusOK, questions)
}

func GetQuestion(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	questionId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid interview ID format",
		})
	}

	question, err := models.GetQuestionById(questionId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch question"})
		return
	}

	if question.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Question does not belong to you",
		})
		return
	}

	context.JSON(http.StatusOK, question)
}

func CreateQuestion(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	var question models.Question
	err = context.ShouldBindJSON(&question)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request body",
		})
		return
	}

	if question.Question == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing question field",
		})
		return
	}

	result, err := internal.GenerateQuestionAnalysis(question.Question)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	question.UserId = userId
	question.Type = result.Type
	question.Difficulty = result.Difficulty
	question.Explanation = result.Explanation
	question.ExpectedLength = result.ExpectedLength
	question.IdealAnswer = result.IdealAnswer

	err = question.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Question created successfully",
		"data":    question,
	})
}

func DeleteQuestion(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	questionId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid interview ID format",
		})
	}

	question, err := models.GetQuestionById(questionId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch question"})
		return
	}

	if question.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Question does not belong to you",
		})
		return
	}

	err = question.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Question deleted successfully",
	})
}
