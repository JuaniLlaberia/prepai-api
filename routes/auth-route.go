package routes

import (
	"github.com/gin-gonic/gin"
	"prepai.app/configs"
	"prepai.app/controllers"
)

func AuthRoute(server *gin.Engine) {
	auth := server.Group("/auth")
	// Email - Password
	auth.POST("/signup", controllers.Signup)
	auth.POST("/login", controllers.Login)
	// Github OAuth
	auth.GET("/github", controllers.OAuthLogin(configs.GetGithubOauthConfig()))
	auth.GET("/github/callback", controllers.GithubCallback)
	// Google OAuth
	auth.GET("/google", controllers.OAuthLogin(configs.GetGoogleOauthConfig()))
	auth.GET("/google/callback", controllers.GoogleCallback)
}
