package admin

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Dashboard renders the admin dashboard overview page.
func Dashboard(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectCount, commentCount, contactCount, unreadCount int64
		db.Model(&model.Project{}).Count(&projectCount)
		db.Model(&model.Comment{}).Where("parent_id IS NULL").Count(&commentCount)
		db.Model(&model.ContactMessage{}).Count(&contactCount)
		db.Model(&model.ContactMessage{}).Where("is_read = ?", false).Count(&unreadCount)

		cfg := config.MyPortfolio.Get()
		return c.Render("admin/dashboard", fiber.Map{
			"Title":          "Dashboard",
			"ProjectCount":   projectCount,
			"CommentCount":   commentCount,
			"ContactCount":   contactCount,
			"UnreadCount":    unreadCount,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}
