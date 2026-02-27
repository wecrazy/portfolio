package router

import (
	"errors"

	"my-portfolio/internal/config"
	"my-portfolio/internal/handler"

	contribhcaptcha "github.com/gofiber/contrib/v3/hcaptcha"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// registerAuthRoutes wires the admin login/logout endpoints.
//
// hCaptcha is always registered but the inner check reads live config at
// request-time, so toggling enabled/disabled via hot-reload takes effect
// without a server restart.
//
//   - loginLimiter — strict IP-based brute-force guard (10 req / 15 min)
func registerAuthRoutes(app *fiber.App, db *gorm.DB, loginLimiter fiber.Handler) {
	cfg := config.MyPortfolio.Get()

	app.Get("/admin/login", handler.AdminLoginPage())

	// hCaptcha middleware: on failure it re-renders the login page with an
	// error and returns a non-nil error so AdminLoginSubmit is never called.
	hcaptchaMiddleware := contribhcaptcha.New(contribhcaptcha.Config{
		SecretKey: cfg.HCaptcha.Secret,
		ResponseKeyFunc: func(c fiber.Ctx) (string, error) {
			return c.FormValue("h-captcha-response"), nil
		},
		ValidateFunc: func(success bool, c fiber.Ctx) error {
			if !success {
				liveCfg := config.MyPortfolio.Get()
				_ = c.Render("admin/login", fiber.Map{
					"Title":          "Admin Login",
					"Error":          "Captcha verification failed. Please try again.",
					"HCaptchaKey":    liveCfg.HCaptcha.SiteKey,
					"HCaptchaEnable": true,
				})
				return errors.New("captcha failed")
			}
			return nil
		},
	})

	app.Post("/admin/login",
		// Brute-force guard: 10 attempts per 15 minutes per IP.
		loginLimiter,
		// hCaptcha gate: skipped when disabled in config (checked at runtime).
		func(c fiber.Ctx) error {
			if !config.MyPortfolio.Get().HCaptcha.Enabled {
				return c.Next()
			}
			return hcaptchaMiddleware(c)
		},
		handler.AdminLoginSubmit(db),
	)

	app.Post("/admin/logout", handler.AdminLogout())
}
