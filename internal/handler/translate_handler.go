package handler

import (
	"my-portfolio/pkg/translate"

	"github.com/gofiber/fiber/v3"
)

// translateRequest is the JSON body expected by TranslateText.
type translateRequest struct {
	Texts    []string `json:"texts"`
	LangFrom string   `json:"from"`
	LangTo   string   `json:"to"`
}

// translateResponse is the JSON body returned by TranslateText.
type translateResponse struct {
	Translations []string `json:"translations"`
}

// TranslateText translates an array of strings from one language to another
// using the provided Translator.
//
// POST /api/translate
//
//	{"texts": ["Hello", "World"], "from": "en", "to": "id"}
//	→ {"translations": ["Halo", "Dunia"]}
func TranslateText(t translate.Translator) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req translateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		}

		if len(req.Texts) == 0 {
			return c.JSON(translateResponse{Translations: []string{}})
		}
		if req.LangFrom == "" {
			req.LangFrom = "en"
		}
		if req.LangTo == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing 'to' field"})
		}

		// If source and target are the same, echo back immediately.
		if req.LangFrom == req.LangTo {
			return c.JSON(translateResponse{Translations: req.Texts})
		}

		translated, err := t.TranslateMany(c.Context(), req.Texts, req.LangFrom, req.LangTo)
		if err != nil {
			// Partial failure: TranslateMany already fills the slice with original
			// text on errors, so we can still return what we have.
			return c.JSON(translateResponse{Translations: translated})
		}

		return c.JSON(translateResponse{Translations: translated})
	}
}
