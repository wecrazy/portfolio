package handler

import (
	"strconv"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// PortfolioPage renders the public portfolio with all published content.
func PortfolioPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var owner model.Owner
		db.Preload("ProfileImage").Preload("ResumeFile").First(&owner)

		// Load first page of projects (6 items).
		const pageSize = 6
		var projects []model.Project
		var totalProjects int64
		db.Model(&model.Project{}).Where("status = ?", "published").Count(&totalProjects)
		db.Where("status = ?", "published").Order("sort_order ASC, created_at DESC").Limit(pageSize).Find(&projects)

		var experiences []model.Experience
		db.Preload("Image").Order("sort_order ASC, start_date DESC").Find(&experiences)

		var skills []model.Skill
		db.Order("category ASC, sort_order ASC").Find(&skills)

		var socialLinks []model.SocialLink
		db.Order("sort_order ASC").Find(&socialLinks)

		var techStacks []model.TechStack
		db.Order("category ASC, sort_order ASC").Find(&techStacks)

		var upcomingItems []model.UpcomingItem
		db.Where("is_visible = ?", true).Order("sort_order ASC, created_at ASC").Find(&upcomingItems)

		// Group skills by category.
		skillsByCategory := make(map[string][]model.Skill)
		for _, s := range skills {
			skillsByCategory[s.Category] = append(skillsByCategory[s.Category], s)
		}

		// Group tech stacks by category.
		techByCategory := make(map[string][]model.TechStack)
		for _, t := range techStacks {
			techByCategory[t.Category] = append(techByCategory[t.Category], t)
		}

		// Get visitor session info for comment section.
		visitorLoggedIn := false
		token := c.Cookies("visitor_session")
		if token != "" {
			visitorLoggedIn = true
		}

		cfg := config.MyPortfolio.Get()
		// Build OG image and description for social sharing.
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
			"Title":            owner.FullName,
			"BaseURL":          cfg.App.BaseURL,
			"OGImage":          ogImage,
			"OGDescription":    ogDesc,
			"Owner":            owner,
			"Projects":         projects,
			"HasMoreProjects":  totalProjects > pageSize,
			"NextPage":         2,
			"Experiences":      experiences,
			"Skills":           skills,
			"SkillsByCategory": skillsByCategory,
			"TechByCategory":   techByCategory,
			"SocialLinks":      socialLinks,
			"UpcomingItems":    upcomingItems,
			"VisitorLoggedIn":  visitorLoggedIn,
			"SupportedLangs":   cfg.I18n.SupportedLangs,
			"DefaultLang":      cfg.I18n.DefaultLang,
		}, "layouts/public_base")
	}
}

// ProjectsPage returns the next batch of projects as an HTMX partial.
func ProjectsPage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		const pageSize = 6
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * pageSize

		var projects []model.Project
		var total int64
		db.Model(&model.Project{}).Where("status = ?", "published").Count(&total)
		db.Where("status = ?", "published").
			Order("sort_order ASC, created_at DESC").
			Offset(offset).Limit(pageSize).
			Find(&projects)

		hasMore := int64(offset+pageSize) < total

		return c.Render("partials/project_cards", fiber.Map{
			"Projects":        projects,
			"HasMoreProjects": hasMore,
			"NextPage":        page + 1,
		})
	}
}
