package handler

import (
	"math"
	"my-portfolio/internal/model"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

const certificatePageSize = 6

// CertificatesPage returns a paginated, searchable certificate grid as an HTMX partial.
func CertificatesPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		q := strings.TrimSpace(c.Query("q", ""))
		offset := (page - 1) * certificatePageSize

		query := db.Model(&model.Certificate{}).Where("is_visible = ?", true)
		if q != "" {
			like := "%" + q + "%"
			query = query.Where("title LIKE ? OR issuer LIKE ? OR description LIKE ?", like, like, like)
		}

		var certificates []model.Certificate
		var totalCertificates int64
		query.Count(&totalCertificates)
		query.Preload("File").Order("sort_order ASC, issue_date DESC").Offset(offset).Limit(certificatePageSize).Find(&certificates)
		certificateTotalPages := int(math.Ceil(float64(totalCertificates) / float64(certificatePageSize)))

		return c.Render("public/certificate_partial", fiber.Map{
			"Certificates":           certificates,
			"CertificateTotalPages":  certificateTotalPages,
			"CertificateCurrentPage": page,
		})
	}
}
