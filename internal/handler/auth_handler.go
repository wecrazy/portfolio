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

		// if secure cookies are required but we're not on HTTPS, warn the user.
		insecureMsg := ""
		if cfg.Admin.CookieSecure && c.Protocol() != "https" {
			insecureMsg, _ = appI18n.T.Localize(c, "login_insecure")
		}
		return c.Render("admin/login", fiber.Map{
			"Title":          "Admin Login",
			"HCaptchaKey":    cfg.HCaptcha.SiteKey,
			"HCaptchaEnable": cfg.HCaptcha.Enabled,
			"InsecureNotice": insecureMsg,
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
			// log.Printf("AdminLoginSubmit: session set error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("session store error")
		}
		// log.Printf("AdminLoginSubmit: token %s stored for adminID=%d ttl=%v", token, admin.ID, ttl)

		// Build a cookie configuration that respects the current request
		// protocol.  In development the config file often sets CookieSecure
		// true, which prevents storage over plain HTTP and results in a
		// perpetual login loop.  Override when the request is not HTTPS.
		secure := cfg.Admin.CookieSecure
		if c.Protocol() != "https" {
			secure = false
		}

		// Helper to apply domain if provided (blank = omit field).
		makeCookie := func(name, value, path string, maxAge int) *fiber.Cookie {
			ck := &fiber.Cookie{
				Name:     name,
				Value:    value,
				Path:     path,
				HTTPOnly: true,
				Secure:   secure,
				SameSite: "Lax", // Lax relaxes strict navigation rules, more forgiving
				MaxAge:   maxAge,
			}
			if cfg.Admin.CookieDomain != "" {
				ck.Domain = cfg.Admin.CookieDomain
			}
			return ck
		}

		// Clear any legacy /admin cookie first (narrow path wins).
		c.Cookie(makeCookie(cfg.Admin.CookieName, "", "/admin", -1))

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
