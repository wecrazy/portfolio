package admin

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ContactListPage renders the contact messages admin page.
func ContactListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/contacts", fiber.Map{
			"Title": "Contact Messages",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// ContactListPartial returns the contact message rows as an HTMX partial.
func ContactListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.ContactMessage
		db.Order("created_at DESC").Find(&items)
		return c.Render("partials/contact_rows", fiber.Map{"Contacts": items})
	}
}

// ContactMarkRead marks a contact message as read.
func ContactMarkRead(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Model(&model.ContactMessage{}).Where("id = ?", c.Params("id")).Update("is_read", true)
		c.Set("HX-Trigger", `{"showToast":"Marked as read"}`)
		return c.SendString(`<span class="badge bg-secondary">Read</span>`)
	}
}

// ContactDelete handles deleting a contact message.
func ContactDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Delete(&model.ContactMessage{}, c.Params("id"))
		c.Set("HX-Trigger", `{"showToast":"Message deleted"}`)
		return c.SendString("")
	}
}
