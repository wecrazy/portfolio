package admin

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// UpcomingListPage renders the upcoming items admin page.
func UpcomingListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/upcoming", fiber.Map{
			"Title":          "Upcoming",
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// UpcomingListPartial returns the upcoming items table rows as an HTMX partial.
func UpcomingListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "sort_order", []string{"sort_order", "title", "type", "status", "created_at"})
		var items []model.UpcomingItem
		query, pageResult := pagination.Paginate(db, &model.UpcomingItem{}, params, []string{"title", "description", "type"})
		query.Find(&items)
		return c.Render("partials/upcoming_rows", fiber.Map{
			"Items":      items,
			"Pagination": pageResult,
		})
	}
}

// UpcomingNewForm renders an empty upcoming item form partial.
func UpcomingNewForm() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Render("partials/upcoming_form", fiber.Map{"Item": model.UpcomingItem{}})
	}
}

// UpcomingEditForm renders a pre-filled upcoming item form partial.
func UpcomingEditForm(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.UpcomingItem
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Item not found")
		}
		return c.Render("partials/upcoming_form", fiber.Map{"Item": item})
	}
}

// UpcomingCreate handles creating a new upcoming item.
func UpcomingCreate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.UpcomingItem
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsVisible = c.FormValue("is_visible") == "on"
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create item")
		}
		setToast(c, "upcoming_created", "success")
		return c.Render("partials/upcoming_row", fiber.Map{"Item": item})
	}
}

// UpcomingUpdate handles updating an existing upcoming item.
func UpcomingUpdate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.UpcomingItem
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Item not found")
		}
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsVisible = c.FormValue("is_visible") == "on"
		db.Save(&item)
		setToast(c, "upcoming_updated", "success")
		return c.Render("partials/upcoming_row", fiber.Map{"Item": item})
	}
}

// UpcomingDelete handles deleting an upcoming item.
func UpcomingDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := db.Delete(&model.UpcomingItem{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		setToast(c, "upcoming_deleted", "success")
		return c.SendString("")
	}
}
