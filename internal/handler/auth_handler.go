// Package handler provides HTTP request handlers for the public and auth routes.
package handler

import (
	"time"

	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/cryptoutil"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLoginPage renders the admin login form.
func AdminLoginPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
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
func AdminLoginSubmit(db *gorm.DB) fiber.Handler {
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

		// Generate session token.
		token := cryptoutil.RandomHex(32)
		now := time.Now().UTC()

		db.Model(&admin).Updates(map[string]interface{}{
			"session_token": token,
			"last_login_at": now,
		})

		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    token,
			Path:     "/admin",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			SameSite: "Strict",
			MaxAge:   cfg.Admin.SessionTTL * 60,
		})

		return c.Redirect().To("/admin?toast=login_success")
	}
}

// AdminLogout clears the admin session.
func AdminLogout() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		c.Cookie(&fiber.Cookie{
			Name:     cfg.Admin.CookieName,
			Value:    "",
			Path:     "/admin",
			HTTPOnly: true,
			Secure:   cfg.Admin.CookieSecure,
			MaxAge:   -1,
		})
		return c.Redirect().To("/admin/login?toast=logout_success")
	}
}
