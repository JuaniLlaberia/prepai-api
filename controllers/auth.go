package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"prepai.app/configs"
	"prepai.app/models"
	"prepai.app/utils"
)

func Signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

func Login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type OAuthState struct {
	State            string
	FrontendRedirect string
	ExpiresAt        time.Time
}

func OAuthLogin(oauthConfig *oauth2.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		frontendRedirect := context.Query("redirect_url")
		if frontendRedirect == "" {
			frontendRedirect = "/"
		}

		stateId := uuid.New().String()
		authURL := oauthConfig.AuthCodeURL(stateId)

		context.JSON(http.StatusOK, gin.H{
			"auth_url": authURL,
		})
	}
}

func GithubCallback(ctx *gin.Context) {
	code := ctx.Query("code")

	token, err := configs.GetGithubOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to exchange token",
		})
		return
	}

	email, err := fetchGithubEmail(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.GetOrCreateUser(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	jwtToken, err := utils.GenerateToken(email, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}

func GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")

	token, err := configs.GetGoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to exchange token",
		})
		return
	}

	email, err := fetchGoogleEmail(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.GetOrCreateUser(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	jwtToken, err := utils.GenerateToken(email, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}

func fetchGithubEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %v", token))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get github")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type GithubEmail struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	var emails []GithubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", errors.New("no verified email found")
}

func fetchGoogleEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get github")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type GoogleUserInfo struct {
		Email string `json:"email"`
	}

	var email GoogleUserInfo
	if err := json.Unmarshal(body, &email); err != nil {
		return "", err
	}

	fmt.Println("GOOGLE EMAIL", email.Email)
	return email.Email, nil
}
