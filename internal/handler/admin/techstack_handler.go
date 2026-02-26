package admin

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// TechStackListPage renders the tech stack admin page.
func TechStackListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/tech_stacks", fiber.Map{
			"Title": "Tech Stack",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// TechStackListPartial returns the tech stack table rows as an HTMX partial.
func TechStackListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.TechStack
		db.Order("category ASC, sort_order ASC").Find(&items)
		return c.Render("partials/techstack_rows", fiber.Map{"TechStacks": items})
	}
}

// TechStackNewForm renders an empty tech stack form partial.
func TechStackNewForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("partials/techstack_form", fiber.Map{"TechStack": model.TechStack{}})
	}
}

// TechStackEditForm renders a pre-filled tech stack form partial.
func TechStackEditForm(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.TechStack
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		return c.Render("partials/techstack_form", fiber.Map{"TechStack": item})
	}
}

// TechStackCreate handles creating a new tech stack entry.
func TechStackCreate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.TechStack
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create")
		}
		c.Set("HX-Trigger", `{"showToast":"Tech stack added"}`)
		return c.Render("partials/techstack_row", fiber.Map{"Tech": item})
	}
}

// TechStackUpdate handles updating an existing tech stack entry.
func TechStackUpdate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.TechStack
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		db.Save(&item)
		c.Set("HX-Trigger", `{"showToast":"Tech stack updated"}`)
		return c.Render("partials/techstack_row", fiber.Map{"Tech": item})
	}
}

// TechStackDelete handles deleting a tech stack entry.
func TechStackDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Delete(&model.TechStack{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		c.Set("HX-Trigger", `{"showToast":"Tech stack deleted"}`)
		return c.SendString("")
	}
}
