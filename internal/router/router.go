// Package router registers all application routes and global middleware.
package router

import (
	"errors"
	"strings"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/handler"
	"my-portfolio/internal/handler/admin"
	"my-portfolio/internal/hub"
	"my-portfolio/internal/middleware"

	contribcb "github.com/gofiber/contrib/v3/circuitbreaker"
	contribfgprof "github.com/gofiber/contrib/v3/fgprof"
	contribhcaptcha "github.com/gofiber/contrib/v3/hcaptcha"
	contribloadshed "github.com/gofiber/contrib/v3/loadshed"
	contribmonitor "github.com/gofiber/contrib/v3/monitor"
	contribswaggo "github.com/gofiber/contrib/v3/swaggo"
	contribws "github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"gorm.io/gorm"
)

// RegisterRoutes wires up every route and middleware in the application.
func RegisterRoutes(app *fiber.App, db *gorm.DB, h *hub.Hub) {
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
			return p == healthcheck.LivenessEndpoint ||
				p == healthcheck.ReadinessEndpoint ||
				strings.HasPrefix(p, "/static") ||
				strings.HasPrefix(p, "/uploads") ||
				strings.HasPrefix(p, "/ws")
		},
		OnShed: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":       "Server is temporarily overloaded",
				"retry_after": 5,
			})
		},
	}))

	// ── 2. Health checks ───────────────────────────────────────────
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(_ fiber.Ctx) bool { return db != nil },
	}))

	// ── 3. Swagger UI (/swagger/*) ─────────────────────────────────
	app.Get("/swagger/*", contribswaggo.HandlerDefault)

	// ── 4. Circuit breaker for public API routes ───────────────────
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
	// Stop the circuit breaker when the server shuts down.
	app.Hooks().OnPreShutdown(func() error {
		cb.Stop()
		return nil
	})

	// ── 5. WebSocket upgrade + real-time comments ──────────────────
	app.Use("/ws", func(c fiber.Ctx) error {
		if contribws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/comments", handler.WSComments(h))

	// ── 6. Public routes ───────────────────────────────────────────
	// Locale JSON — single source of truth from web/locales/*.yaml
	app.Get("/lang/:code", handler.LangJSON("web/locales"))

	app.Get("/", handler.PortfolioPage(db))
	app.Get("/resume", handler.ServeResumePDF(db))
	app.Get("/resume/download", handler.DownloadResumePDF(db))

	// OAuth
	app.Get("/auth/google", handler.GoogleLogin())
	app.Get("/auth/google/callback", handler.GoogleCallback(db))
	app.Get("/auth/github", handler.GitHubLogin())
	app.Get("/auth/github/callback", handler.GitHubCallback(db))
	app.Get("/auth/logout", handler.OAuthLogout())

	// Comments (circuit-breaker protected)
	app.Get("/comments", contribcb.Middleware(cb), handler.GetComments(db))
	app.Post("/comments", contribcb.Middleware(cb), middleware.OAuthAuth(), handler.PostComment(db, h))

	// Projects (HTMX lazy load)
	app.Get("/projects", handler.ProjectsPage(db))

	// Blog (circuit-breaker protected)
	app.Get("/blog/more", contribcb.Middleware(cb), handler.BlogPostsPartial(db))
	app.Get("/blog", handler.BlogListPage(db))
	app.Get("/blog/:slug", handler.BlogPostPage(db))

	// Contact (circuit-breaker protected)
	app.Post("/contact", contribcb.Middleware(cb), handler.SubmitContact(db))

	// ── 7. Admin auth (unauthenticated) ───────────────────────────
	app.Get("/admin/login", handler.AdminLoginPage())

	// hCaptcha middleware on the login POST — always registered.
	// The wrapper checks live config at request-time so toggling enabled/disabled
	// via hot-reload takes effect without a server restart.
	hcaptchaHandler := contribhcaptcha.New(contribhcaptcha.Config{
		SecretKey: cfg.HCaptcha.Secret,
		// Read the token from the HTML form field (not default JSON body).
		ResponseKeyFunc: func(c fiber.Ctx) (string, error) {
			return c.FormValue("h-captcha-response"), nil
		},
		// On failure: re-render the login page with an error.
		// Must return a non-nil error so the middleware does NOT call c.Next().
		ValidateFunc: func(success bool, c fiber.Ctx) error {
			if !success {
				liveCfg := config.MyPortfolio.Get()
				// Write the response body first, then return a sentinel error.
				// The middleware sees the body is non-empty and returns nil,
				// skipping c.Next() so AdminLoginSubmit never runs.
				_ = c.Render("admin/login", fiber.Map{
					"Title":          "Admin Login",
					"Error":          "Captcha verification failed. Please try again.",
					"HCaptchaKey":    liveCfg.HCaptcha.SiteKey,
					"HCaptchaEnable": true,
				})
				return errors.New("captcha failed")
			}
			// Return nil → middleware calls c.Next() → AdminLoginSubmit runs.
			return nil
		},
	})
	app.Post("/admin/login",
		func(c fiber.Ctx) error {
			if !config.MyPortfolio.Get().HCaptcha.Enabled {
				return c.Next()
			}
			return hcaptchaHandler(c)
		},
		handler.AdminLoginSubmit(db),
	)
	app.Post("/admin/logout", handler.AdminLogout())

	// ── 8. Admin routes (session-protected) ───────────────────────
	adm := app.Group("/admin", middleware.AdminAuth(db))

	adm.Get("/", admin.Dashboard(db))

	// Server monitor (live metrics dashboard)
	adm.Get("/monitor", contribmonitor.New(contribmonitor.Config{
		Title: "Portfolio · Server Monitor",
	}))

	// Full goroutine-aware profiler — access /admin/debug/fgprof?seconds=10
	adm.Use(contribfgprof.New(contribfgprof.Config{Prefix: "/admin"}))

	// Owner / About
	adm.Get("/owner", admin.OwnerEditPage(db))
	adm.Put("/owner", admin.OwnerUpdate(db))
	adm.Post("/owner/upload-image", admin.OwnerUploadImage(db))
	adm.Post("/owner/upload-resume", admin.OwnerUploadResume(db))

	// Projects CRUD
	adm.Get("/projects", admin.ProjectListPage())
	adm.Get("/projects/list", admin.ProjectListPartial(db))
	adm.Get("/projects/new", admin.ProjectNewForm())
	adm.Get("/projects/:id/edit", admin.ProjectEditForm(db))
	adm.Post("/projects", admin.ProjectCreate(db))
	adm.Put("/projects/:id", admin.ProjectUpdate(db))
	adm.Delete("/projects/:id", admin.ProjectDelete(db))

	// Experience CRUD
	adm.Get("/experience", admin.ExperienceListPage())
	adm.Get("/experience/list", admin.ExperienceListPartial(db))
	adm.Get("/experience/new", admin.ExperienceNewForm())
	adm.Get("/experience/:id/edit", admin.ExperienceEditForm(db))
	adm.Post("/experience", admin.ExperienceCreate(db))
	adm.Put("/experience/:id", admin.ExperienceUpdate(db))
	adm.Delete("/experience/:id", admin.ExperienceDelete(db))
	adm.Post("/experience/upload-image", admin.ExperienceUploadImage(db))

	// Skills CRUD
	adm.Get("/skills", admin.SkillListPage())
	adm.Get("/skills/list", admin.SkillListPartial(db))
	adm.Get("/skills/new", admin.SkillNewForm())
	adm.Get("/skills/:id/edit", admin.SkillEditForm(db))
	adm.Post("/skills", admin.SkillCreate(db))
	adm.Put("/skills/:id", admin.SkillUpdate(db))
	adm.Delete("/skills/:id", admin.SkillDelete(db))

	// Social Links CRUD
	adm.Get("/social-links", admin.SocialListPage())
	adm.Get("/social-links/list", admin.SocialListPartial(db))
	adm.Get("/social-links/new", admin.SocialNewForm())
	adm.Get("/social-links/:id/edit", admin.SocialEditForm(db))
	adm.Post("/social-links", admin.SocialCreate(db))
	adm.Put("/social-links/:id", admin.SocialUpdate(db))
	adm.Delete("/social-links/:id", admin.SocialDelete(db))

	// Uploads
	adm.Get("/uploads", admin.UploadListPage())
	adm.Get("/uploads/list", admin.UploadListPartial(db))
	adm.Post("/uploads", admin.UploadCreate(db))
	adm.Delete("/uploads/:id", admin.UploadDelete(db))

	// Comment moderation
	adm.Get("/comments", admin.CommentListPage())
	adm.Get("/comments/list", admin.CommentListPartial(db))
	adm.Put("/comments/:id/approve", admin.CommentApprove(db))
	adm.Put("/comments/:id/reject", admin.CommentReject(db))
	adm.Delete("/comments/:id", admin.CommentDelete(db))
	adm.Post("/comments/:id/reply", admin.CommentReply(db))

	// Contact messages
	adm.Get("/contacts", admin.ContactListPage())
	adm.Get("/contacts/list", admin.ContactListPartial(db))
	adm.Put("/contacts/:id/read", admin.ContactMarkRead(db))
	adm.Delete("/contacts/:id", admin.ContactDelete(db))

	// Tech Stack CRUD
	adm.Get("/tech-stacks", admin.TechStackListPage())
	adm.Get("/tech-stacks/list", admin.TechStackListPartial(db))
	adm.Get("/tech-stacks/new", admin.TechStackNewForm())
	adm.Get("/tech-stacks/:id/edit", admin.TechStackEditForm(db))
	adm.Post("/tech-stacks", admin.TechStackCreate(db))
	adm.Put("/tech-stacks/:id", admin.TechStackUpdate(db))
	adm.Delete("/tech-stacks/:id", admin.TechStackDelete(db))

	// Blog Posts CRUD
	adm.Get("/posts", admin.PostListPage())
	adm.Get("/posts/list", admin.PostListPartial(db))
	adm.Get("/posts/new", admin.PostNewForm())
	adm.Get("/posts/:id/edit", admin.PostEditForm(db))
	adm.Post("/posts", admin.PostCreate(db))
	adm.Put("/posts/:id", admin.PostUpdate(db))
	adm.Delete("/posts/:id", admin.PostDelete(db))
	adm.Post("/posts/upload-thumbnail", admin.PostUploadThumbnail(db))

	// Upcoming Items CRUD
	adm.Get("/upcoming", admin.UpcomingListPage())
	adm.Get("/upcoming/list", admin.UpcomingListPartial(db))
	adm.Get("/upcoming/new", admin.UpcomingNewForm())
	adm.Get("/upcoming/:id/edit", admin.UpcomingEditForm(db))
	adm.Post("/upcoming", admin.UpcomingCreate(db))
	adm.Put("/upcoming/:id", admin.UpcomingUpdate(db))
	adm.Delete("/upcoming/:id", admin.UpcomingDelete(db))
}
