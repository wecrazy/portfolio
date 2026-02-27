// Package i18n exposes a global, thread-safe translator backed by
// github.com/gofiber/contrib/v3/i18n (go-i18n v2).
//
// Usage in a handler:
//
//	msg, _ := i18n.T.Localize(c, "comment_submitted")
package i18n

import (
	contribi18n "github.com/gofiber/contrib/v3/i18n"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/text/language"
)

// T is the global translator instance. It is safe for concurrent use.
var T *contribi18n.I18n

func init() {
	T = contribi18n.New(&contribi18n.Config{
		RootPath:         "./web/locales",
		FormatBundleFile: "yaml",
		DefaultLanguage:  language.English,
		AcceptLanguages:  []language.Tag{language.English, language.Indonesian},
		// LangHandler resolves the request language in priority order:
		//  1. "lang" cookie — set by the JS language switcher (persists across HTMX requests)
		//  2. ?lang= query param
		//  3. Accept-Language header (first 2-char tag)
		//  4. Configured default language
		LangHandler: func(c fiber.Ctx, defaultLang string) string {
			if lang := c.Cookies("lang"); lang != "" {
				return lang
			}
			if lang := c.Query("lang"); lang != "" {
				return lang
			}
			if al := c.Get("Accept-Language"); len(al) >= 2 {
				return al[:2]
			}
			return defaultLang
		},
	})
}
