package handler

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// careerStart is the date Wegil started his professional career.
var careerStart = time.Date(2023, time.February, 13, 0, 0, 0, 0, time.UTC)

// expandText replaces dynamic placeholders in owner text fields so the values
// stay accurate over time without any manual edits to the database.
//
//	{years_experience} → "+3 years", "+4 years", etc.
func expandText(s string) string {
	years := int(time.Since(careerStart).Hours() / (24 * 365.25))
	return strings.ReplaceAll(s, "{years_experience}", fmt.Sprintf("+%d years", years))
}

const projectPageSize = 6
const upcomingPageSize = 6

// PortfolioPage renders the public portfolio with all published content.
func PortfolioPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ProfileImage").Preload("ResumeFile").First(&owner)

		// Expand dynamic placeholders in owner text fields.
		owner.Tagline = expandText(owner.Tagline)
		owner.Bio = expandText(owner.Bio)

		// ── Projects (first page) ──────────────────────────────────
		var projects []model.Project
		var totalProjects int64
		db.Model(&model.Project{}).Where("status = ?", "published").Count(&totalProjects)
		db.Where("status = ?", "published").
			Order("sort_order ASC, created_at DESC").
			Limit(projectPageSize).Find(&projects)
		projectTotalPages := int(math.Ceil(float64(totalProjects) / float64(projectPageSize)))

		// ── Upcoming (first page) ──────────────────────────────────
		var upcomingItems []model.UpcomingItem
		var totalUpcoming int64
		db.Model(&model.UpcomingItem{}).Where("is_visible = ?", true).Count(&totalUpcoming)
		db.Where("is_visible = ?", true).
			Order("sort_order ASC, created_at ASC").
			Limit(upcomingPageSize).Find(&upcomingItems)
		upcomingTotalPages := int(math.Ceil(float64(totalUpcoming) / float64(upcomingPageSize)))

		// ── Other data ─────────────────────────────────────────────
		var experiences []model.Experience
		db.Preload("Image").Order("sort_order ASC, start_date DESC").Find(&experiences)

		var skills []model.Skill
		db.Order("category ASC, sort_order ASC").Find(&skills)

		var socialLinks []model.SocialLink
		db.Order("sort_order ASC").Find(&socialLinks)

		var techStacks []model.TechStack
		db.Order("category ASC, sort_order ASC").Find(&techStacks)

		skillsByCategory := make(map[string][]model.Skill)
		for _, s := range skills {
			skillsByCategory[s.Category] = append(skillsByCategory[s.Category], s)
		}

		techByCategory := make(map[string][]model.TechStack)
		for _, t := range techStacks {
			techByCategory[t.Category] = append(techByCategory[t.Category], t)
		}

		visitorLoggedIn := c.Cookies("visitor_session") != ""

		// Certificates (first page, min 6, paginated, searchable)
		var certificates []model.Certificate
		var totalCertificates int64
		db.Model(&model.Certificate{}).Where("is_visible = ?", true).Count(&totalCertificates)
		db.Where("is_visible = ?", true).
			Preload("File").
			Order("sort_order ASC, issue_date DESC").
			Limit(6).Find(&certificates)
		certificateTotalPages := int(math.Ceil(float64(totalCertificates) / 6.0))

		cfg := config.MyPortfolio.Get()
		ogImage := cfg.App.BaseURL + "/static/img/favicon.svg"
		if owner.ProfileImage != nil {
			ogImage = cfg.App.BaseURL + "/uploads/images/" + owner.ProfileImage.StoredName
		}
		ogDesc := owner.Bio
		if len(ogDesc) > 160 {
			ogDesc = ogDesc[:157] + "..."
		}
		if ogDesc == "" {
			ogDesc = owner.Title
		}

		return c.Render("public/portfolio", fiber.Map{
			"Title":                 owner.FullName,
			"BaseURL":               cfg.App.BaseURL,
			"OGImage":               ogImage,
			"OGDescription":         ogDesc,
			"Owner":                 owner,
			"Projects":              projects,
			"ProjectCurrentPage":    1,
			"ProjectTotalPages":     projectTotalPages,
			"Experiences":           experiences,
			"Skills":                skills,
			"SkillsByCategory":      skillsByCategory,
			"TechByCategory":        techByCategory,
			"SocialLinks":           socialLinks,
			"UpcomingItems":         upcomingItems,
			"UpcomingCurrentPage":   1,
			"UpcomingTotalPages":    upcomingTotalPages,
			"Certificates":          certificates,
			"CertificateTotalPages": certificateTotalPages,
			"VisitorLoggedIn":       visitorLoggedIn,
			"SupportedLangs":        cfg.I18n.SupportedLangs,
			"DefaultLang":           cfg.I18n.DefaultLang,
			"IsPortfolio":           true,
			"HCaptchaEnabled":       cfg.HCaptcha.Enabled,
			"HCaptchaKey":           cfg.HCaptcha.SiteKey,
		}, "layouts/public_base")
	}
}

// ProjectsPage returns a paginated, searchable project grid as an HTMX partial.
func ProjectsPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		q := strings.TrimSpace(c.Query("q", ""))
		offset := (page - 1) * projectPageSize

		query := db.Model(&model.Project{}).Where("status = ?", "published")
		if q != "" {
			like := "%" + q + "%"
			query = query.Where("title LIKE ? OR tags LIKE ? OR description LIKE ?", like, like, like)
		}

		var total int64
		query.Count(&total)

		var projects []model.Project
		query.Order("sort_order ASC, created_at DESC").
			Offset(offset).Limit(projectPageSize).
			Find(&projects)

		totalPages := int(math.Ceil(float64(total) / float64(projectPageSize)))

		return c.Render("partials/project_cards", fiber.Map{
			"Projects":    projects,
			"CurrentPage": page,
			"TotalPages":  totalPages,
		})
	}
}

// UpcomingPage returns a paginated, searchable upcoming-items grid as an HTMX partial.
func UpcomingPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		q := strings.TrimSpace(c.Query("q", ""))
		offset := (page - 1) * upcomingPageSize

		query := db.Model(&model.UpcomingItem{}).Where("is_visible = ?", true)
		if q != "" {
			like := "%" + q + "%"
			query = query.Where("title LIKE ? OR description LIKE ?", like, like)
		}

		var total int64
		query.Count(&total)

		var upcomingItems []model.UpcomingItem
		query.Order("sort_order ASC, created_at ASC").
			Offset(offset).Limit(upcomingPageSize).
			Find(&upcomingItems)

		totalPages := int(math.Ceil(float64(total) / float64(upcomingPageSize)))

		return c.Render("partials/upcoming_cards", fiber.Map{
			"UpcomingItems": upcomingItems,
			"CurrentPage":   page,
			"TotalPages":    totalPages,
		})
	}
}
