package router

import (
	"my-portfolio/internal/handler"
	"my-portfolio/internal/hub"
	"my-portfolio/internal/middleware"

	contribswaggo "github.com/gofiber/contrib/v3/swaggo"
	contribws "github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"gorm.io/gorm"
)

// registerPublicRoutes wires health checks, Swagger, WebSocket, and all
// unauthenticated public-facing routes.
//
//   - cb             — circuit-breaker middleware (pre-built from router.go)
//   - contactLimiter — per-IP rate limiter for the contact POST
//   - commentLimiter — per-IP rate limiter for comment POSTs
func registerPublicRoutes(
	app *fiber.App,
	db *gorm.DB,
	h *hub.Hub,
	cb fiber.Handler,
	contactLimiter, commentLimiter fiber.Handler,
) {
	// ── Health checks ──────────────────────────────────────────────
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(_ fiber.Ctx) bool { return db != nil },
	}))

	// ── Swagger UI ─────────────────────────────────────────────────
	// Protected behind admin session auth — API docs should not be public.
	app.Get("/swagger/*", middleware.AdminAuth(db), contribswaggo.HandlerDefault)

	// ── WebSocket upgrade ──────────────────────────────────────────
	app.Use("/ws", func(c fiber.Ctx) error {
		if contribws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/comments", handler.WSComments(h))

	// ── Locale JSON ────────────────────────────────────────────────
	// Single source of truth served from web/locales/*.yaml.
	app.Get("/lang/:code", handler.LangJSON("web/locales"))
	// ── Portfolio & resume ─────────────────────────────────────────
	app.Get("/", handler.PortfolioPage(db))
	app.Get("/resume", handler.ServeResumePDF(db))
	app.Get("/resume/download", handler.DownloadResumePDF(db))
	// HTMX partials for project + upcoming pagination/search
	app.Get("/projects", handler.ProjectsPage(db))
	app.Get("/upcoming", handler.UpcomingPage(db))

	// ── OAuth ──────────────────────────────────────────────────────
	app.Get("/auth/google", handler.GoogleLogin())
	app.Get("/auth/google/callback", handler.GoogleCallback(db))
	app.Get("/auth/github", handler.GitHubLogin())
	app.Get("/auth/github/callback", handler.GitHubCallback(db))
	app.Get("/auth/logout", handler.OAuthLogout())

	// ── Comments ───────────────────────────────────────────────────
	app.Get("/comments", cb, handler.GetComments(db))
	app.Post("/comments", commentLimiter, cb, middleware.OAuthAuth(), handler.PostComment(db, h))

	// ── Blog ───────────────────────────────────────────────────────
	app.Get("/blog", handler.BlogListPage(db))
	app.Get("/blog/more", cb, handler.BlogPostsPartial(db))
	app.Get("/blog/:slug", handler.BlogPostPage(db))

	// ── Contact ────────────────────────────────────────────────────
	app.Post("/contact", contactLimiter, cb, handler.SubmitContact(db))
}
