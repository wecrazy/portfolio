// Package admin provides HTTP handlers for the admin dashboard.
package admin

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

var commentSanitizer = bluemonday.StrictPolicy()

// CommentListPage renders the comments moderation admin page.
func CommentListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/comments", fiber.Map{
			"Title": "Comments",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// CommentListPartial returns the comment cards as an HTMX partial.
func CommentListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var comments []model.Comment
		db.Where("parent_id IS NULL").
			Preload("OAuthUser").
			Preload("Replies", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("OAuthUser").Order("created_at ASC")
			}).
			Order("created_at DESC").
			Find(&comments)
		return c.Render("partials/admin_comment_rows", fiber.Map{"Comments": comments})
	}
}

// CommentApprove marks a comment as approved.
func CommentApprove(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Model(&model.Comment{}).Where("id = ?", c.Params("id")).Update("is_approved", true)
		c.Set("HX-Trigger", `{"showToast":"Comment approved"}`)
		return c.SendString(`<span class="badge bg-success">Approved</span>`)
	}
}

// CommentReject marks a comment as rejected.
func CommentReject(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Model(&model.Comment{}).Where("id = ?", c.Params("id")).Update("is_approved", false)
		c.Set("HX-Trigger", `{"showToast":"Comment rejected"}`)
		return c.SendString(`<span class="badge bg-warning">Rejected</span>`)
	}
}

// CommentDelete handles deleting a comment and its replies.
func CommentDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Delete replies first, then the comment.
		db.Where("parent_id = ?", c.Params("id")).Delete(&model.Comment{})
		db.Delete(&model.Comment{}, c.Params("id"))
		c.Set("HX-Trigger", `{"showToast":"Comment deleted"}`)
		return c.SendString("")
	}
}

// CommentReply handles posting an owner reply to a comment.
func CommentReply(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		parentID := c.Params("id")
		body := commentSanitizer.Sanitize(c.FormValue("body"))
		if body == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Reply body is required")
		}

		// Find or create a system OAuthUser for owner replies.
		var ownerOAuth model.OAuthUser
		db.Where("provider = ? AND provider_id = ?", "system", "owner").First(&ownerOAuth)
		if ownerOAuth.ID == 0 {
			ownerOAuth = model.OAuthUser{
				Provider:    "system",
				ProviderID:  "owner",
				DisplayName: "Owner",
			}
			db.Create(&ownerOAuth)
		}

		var pid uint
		if _, err := parseUint(parentID); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid comment ID")
		}
		pid = parseUintVal(parentID)

		reply := model.Comment{
			OAuthUserID:  ownerOAuth.ID,
			ParentID:     &pid,
			Body:         body,
			IsOwnerReply: true,
			IsApproved:   true,
		}
		db.Create(&reply)
		db.Preload("OAuthUser").First(&reply, reply.ID)

		c.Set("HX-Trigger", `{"showToast":"Reply posted"}`)
		return c.Render("partials/comment_reply", fiber.Map{"Reply": reply})
	}
}

func parseUint(s string) (uint, error) {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fiber.ErrBadRequest
		}
		n = n*10 + uint(c-'0')
	}
	return n, nil
}

func parseUintVal(s string) uint {
	n, _ := parseUint(s)
	return n
}
