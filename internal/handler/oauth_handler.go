package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"my-portfolio/internal/middleware"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GoogleLogin redirects the user to Google's consent screen.
func GoogleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthCfg := service.GoogleOAuthConfig()
		state := generateState()
		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Value:    state,
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge:   300,
		})
		url := oauthCfg.AuthCodeURL(state)
		return c.Redirect(url)
	}
}

// GoogleCallback handles the OAuth2 callback from Google.
func GoogleCallback(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		savedState := c.Cookies("oauth_state")
		if c.Query("state") != savedState || savedState == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid OAuth state")
		}

		oauthCfg := service.GoogleOAuthConfig()
		token, err := oauthCfg.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to exchange token")
		}

		client := oauthCfg.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var userInfo struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		json.Unmarshal(body, &userInfo)

		user := upsertOAuthUser(db, "google", userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
		setVisitorSessionCookie(c, user)
		return c.Redirect("/#comments")
	}
}

// GitHubLogin redirects the user to GitHub's consent screen.
func GitHubLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthCfg := service.GitHubOAuthConfig()
		state := generateState()
		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Value:    state,
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge:   300,
		})
		url := oauthCfg.AuthCodeURL(state)
		return c.Redirect(url)
	}
}

// GitHubCallback handles the OAuth2 callback from GitHub.
func GitHubCallback(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		savedState := c.Cookies("oauth_state")
		if c.Query("state") != savedState || savedState == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid OAuth state")
		}

		oauthCfg := service.GitHubOAuthConfig()
		token, err := oauthCfg.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to exchange token")
		}

		client := oauthCfg.Client(context.Background(), token)

		// Fetch user profile.
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var ghUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		}
		json.Unmarshal(body, &ghUser)

		displayName := ghUser.Name
		if displayName == "" {
			displayName = ghUser.Login
		}

		// If email is empty, fetch from emails API.
		email := ghUser.Email
		if email == "" {
			email = fetchGitHubPrimaryEmail(client)
		}

		providerID := fmt.Sprintf("%d", ghUser.ID)
		user := upsertOAuthUser(db, "github", providerID, email, displayName, ghUser.AvatarURL)
		setVisitorSessionCookie(c, user)
		return c.Redirect("/#comments")
	}
}

// OAuthLogout clears the visitor session.
func OAuthLogout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("visitor_session")
		if token != "" {
			middleware.DeleteVisitorSession(token)
		}
		c.Cookie(&fiber.Cookie{
			Name:     "visitor_session",
			Value:    "",
			HTTPOnly: true,
			MaxAge:   -1,
		})
		return c.Redirect("/")
	}
}

// upsertOAuthUser finds or creates an OAuthUser record.
func upsertOAuthUser(db *gorm.DB, provider, providerID, email, name, avatar string) model.OAuthUser {
	var user model.OAuthUser
	db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user)

	if user.ID == 0 {
		user = model.OAuthUser{
			Provider:    provider,
			ProviderID:  providerID,
			Email:       email,
			DisplayName: name,
			AvatarURL:   avatar,
		}
		db.Create(&user)
	} else {
		db.Model(&user).Updates(map[string]interface{}{
			"email":        email,
			"display_name": name,
			"avatar_url":   avatar,
		})
	}
	return user
}

func setVisitorSessionCookie(c *fiber.Ctx, user model.OAuthUser) {
	token := generateState()
	middleware.SetVisitorSession(token, user)
	c.Cookie(&fiber.Cookie{
		Name:     "visitor_session",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   86400 * 7, // 7 days
	})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func fetchGitHubPrimaryEmail(client *http.Client) string {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	json.Unmarshal(body, &emails)

	for _, e := range emails {
		if e.Primary {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}
