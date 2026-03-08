package admin

import (
	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ContactListPage renders the contact messages admin page.
func ContactListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		title, _ := appI18n.T.Localize(c, "admin.contact.title")
		return c.Render("admin/contacts", fiber.Map{
			"Title":          title,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// ContactListPartial returns the contact message rows as an HTMX partial.
func ContactListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "created_at", []string{"created_at", "name", "email", "subject"})
		if c.Query("sort_dir") == "" {
			params.SortDir = "DESC"
		}
		var items []model.ContactMessage
		query, pageResult := pagination.Paginate(db, &model.ContactMessage{}, params, []string{"name", "email", "subject", "message"})
		query.Find(&items)
		return c.Render("partials/contact_rows", fiber.Map{
			"Contacts":   items,
			"Pagination": pageResult,
		})
	}
}

// ContactMarkRead marks a contact message as read.
func ContactMarkRead(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		db.Model(&model.ContactMessage{}).Where("id = ?", c.Params("id")).Update("is_read", true)
		setToast(c, "contact_read", "info")
		return c.SendString(`<span class="badge bg-secondary">Read</span>`)
	}
}

// ContactDelete handles deleting a contact message.
func ContactDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		db.Delete(&model.ContactMessage{}, c.Params("id"))
		setToast(c, "contact_deleted", "success")
		return c.SendString("")
	}
}
