// Package handler provides HTTP request handlers for the public and auth routes.
package handler

import (
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/cryptoutil"

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
		email := c.FormValue("email")
		password := c.FormValue("password")

		var admin model.Admin
		if err := db.Where("email = ?", email).First(&admin).Error; err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title": "Admin Login",
				"Error": "Invalid email or password",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title": "Admin Login",
				"Error": "Invalid email or password",
			})
		}

		// Generate session token.
		token := cryptoutil.RandomHex(32)
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

		return c.Redirect("/admin?toast=login_success")
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
		return c.Redirect("/admin/login?toast=logout_success")
	}
}
