package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
	"prepai.app/middlewares"
)

func ExamRoute(server *gin.Engine) {
	authExam := server.Group("/exams")
	authExam.Use(middlewares.Authenticate)

	// GET
	authExam.GET("", controllers.GetExams)
	authExam.GET("/:id", controllers.GetExam)
	authExam.GET("/:id/attempt", controllers.GetExamAttempt)
	// POST
	authExam.POST("", controllers.CreateExam)
	authExam.POST("/:id/attempt", controllers.CreateExamAttempt)

	// PATCH
	authExam.PATCH("/:id", controllers.UpdateExam)
	authExam.PATCH("/:id/regenerate", controllers.RegenerateExam)
	authExam.PATCH("/:id/attempt/submit", controllers.SubmitExamAttempt)
	// DELETE
	authExam.DELETE("/:id", controllers.DeleteExam)
}
