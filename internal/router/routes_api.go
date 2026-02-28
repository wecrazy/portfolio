package router

import (
	"my-portfolio/internal/handler"
	"my-portfolio/pkg/translate"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// registerAPIRoutes mounts all /api/* endpoints.
//
// Keep each logical group in its own block comment so new routes
// are easy to add later without touching other files.
func registerAPIRoutes(
	app *fiber.App,
	db *gorm.DB,
	cb fiber.Handler,
	t translate.Translator,
) {
	api := app.Group("/api")

	// ── Translation ────────────────────────────────────────────────
	// POST /api/translate  — machine-translate dynamic DB content on the client.
	api.Post("/translate", handler.TranslateText(t))

	// Future groups go here, e.g.:
	// api.Get("/posts",    handler.APIPostList(db))
	// api.Get("/projects", handler.APIProjectList(db))
}
