// Package handler provides HTTP request handlers for the public and auth routes.
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLoginPage renders the admin login form.
func AdminLoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/login", fiber.Map{
			"Title": "Admin Login",
		})
	}
}

// AdminLoginSubmit processes the admin login form.
func AdminLoginSubmit(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		var admin model.Admin
		if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title": "Admin Login",
				"Error": "Invalid username or password",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title": "Admin Login",
				"Error": "Invalid username or password",
			})
		}

		// Generate session token.
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			log.Printf("Failed to generate session token: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
		}
		token := hex.EncodeToString(tokenBytes)
		now := time.Now().UTC()

		db.Model(&admin).Updates(map[string]interface{}{
			"session_token": token,
			"last_login_at": now,
		})

		cfg := config.MyPortfolio.Get()
		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    token,
			Path:     "/admin",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			SameSite: "Strict",
			MaxAge:   cfg.Admin.SessionTTL * 60,
		})

		return c.Redirect("/admin")
	}
}

// AdminLogout clears the admin session.
func AdminLogout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    "",
			Path:     "/admin",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			MaxAge:   -1,
		})
		return c.Redirect("/admin/login")
	}
}

func generateSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
