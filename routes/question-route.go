package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
	"prepai.app/middlewares"
)

func QuestionRoute(server *gin.Engine) {
	authQuestion := server.Group("/questions")
	authQuestion.Use(middlewares.Authenticate)

	// GET
	authQuestion.GET("", controllers.GetQuestions)
	authQuestion.GET("/:id", controllers.GetQuestion)
	// POST
	authQuestion.POST("", controllers.CreateQuestion)
	// DELETE
	authQuestion.DELETE("/:id", controllers.DeleteQuestion)
}
