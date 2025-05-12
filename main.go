package main

import (
	"github.com/gin-gonic/gin"
	"prepai.app/configs"
	"prepai.app/routes"
)

func main() {
	// Setting up database
	configs.ConnectDB()
	configs.InitDatabase()

	server := gin.Default()

	// Routes
	routes.UserRoute(server)
	routes.AuthRoute(server)
	routes.InterviewRoute(server)
	routes.ExamRoute(server)
	routes.QuestionRoute(server)

	server.Run(":8080")
}
