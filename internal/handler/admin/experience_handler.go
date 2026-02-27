package admin

import (
	"fmt"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ExperienceListPage renders the experience admin page.
func ExperienceListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/experience", fiber.Map{
			"Title":          "Experience",
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// ExperienceListPartial returns the experience table rows as an HTMX partial.
func ExperienceListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "sort_order", []string{"sort_order", "title", "org", "type", "start_date"})
		var items []model.Experience
		query, pageResult := pagination.Paginate(db, &model.Experience{}, params, []string{"title", "org", "location", "type"})
		query.Preload("Image").Find(&items)
		return c.Render("partials/experience_rows", fiber.Map{
			"Experiences": items,
			"Pagination":  pageResult,
		})
	}
}

// ExperienceNewForm renders an empty experience form partial.
func ExperienceNewForm() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Render("partials/experience_form", fiber.Map{"Experience": model.Experience{}})
	}
}

// ExperienceEditForm renders a pre-filled experience form partial.
func ExperienceEditForm(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Experience
		if err := db.Preload("Image").First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		return c.Render("partials/experience_form", fiber.Map{"Experience": item})
	}
}

// ExperienceCreate handles creating a new experience entry.
func ExperienceCreate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Experience
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsCurrent = c.FormValue("is_current") == "on"
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create")
		}
		setToast(c, "experience_created", "success")
		return c.Render("partials/experience_row", fiber.Map{"Experience": item})
	}
}

// ExperienceUpdate handles updating an existing experience entry.
func ExperienceUpdate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Experience
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		item.IsCurrent = c.FormValue("is_current") == "on"
		db.Save(&item)
		setToast(c, "experience_updated", "success")
		return c.Render("partials/experience_row", fiber.Map{"Experience": item})
	}
}

// ExperienceDelete handles deleting an experience entry.
func ExperienceDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := db.Delete(&model.Experience{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		setToast(c, "experience_deleted", "success")
		return c.SendString("")
	}
}

// ExperienceUploadImage handles uploading a thumbnail image for an experience entry.
func ExperienceUploadImage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file uploaded")
		}
		uploaded, err := service.ProcessUpload(file, "images")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := c.SaveFile(file, uploaded.FilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}
		db.Create(uploaded)

		setToast(c, "image_uploaded", "success")
		html := fmt.Sprintf(
			`<img src="/uploads/images/%s" class="img-thumbnail mt-2" style="max-height:120px" alt="Experience Image"><input type="hidden" id="experience_image_id" name="image_id" value="%d" hx-swap-oob="outerHTML:#experience_image_id">`,
			uploaded.StoredName, uploaded.ID,
		)
		return c.SendString(html)
	}
}
