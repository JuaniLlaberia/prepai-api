package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
	"prepai.app/middlewares"
)

func InterviewRoute(server *gin.Engine) {
	authInterview := server.Group("/interviews")
	authInterview.Use(middlewares.Authenticate)

	// GET
	authInterview.GET("", controllers.GetInterviews)
	authInterview.GET("/:id", controllers.GetInterview)
	authInterview.GET("/:id/attempt", controllers.GetInterviewAttempt)
	// POST
	authInterview.POST("", controllers.CreateInterview)
	authInterview.POST("/:id/attempt", controllers.CreateInterviewAttempt)
	// PATCH
	authInterview.PATCH("/:id", controllers.UpdateInterview)
	authInterview.PATCH("/:id/regenerate", controllers.RegenerateInterview)
	authInterview.PATCH("/:id/attempt/feedback", controllers.CreateInterviewAttemptFeedback)
	// DELETE
	authInterview.DELETE("/:id", controllers.DeleteInterview)
}
