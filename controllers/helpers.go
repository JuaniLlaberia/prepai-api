package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetUserId(context *gin.Context) (bson.ObjectID, error) {
	userIdInterface, exists := context.Get("userId")
	if !exists {
		return bson.NilObjectID, errors.New("user id not found in context")
	}

	userId, ok := userIdInterface.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, errors.New("invalid user id format")
	}

	return userId, nil
}
