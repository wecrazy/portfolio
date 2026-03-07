package admin

import (
	"io"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ListCertificates handles admin listing of certificates.
func ListCertificates(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var certificates []model.Certificate
		db.Preload("File").Order("sort_order ASC, issue_date DESC").Find(&certificates)
		return c.Render("admin/certificates", fiber.Map{
			"Certificates": certificates,
		})
	}
}

// CertificateForm handles add/edit form rendering.  It does not access the database;
// db parameter was previously included for consistency with other handlers but is
// unused and therefore removed to satisfy linters.
func CertificateForm(cert *model.Certificate, action string) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Render("admin/certificate_form", fiber.Map{
			"Certificate": cert,
			"FormAction":  action,
		})
	}
}

// CreateCertificate handles creation of a new certificate.
func CreateCertificate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		cert := model.Certificate{
			Title:       c.FormValue("title"),
			Issuer:      c.FormValue("issuer"),
			Description: c.FormValue("description"),
			CertURL:     c.FormValue("cert_url"),
			IsVisible:   c.FormValue("is_visible") == "on",
			SortOrder:   0,
		}
		if dateStr := c.FormValue("issue_date"); dateStr != "" {
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				cert.IssueDate = t
			}
		}
		fileHeader, err := c.FormFile("file")
		if err == nil && fileHeader != nil {
			uploaded, err := service.ProcessUpload(fileHeader, "certificates")
			if err != nil {
				return c.Status(400).SendString("File upload error: " + err.Error())
			}
			file, err := fileHeader.Open()
			if err != nil {
				return c.Status(400).SendString("File open error: " + err.Error())
			}
			defer file.Close()
			out, err := os.Create(uploaded.FilePath)
			if err != nil {
				return c.Status(500).SendString("Failed to save file")
			}
			defer out.Close()
			if _, err := io.Copy(out, file); err != nil {
				return c.Status(500).SendString("Failed to write file")
			}
			if err := db.Create(uploaded).Error; err == nil {
				cert.FileID = &uploaded.ID
			}
		}
		if err := db.Create(&cert).Error; err != nil {
			return c.Status(500).SendString("Failed to create certificate")
		}
		c.Redirect().To("/admin/certificates")
		return nil
	}
}

// EditCertificate handles editing an existing certificate.
func EditCertificate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		var cert model.Certificate
		if err := db.First(&cert, id).Error; err != nil {
			return c.Status(404).SendString("Certificate not found")
		}
		cert.Title = c.FormValue("title")
		cert.Issuer = c.FormValue("issuer")
		cert.Description = c.FormValue("description")
		cert.CertURL = c.FormValue("cert_url")
		cert.IsVisible = c.FormValue("is_visible") == "on"
		if dateStr := c.FormValue("issue_date"); dateStr != "" {
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				cert.IssueDate = t
			}
		}
		fileHeader, err := c.FormFile("file")
		if err == nil && fileHeader != nil {
			uploaded, err := service.ProcessUpload(fileHeader, "certificates")
			if err != nil {
				return c.Status(400).SendString("File upload error: " + err.Error())
			}
			file, err := fileHeader.Open()
			if err != nil {
				return c.Status(400).SendString("File open error: " + err.Error())
			}
			defer file.Close()
			out, err := os.Create(uploaded.FilePath)
			if err != nil {
				return c.Status(500).SendString("Failed to save file")
			}
			defer out.Close()
			if _, err := io.Copy(out, file); err != nil {
				return c.Status(500).SendString("Failed to write file")
			}
			if err := db.Create(uploaded).Error; err == nil {
				cert.FileID = &uploaded.ID
			}
		}
		if err := db.Save(&cert).Error; err != nil {
			return c.Status(500).SendString("Failed to update certificate")
		}
		c.Redirect().To("/admin/certificates")
		return nil
	}
}

// DeleteCertificate handles deletion of a certificate.
func DeleteCertificate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if err := db.Delete(&model.Certificate{}, id).Error; err != nil {
			return c.Status(500).SendString("Failed to delete certificate")
		}
		c.Redirect().To("/admin/certificates")
		return nil
	}
}
