package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/internal"
	"prepai.app/models"
)

func GetResumes(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	resumes, err := models.GetAllUserResumes(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch resumes. Try again later."})
		return
	}

	context.JSON(http.StatusOK, resumes)
}

func GetResume(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	resumeId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid interview ID format",
		})
		return
	}

	resume, err := models.GetResumeById(resumeId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch resume. Try again later."})
		return
	}

	if resume.UserId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Resume does not belong to you",
		})
		return
	}

	context.JSON(http.StatusOK, resume)
}

func CreateResume(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	header, err := context.FormFile("resume")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing resume file or error uploading",
		})
		return
	}

	// Functionality to upload it somewhere (NOT IMPLEMENTED YET)
	//

	const maxSize = 5 << 20 // 5 MB
	if header.Size > maxSize {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "File size exceeds the limit (5MB)",
		})
		return
	}

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Only PDF files are allowed",
		})
		return
	}

	file, err := header.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error opening file",
		})
		return
	}
	defer file.Close()

	jobDescription := context.PostForm("job_description")
	if jobDescription == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Job description is required",
		})
		return
	}

	if jobDescription == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Job description is required",
		})
		return
	}

	result, err := internal.ResumeAnalyzer(file, header, jobDescription)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error analyzing resume: " + err.Error(),
		})
		return
	}

	var resume models.Resume
	resume.UserId = userId
	resume.Title = result.Title
	resume.OverallScore = result.OverallScore
	resume.AnalysisSummary = result.AnalysisSummary
	resume.ImprovementSuggestions = result.ImprovementSuggestions
	resume.Metrics = result.Metrics

	err = resume.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Resume analyzed successfully",
		"data":    resume,
	})
}

func DeleteResume(context *gin.Context) {
	userId, err := GetUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	resumeId, err := bson.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid resume ID format",
		})
		return
	}

	resume, err := models.GetResumeById(resumeId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch resume document",
		})
		return
	}

	if resume.UserId != userId {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Resume does not belong to you",
		})
		return
	}

	err = resume.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Resume deleted successfully",
	})
}
