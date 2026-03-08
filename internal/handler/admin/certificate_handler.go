package admin

import (
	"io"
	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"
	"my-portfolio/pkg/pagination"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// CertificateListPage renders the skeleton certificates page.  The
// actual table rows are loaded via HTMX from CertificateListPartial so that
// search, sorting and pagination behave the same as other admin sections.
func CertificateListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		title, _ := appI18n.T.Localize(c, "admin.certificates.title")
		return c.Render("admin/certificates", fiber.Map{
			"Title":          title,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// CertificateListPartial returns the certificates table rows as an HTMX
// partial, including pagination controls.
func CertificateListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "issue_date", []string{"issue_date", "title", "issuer"})
		var certificates []model.Certificate
		query, pageResult := pagination.Paginate(db.Preload("File"), &model.Certificate{}, params, []string{"title", "issuer", "description"})
		query.Order("sort_order ASC, issue_date DESC").Find(&certificates)
		return c.Render("partials/certificate_rows", fiber.Map{
			"Certificates": certificates,
			"Pagination":   pageResult,
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

// DeleteCertificate handles deletion of a certificate.  When invoked via
// HTMX (DELETE) it responds with an empty body so the row can be removed
// client-side; for plain GET requests we fall back to a redirect for
// backwards compatibility.
func DeleteCertificate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if err := db.Delete(&model.Certificate{}, id).Error; err != nil {
			return c.Status(500).SendString("Failed to delete certificate")
		}
		setToast(c, "certificate_deleted", "success")
		if c.Method() == "DELETE" {
			return c.SendString("")
		}
		// legacy GET route
		c.Redirect().To("/admin/certificates")
		return nil
	}
}
