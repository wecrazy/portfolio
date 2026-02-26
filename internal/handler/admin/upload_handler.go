package admin

import (
	"os"

	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UploadListPage renders the uploads admin page.
func UploadListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/uploads", fiber.Map{
			"Title": "Uploads",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// UploadListPartial returns the upload cards as an HTMX partial.
func UploadListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.UploadedFile
		db.Order("created_at DESC").Find(&items)
		return c.Render("partials/upload_rows", fiber.Map{"Uploads": items})
	}
}

// UploadCreate handles uploading a new file.
func UploadCreate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		c.Set("HX-Trigger", `{"showToast":"File uploaded"}`)
		return c.Render("partials/upload_card", fiber.Map{"Upload": *uploaded})
	}
}

// UploadDelete handles deleting an uploaded file.
func UploadDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item model.UploadedFile
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		// Remove from disk.
		os.Remove(item.FilePath)
		db.Delete(&item)
		c.Set("HX-Trigger", `{"showToast":"File deleted"}`)
		return c.SendString("")
	}
}
