// Package handler provides HTTP request handlers for the public and auth routes.
package handler

import (
	"time"

	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/internal/session"
	"my-portfolio/pkg/cryptoutil"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLoginPage renders the admin login form.
// If the user already has a valid admin session in Redis they are redirected
// straight to the dashboard so they never see the login page twice.
func AdminLoginPage(_ *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()

		// Already logged in? Forward to dashboard.
		token := c.Cookies(cfg.Admin.CookieName)
		if token != "" {
			adminID, err := session.Get(rdb, token)
			if err == nil && adminID > 0 {
				return c.Redirect().To("/admin/")
			}
		}

		return c.Render("admin/login", fiber.Map{
			"Title":          "Admin Login",
			"HCaptchaKey":    cfg.HCaptcha.SiteKey,
			"HCaptchaEnable": cfg.HCaptcha.Enabled,
		})
	}
}

// AdminLoginSubmit processes the admin login form.
// hCaptcha verification (when enabled) is handled upstream by the
// gofiber/contrib/v3/hcaptcha middleware registered in router.go.
// On success the session is stored in Redis (not in the SQLite admin row),
// so it survives Go server restarts.
func AdminLoginSubmit(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()

		email := c.FormValue("email")
		password := c.FormValue("password")

		loginErrMsg, _ := appI18n.T.Localize(c, "login_invalid")
		if loginErrMsg == "" {
			loginErrMsg = "Invalid email or password."
		}

		var admin model.Admin
		if err := db.Where("email = ?", email).First(&admin).Error; err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title":          "Admin Login",
				"Error":          loginErrMsg,
				"HCaptchaKey":    cfg.HCaptcha.SiteKey,
				"HCaptchaEnable": cfg.HCaptcha.Enabled,
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
			return c.Render("admin/login", fiber.Map{
				"Title":          "Admin Login",
				"Error":          loginErrMsg,
				"HCaptchaKey":    cfg.HCaptcha.SiteKey,
				"HCaptchaEnable": cfg.HCaptcha.Enabled,
			})
		}

		// Generate session token and persist to Redis with the configured TTL.
		token := cryptoutil.RandomHex(32)
		ttl := time.Duration(cfg.Admin.SessionTTL) * time.Minute
		if err := session.Set(rdb, token, admin.ID, ttl); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("session store error")
		}

		// Only update last_login_at in SQLite (for audit); token is in Redis.
		now := time.Now().UTC()
		db.Model(&admin).Update("last_login_at", now)

		// Expire any old cookie that was previously set with Path: "/admin"
		// (narrower path takes precedence in browsers, so we must clear it first).
		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    "",
			Path:     "/admin",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			SameSite: "Strict",
			MaxAge:   -1,
		})

		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    token,
			Path:     "/",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			SameSite: "Strict",
			MaxAge:   cfg.Admin.SessionTTL * 60,
		})

		return c.Redirect().To("/admin?toast=login_success")
	}
}

// AdminLogout deletes the Redis session and clears all cookie paths.
func AdminLogout(rdb *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		token := c.Cookies(cfg.Admin.CookieName)
		if token != "" {
			_ = session.Delete(rdb, token)
		}
		// Clear both possible cookie paths (Path: "/admin" legacy + Path: "/" current).
		for _, path := range []string{"/", "/admin"} {
			c.Cookie(&fiber.Cookie{
				Name:     cfg.Admin.CookieName,
				Value:    "",
				Path:     path,
				HTTPOnly: true,
				Secure:   cfg.Admin.CookieSecure,
				MaxAge:   -1,
			})
		}
		return c.Redirect().To("/admin/login?toast=logout_success")
	}
}
