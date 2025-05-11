package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
	"prepai.app/models"
)

func GetInterviews(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	interviews, err := models.GetAllUserInterviews(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch interviews. Try again later."})
		return
	}

	context.JSON(http.StatusOK, interviews)
}

func GetInterview(context *gin.Context) {
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

	interview, err := models.GetInterviewById(interviewId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch interview. Try again later."})
		return
	}

	if interview.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Interview does not belong to you",
		})
		return
	}

	context.JSON(http.StatusOK, interview)
}

func CreateInterview(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	var interview models.Interview
	interview.UserId = userId

	err = context.ShouldBindJSON(&interview)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
	}

	results, err := internal.GenerateInterview(interview.JobRole, interview.JobLevel, interview.Topics)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	interview.Title = results.Title
	interview.Questions = results.Questions

	err = interview.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Interview created successfully",
		"data":    interview,
	})
}

func RegenerateInterview(context *gin.Context) {
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

	interview, err := models.GetInterviewById(interviewId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch interview. Try again later."})
		return
	}

	if interview.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Interview does not belong to you",
		})
		return
	}

	results, err := internal.GenerateInterview(interview.JobRole, interview.JobLevel, interview.Topics)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	interview.Title = results.Title
	interview.Questions = results.Questions

	err = interview.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Interview questions regenerated successfully",
		"data":    interview,
	})
}

func UpdateInterview(context *gin.Context) {
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

	interview, err := models.GetInterviewById(interviewId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch interview. Try again later."})
		return
	}

	if interview.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Interview does not belong to you",
		})
		return
	}

	err = context.ShouldBindJSON(&interview)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
	}

	err = interview.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Interview updated successfully",
		"data":    interview,
	})
}

func DeleteInterview(context *gin.Context) {
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

	interview, err := models.GetInterviewById(interviewId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch interview. Try again later."})
		return
	}

	if interview.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Interview does not belong to you",
		})
		return
	}

	err = interview.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Interview deleted successfully",
	})
}
