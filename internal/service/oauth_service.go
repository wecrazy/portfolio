package service

import (
	"my-portfolio/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthConfig returns the OAuth2 config for Google.
func GoogleOAuthConfig() *oauth2.Config {
	cfg := config.MyPortfolio.Get()
	return &oauth2.Config{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

// GitHubOAuthConfig returns the OAuth2 config for GitHub.
func GitHubOAuthConfig() *oauth2.Config {
	cfg := config.MyPortfolio.Get()
	return &oauth2.Config{
		ClientID:     cfg.OAuth.GitHub.ClientID,
		ClientSecret: cfg.OAuth.GitHub.ClientSecret,
		RedirectURL:  cfg.OAuth.GitHub.RedirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
}
