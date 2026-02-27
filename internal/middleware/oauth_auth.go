package middleware

import (
	"sync"

	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
)

// VisitorSessions maps session token → OAuthUser ID.
// Kept in-memory; visitors re-authenticate on restart.
var VisitorSessions sync.Map

// OAuthAuth checks for a valid visitor OAuth session cookie and loads the
// OAuthUser into c.Locals("oauth_user"). Returns 401 when not authenticated.
func OAuthAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("visitor_session")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Login required to perform this action",
			})
		}

		userID, ok := VisitorSessions.Load(token)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session expired, please login again",
			})
		}

		c.Locals("oauth_user_id", userID.(uint))
		return c.Next()
	}
}

// SetVisitorSession stores a session token for an OAuth user.
func SetVisitorSession(token string, user model.OAuthUser) {
	VisitorSessions.Store(token, user.ID)
}

// DeleteVisitorSession removes a session token.
func DeleteVisitorSession(token string) {
	VisitorSessions.Delete(token)
}
