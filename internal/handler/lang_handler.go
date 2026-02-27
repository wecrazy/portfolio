package handler

import (
	"os"
	"path/filepath"

	goyaml "github.com/goccy/go-yaml"
	"github.com/gofiber/fiber/v3"
)

// LangJSON serves a locale YAML file as a flat JSON object so the
// frontend can use a single source of truth at web/locales/*.yaml.
//
//	GET /lang/en  →  {"key": "value", ...}
//	GET /lang/id  →  {"key": "valeur", ...}
func LangJSON(localesDir string) fiber.Handler {
	allowed := map[string]bool{"en": true, "id": true}
	return func(c fiber.Ctx) error {
		code := c.Params("code")
		if !allowed[code] {
			return c.Status(fiber.StatusNotFound).SendString("locale not found")
		}
		data, err := os.ReadFile(filepath.Join(localesDir, code+".yaml"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("locale not found")
		}
		var m map[string]string
		if err := goyaml.Unmarshal(data, &m); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("parse error")
		}
		return c.JSON(m)
	}
}
