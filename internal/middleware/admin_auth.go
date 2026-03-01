// Package middleware provides Fiber middleware for authentication and security.
package middleware

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/internal/session"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AdminAuth checks for a valid admin session in Redis and loads the admin into
// c.Locals("admin"). Redirects to /admin/login when not authenticated.
// Sessions survive Go server restarts because they live in Redis.
func AdminAuth(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		token := c.Cookies(cfg.Admin.CookieName)
		if token == "" {
			return c.Redirect().To("/admin/login")
		}

		// Look up session in Redis — returns 0 when token is missing / expired.
		adminID, err := session.Get(rdb, token)
		if err != nil || adminID == 0 {
			return c.Redirect().To("/admin/login")
		}

		var admin model.Admin
		if err := db.First(&admin, adminID).Error; err != nil {
			return c.Redirect().To("/admin/login")
		}

		c.Locals("admin", admin)
		return c.Next()
	}
}
