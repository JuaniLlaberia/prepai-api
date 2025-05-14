package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
	"prepai.app/middlewares"
)

func ResumeRoute(server *gin.Engine) {
	authResume := server.Group("/resumes")
	authResume.Use(middlewares.Authenticate)

	// GET
	authResume.GET("", controllers.GetResumes)
	authResume.GET("/:id", controllers.GetResume)
	// POST
	authResume.POST("", controllers.CreateResume)
	// DELETE
	authResume.DELETE("/:id", controllers.DeleteResume)
}
