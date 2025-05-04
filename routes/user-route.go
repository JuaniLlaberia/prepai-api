package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/controllers"
)

func UserRoute(server *gin.Engine) {
	server.POST("/signup", controllers.Signup)
	server.POST("/login", controllers.Login)
}
