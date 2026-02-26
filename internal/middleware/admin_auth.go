// Package middleware provides Fiber middleware for authentication and security.
package middleware

import (
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AdminAuth checks for a valid admin session cookie and loads the admin into
// c.Locals("admin"). Redirects to /admin/login when not authenticated.
func AdminAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		token := c.Cookies(cfg.Admin.CookieName)
		if token == "" {
			return c.Redirect("/admin/login")
		}

		var admin model.Admin
		if err := db.Where("session_token = ?", token).First(&admin).Error; err != nil {
			return c.Redirect("/admin/login")
		}

		// Check session TTL.
		ttl := time.Duration(cfg.Admin.SessionTTL) * time.Minute
		if admin.LastLoginAt != nil && time.Since(*admin.LastLoginAt) > ttl {
			// Session expired – clear it.
			db.Model(&admin).Updates(map[string]interface{}{"session_token": ""})
			return c.Redirect("/admin/login")
		}

		c.Locals("admin", admin)
		return c.Next()
	}
}
