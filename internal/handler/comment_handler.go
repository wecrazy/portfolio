package handler

import (
	"my-portfolio/internal/hub"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/sanitize"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// GetComments returns the comments list as an HTML partial.
func GetComments(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var comments []model.Comment
		db.Where("parent_id IS NULL AND is_approved = ?", true).
			Preload("OAuthUser").
			Preload("Replies", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("OAuthUser").Order("created_at ASC")
			}).
			Order("created_at DESC").
			Find(&comments)

		return c.Render("partials/comment_list", fiber.Map{
			"Comments": comments,
		})
	}
}

// PostComment creates a new comment (requires OAuth session via middleware).
// On success it broadcasts a "comment" SSE event so live clients update.
func PostComment(db *gorm.DB, h *hub.Hub) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("oauth_user_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authenticated"})
		}

		body := sanitize.Strict(c.FormValue("body"))
		if body == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Comment body is required")
		}

		comment := model.Comment{
			OAuthUserID: userID,
			Body:        body,
			IsApproved:  true,
		}
		if err := db.Create(&comment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to post comment")
		}

		// Reload with association for rendering.
		db.Preload("OAuthUser").First(&comment, comment.ID)

		// Notify all SSE clients that a new comment is available.
		h.Broadcast(hub.Event{
			Type: hub.EventComment,
			Data: map[string]any{
				"id":      comment.ID,
				"user":    comment.OAuthUser.DisplayName,
				"refresh": true,
			},
		})

		return c.Render("partials/comment_card", fiber.Map{
			"Comment": comment,
		})
	}
}
