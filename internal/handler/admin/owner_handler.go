package admin

import (
	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// OwnerEditPage renders the owner profile edit form.
func OwnerEditPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ProfileImage").Preload("ResumeFile").First(&owner)
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/owner", fiber.Map{
			"Title":          "Owner Profile",
			"Owner":          owner,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// OwnerUpdate saves owner profile changes.
func OwnerUpdate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.First(&owner)

		owner.FullName = c.FormValue("full_name")
		owner.Title = c.FormValue("title")
		owner.Bio = c.FormValue("bio")
		owner.Email = c.FormValue("email")
		owner.Phone = c.FormValue("phone")
		owner.Location = c.FormValue("location")

		db.Save(&owner)

		setToast(c, "owner_updated", "success")
		db.Preload("ProfileImage").Preload("ResumeFile").First(&owner)
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/owner", fiber.Map{
			"Title":          "Owner Profile",
			"Owner":          owner,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// OwnerUploadImage handles profile image upload.
func OwnerUploadImage(db *gorm.DB) fiber.Handler {
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

		var owner model.Owner
		db.First(&owner)
		db.Model(&owner).Update("profile_image_id", uploaded.ID)

		setToast(c, "owner_image_updated", "success")
		return c.SendString(`<img src="/uploads/images/` + uploaded.StoredName + `" class="img-fluid rounded" alt="Profile">`)
	}
}

// OwnerUploadResume handles resume PDF upload.
func OwnerUploadResume(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("resume")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file uploaded")
		}

		uploaded, err := service.ProcessUpload(file, "resume")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := c.SaveFile(file, uploaded.FilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}
		db.Create(uploaded)

		var owner model.Owner
		db.First(&owner)
		db.Model(&owner).Update("resume_file_id", uploaded.ID)

		setToast(c, "resume_uploaded", "success")
		return c.SendString(`<a href="/resume" target="_blank" class="btn btn-sm btn-outline-primary">` + uploaded.OriginalName + `</a>`)
	}
}
