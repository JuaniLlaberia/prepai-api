package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func GetGithubOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ProcessEnv("GITHUB_CLIENT_ID"),
		ClientSecret: ProcessEnv("GITHUB_SECRET_ID"),
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user"},
	}

}

func GetGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ProcessEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: ProcessEnv("GOOGLE_SECRET_ID"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}
