package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"prepai.app/configs"
)

var secretKey = configs.ProcessEnv("JWT_SECRET")

func GenerateToken(email string, userId bson.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (bson.ObjectID, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Not authorized")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return bson.NilObjectID, errors.New("Not authorized")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return bson.NilObjectID, errors.New("Not authorized")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return bson.NilObjectID, errors.New("Not authorized")
	}

	userId := claims["user"].(bson.ObjectID)
	return userId, nil
}
