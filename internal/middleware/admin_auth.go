// Package middleware provides Fiber middleware for authentication and security.
package middleware

import (
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// AdminAuth checks for a valid admin session cookie and loads the admin into
// c.Locals("admin"). Redirects to /admin/login when not authenticated.
func AdminAuth(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		token := c.Cookies(cfg.Admin.CookieName)
		if token == "" {
			return c.Redirect().To("/admin/login")
		}

		// Use Find+Limit instead of First to avoid GORM logging ErrRecordNotFound
		// for every request with a stale or missing session cookie.
		var admins []model.Admin
		db.Where("session_token = ?", token).Limit(1).Find(&admins)
		if len(admins) == 0 {
			return c.Redirect().To("/admin/login")
		}
		admin := admins[0]

		// Check session TTL.
		ttl := time.Duration(cfg.Admin.SessionTTL) * time.Minute
		if admin.LastLoginAt != nil && time.Since(*admin.LastLoginAt) > ttl {
			// Session expired – clear it.
			db.Model(&admin).Updates(map[string]interface{}{"session_token": ""})
			return c.Redirect().To("/admin/login")
		}

		c.Locals("admin", admin)
		return c.Next()
	}
}
