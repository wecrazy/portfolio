package admin

import (
	"encoding/json"
	"fmt"

	appI18n "my-portfolio/internal/i18n"

	"github.com/gofiber/fiber/v3"
)

// setToast fires a showToast event on the client via the HX-Trigger response header.
//
// msgID is looked up in the YAML message bundle for the request language (resolved
// from the "lang" cookie first, then Accept-Language, then the configured default).
// kind is one of "success", "info", "warning", or "danger".
//
// The generated header format is:
//
//	HX-Trigger: {"showToast":{"message":"...","type":"..."}}
//
// This matches the showToast listener in admin.js which reads e.detail.message and
// e.detail.type.
func setToast(c fiber.Ctx, msgID, kind string) {
	msg, err := appI18n.T.Localize(c, msgID)
	if err != nil || msg == "" {
		msg = msgID // visible fallback — key itself, not silent empty
	}

	type payload struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}
	b, _ := json.Marshal(payload{Message: msg, Type: kind})
	c.Set("HX-Trigger", fmt.Sprintf(`{"showToast":%s}`, b))
}
