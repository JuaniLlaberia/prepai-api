package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
	"prepai.app/middlewares"
)

func UserRoute(server *gin.Engine) {
	authenticatedRoute := server.Group("/user")
	authenticatedRoute.Use(middlewares.Authenticate)
	authenticatedRoute.GET("", controllers.GetUser)
	authenticatedRoute.PATCH("/update", controllers.UpdateUser)
	authenticatedRoute.DELETE("/delete", controllers.DeleteUser)
}
