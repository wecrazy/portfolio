package handler

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ServeResumePDF renders a page with an embedded PDF viewer.
func ServeResumePDF(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ResumeFile").First(&owner)
		if owner.ResumeFile == nil {
			return c.Status(fiber.StatusNotFound).SendString("No resume uploaded")
		}
		return c.Render("public/pdf_viewer", fiber.Map{
			"Title":          "Resume",
			"ResumeFilename": owner.ResumeFile.StoredName,
		}, "layouts/public_base")
	}
}

// DownloadResumePDF forces a file download for the resume.
func DownloadResumePDF(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ResumeFile").First(&owner)
		if owner.ResumeFile == nil {
			return c.Status(fiber.StatusNotFound).SendString("No resume uploaded")
		}
		return c.Download(owner.ResumeFile.FilePath, owner.ResumeFile.OriginalName)
	}
}
