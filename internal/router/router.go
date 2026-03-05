// Package router registers all application routes and global middleware.
package router

import (
	"encoding/json"
	"strings"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/hub"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/pkg/translate"

	contribcb "github.com/gofiber/contrib/v3/circuitbreaker"
	contribloadshed "github.com/gofiber/contrib/v3/loadshed"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterRoutes wires up every route and middleware in the application.
func RegisterRoutes(app *fiber.App, db *gorm.DB, rdb *redis.Client, h *hub.Hub) {
	cfg := config.MyPortfolio.Get()

	// ── 1. Load shedding ───────────────────────────────────────────
	// Shed requests when CPU > 90%; skip health, static and WebSocket paths.
	app.Use(contribloadshed.New(contribloadshed.Config{
		Criteria: &contribloadshed.CPULoadCriteria{
			LowerThreshold: 0.75,
			UpperThreshold: 0.90,
			Interval:       500 * time.Millisecond,
			Getter:         &contribloadshed.DefaultCPUPercentGetter{},
		},
		Next: func(c fiber.Ctx) bool {
			p := c.Path()
			// these routes are lightweight and should always be served even
			// when we are shedding load; clients rely on them for i18n and
			// comment display. skipping them avoids 503 responses that broke
			// the frontend.
			return p == healthcheck.LivenessEndpoint ||
				p == healthcheck.ReadinessEndpoint ||
				strings.HasPrefix(p, "/static") ||
				strings.HasPrefix(p, "/uploads") ||
				strings.HasPrefix(p, "/ws") ||
				strings.HasPrefix(p, "/lang") ||
				strings.HasPrefix(p, "/comments") ||
				strings.HasPrefix(p, "/api/translate") // TODO: exclude the /api/translate soon if it will be overloaded
		},
		OnShed: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":       "Server is temporarily overloaded",
				"retry_after": 5,
			})
		},
	}))

	// ── 2. Rate limiting ───────────────────────────────────────────
	// makeLimiter builds an IP-keyed rate limiter with a standard 429 response.
	makeLimiter := func(maxReqs int, exp time.Duration) fiber.Handler {
		return limiter.New(limiter.Config{
			Max:          maxReqs,
			Expiration:   exp,
			KeyGenerator: func(c fiber.Ctx) string { return c.IP() },
			LimitReached: func(c fiber.Ctx) error {
				c.Set("Retry-After", exp.String())
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":       "Too many requests — please slow down",
					"retry_after": int(exp.Seconds()),
				})
			},
		})
	}

	// Global: 200 req/min per IP. Static assets / health endpoints are exempt.
	app.Use(limiter.New(limiter.Config{
		Max:          200,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string { return c.IP() },
		Next: func(c fiber.Ctx) bool {
			p := c.Path()
			return strings.HasPrefix(p, "/static") ||
				strings.HasPrefix(p, "/uploads") ||
				p == healthcheck.LivenessEndpoint ||
				p == healthcheck.ReadinessEndpoint
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many requests — please slow down",
				"retry_after": 60,
			})
		},
	}))

	contactMax := cfg.RateLimit.ContactForm
	if contactMax <= 0 {
		contactMax = 5
	}
	commentMax := cfg.RateLimit.Comments
	if commentMax <= 0 {
		commentMax = 20
	}

	// Contact: stricter window (10 min) + HTMX-aware too-many-requests toast.
	contactLimiter := limiter.New(limiter.Config{
		Max:          contactMax,
		Expiration:   10 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string { return c.IP() },
		LimitReached: func(c fiber.Ctx) error {
			c.Set("Retry-After", "600")
			msg, _ := appI18n.T.Localize(c, "contact_rate_limit")
			if msg == "" {
				msg = "Too many messages sent. Please wait 10 minutes before trying again."
			}
			if c.Get("HX-Request") == "true" {
				type payload struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				}
				type trigger struct {
					ContactToast payload `json:"contactToast"`
				}
				b, _ := json.Marshal(trigger{ContactToast: payload{Msg: msg, Type: "warning"}})
				c.Set("HX-Trigger", string(b))
				return c.Status(fiber.StatusTooManyRequests).SendString("")
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       msg,
				"retry_after": 600,
			})
		},
	})
	commentLimiter := makeLimiter(commentMax, 1*time.Minute)
	// Admin login: strict brute-force guard — 10 attempts per 15 minutes per IP.
	loginLimiter := makeLimiter(10, 15*time.Minute)

	// ── 3. Circuit breaker ─────────────────────────────────────────
	cb := contribcb.New(contribcb.Config{
		FailureThreshold: 5,
		Timeout:          10 * time.Second,
		SuccessThreshold: 2,
		OnOpen: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":       "Service temporarily unavailable",
				"retry_after": 10,
			})
		},
	})
	app.Hooks().OnPreShutdown(func() error {
		cb.Stop()
		return nil
	})
	cbMiddleware := contribcb.Middleware(cb)

	// ── 4. Route groups ────────────────────────────────────────────
	translator := translate.New()
	registerAPIRoutes(app, db, cbMiddleware, translator)
	registerPublicRoutes(app, db, rdb, h, cbMiddleware, contactLimiter, commentLimiter)
	registerAuthRoutes(app, db, rdb, loginLimiter)
	registerAdminRoutes(app, db, rdb)
}
