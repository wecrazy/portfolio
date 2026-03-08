package admin

import (
	"os"

	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// UploadListPage renders the uploads admin page.
func UploadListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		title, _ := appI18n.T.Localize(c, "admin.uploads.title")
		return c.Render("admin/uploads", fiber.Map{
			"Title":          title,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// UploadListPartial returns the upload cards as an HTMX partial.
func UploadListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "created_at", []string{"created_at", "original_name", "category", "mime_type"})
		if c.Query("sort_dir") == "" {
			params.SortDir = "DESC"
		}
		var items []model.UploadedFile
		query, pageResult := pagination.Paginate(db, &model.UploadedFile{}, params, []string{"original_name", "category", "mime_type"})
		query.Find(&items)
		return c.Render("partials/upload_rows", fiber.Map{
			"Uploads":    items,
			"Pagination": pageResult,
		})
	}
}

// UploadCreate handles uploading a new file.
func UploadCreate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file uploaded")
		}
		category := c.FormValue("category", "images")
		uploaded, err := service.ProcessUpload(file, category)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := c.SaveFile(file, uploaded.FilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}
		db.Create(uploaded)
		setToast(c, "upload_success", "success")
		return c.Render("partials/upload_card", fiber.Map{"Upload": *uploaded})
	}
}

// UploadDelete handles deleting an uploaded file.
func UploadDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.UploadedFile
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		// Remove from disk.
		os.Remove(item.FilePath)
		db.Delete(&item)
		setToast(c, "file_deleted", "success")
		return c.SendString("")
	}
}
