package admin

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ExperienceListPage renders the experience admin page.
func ExperienceListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/experience", fiber.Map{
			"Title": "Experience",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// ExperienceListPartial returns the experience table rows as an HTMX partial.
func ExperienceListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.Experience
		db.Order("sort_order ASC, start_date DESC").Find(&items)
		return c.Render("partials/experience_rows", fiber.Map{"Experiences": items})
	}
}

// ExperienceNewForm renders an empty experience form partial.
func ExperienceNewForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("partials/experience_form", fiber.Map{"Experience": model.Experience{}})
	}
}

// ExperienceEditForm renders a pre-filled experience form partial.
func ExperienceEditForm(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.Experience
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		return c.Render("partials/experience_form", fiber.Map{"Experience": item})
	}
}

// ExperienceCreate handles creating a new experience entry.
func ExperienceCreate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.Experience
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsCurrent = c.FormValue("is_current") == "on"
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create")
		}
		c.Set("HX-Trigger", `{"showToast":"Experience created"}`)
		return c.Render("partials/experience_row", fiber.Map{"Experience": item})
	}
}

// ExperienceUpdate handles updating an existing experience entry.
func ExperienceUpdate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.Experience
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsCurrent = c.FormValue("is_current") == "on"
		db.Save(&item)
		c.Set("HX-Trigger", `{"showToast":"Experience updated"}`)
		return c.Render("partials/experience_row", fiber.Map{"Experience": item})
	}
}

// ExperienceDelete handles deleting an experience entry.
func ExperienceDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Delete(&model.Experience{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		c.Set("HX-Trigger", `{"showToast":"Experience deleted"}`)
		return c.SendString("")
	}
}
