package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/models"
)

func GetUser(context *gin.Context) {
	userIdInterface, exists := context.Get("userId")
	if !exists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}

	userId, ok := userIdInterface.(bson.ObjectID)
	if !ok {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}

	user, err := models.GetUser(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get user.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully.",
		"user":    user,
	})
}

func UpdateUser(context *gin.Context) {
	// Logic to get and validate bsonObjectID (needs extra step)
	userIdInterface, exists := context.Get("userId")
	if !exists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}

	userId, ok := userIdInterface.(bson.ObjectID)
	if !ok {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}

	// Get and validate new user data
	var updatedUser models.User
	err := context.ShouldBindJSON(&updatedUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updatedUser.Id = userId

	err = updatedUser.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully.", "user": updatedUser})
}

func DeleteUser(context *gin.Context) {
	userIdInterface, exists := context.Get("userId")
	if !exists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}

	userId, ok := userIdInterface.(bson.ObjectID)
	if !ok {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}

	var user models.User
	user.Id = userId
	user.Delete()

	context.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully.",
	})
}
