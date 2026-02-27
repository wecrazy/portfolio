package admin

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SocialListPage renders the social links admin page.
func SocialListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/social_links", fiber.Map{
			"Title":          "Social Links",
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// SocialListPartial returns the social links table rows as an HTMX partial.
func SocialListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := pagination.ParseParams(c, "sort_order", []string{"sort_order", "platform", "label"})
		var items []model.SocialLink
		query, pageResult := pagination.Paginate(db, &model.SocialLink{}, params, []string{"platform", "label", "url"})
		query.Find(&items)
		return c.Render("partials/social_rows", fiber.Map{
			"SocialLinks": items,
			"Pagination":  pageResult,
		})
	}
}

// SocialNewForm renders an empty social link form partial.
func SocialNewForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("partials/social_form", fiber.Map{"Social": model.SocialLink{}})
	}
}

// SocialEditForm renders a pre-filled social link form partial.
func SocialEditForm(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.SocialLink
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		return c.Render("partials/social_form", fiber.Map{"Social": item})
	}
}

// SocialCreate handles creating a new social link.
func SocialCreate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.SocialLink
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create")
		}
		c.Set("HX-Trigger", `{"showToast":"Social link created"}`)
		return c.Render("partials/social_row", fiber.Map{"Social": item})
	}
}

// SocialUpdate handles updating an existing social link.
func SocialUpdate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.SocialLink
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		db.Save(&item)
		c.Set("HX-Trigger", `{"showToast":"Social link updated"}`)
		return c.Render("partials/social_row", fiber.Map{"Social": item})
	}
}

// SocialDelete handles deleting a social link.
func SocialDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Delete(&model.SocialLink{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		c.Set("HX-Trigger", `{"showToast":"Social link deleted"}`)
		return c.SendString("")
	}
}
