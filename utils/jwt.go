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
		"userId": userId.Hex(),
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (bson.ObjectID, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("not authorized")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return bson.NilObjectID, errors.New("not authorized")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userClaim, exists := claims["userId"]
		if !exists {
			return bson.NilObjectID, errors.New("user claim not found in token")
		}

		switch userID := userClaim.(type) {
		case string:
			objID, err := bson.ObjectIDFromHex(userID)
			if err != nil {
				return bson.NilObjectID, errors.New("invalid user ID format")
			}

			return objID, nil

		default:
			return bson.NilObjectID, errors.New("user claim has unexpected format")
		}
	}

	return bson.NilObjectID, errors.New("invalid token")
}
