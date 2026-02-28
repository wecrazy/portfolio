package handler

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ServeResumePDF renders a page with an embedded PDF viewer.
func ServeResumePDF(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ResumeFile").First(&owner)
		if owner.ResumeFile == nil {
			return c.Status(fiber.StatusNotFound).SendString("No resume uploaded")
		}
		cfg := config.MyPortfolio.Get()
		return c.Render("public/pdf_viewer", fiber.Map{
			"Title":          "Resume",
			"ResumeFilename": owner.ResumeFile.StoredName,
			"Owner":          owner,
			"BaseURL":        cfg.App.BaseURL,
			"DefaultLang":    cfg.I18n.DefaultLang,
			"SupportedLangs": cfg.I18n.SupportedLangs,
		}, "layouts/public_base")
	}
}

// DownloadResumePDF forces a file download for the resume.
func DownloadResumePDF(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ResumeFile").First(&owner)
		if owner.ResumeFile == nil {
			return c.Status(fiber.StatusNotFound).SendString("No resume uploaded")
		}
		return c.Download(owner.ResumeFile.FilePath, owner.ResumeFile.OriginalName)
	}
}
