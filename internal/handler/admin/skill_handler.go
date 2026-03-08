package admin

import (
	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// SkillListPage renders the skills admin page.
func SkillListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		title, _ := appI18n.T.Localize(c, "admin.skills.title")
		return c.Render("admin/skills", fiber.Map{
			"Title":          title,
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// SkillListPartial returns the skills table rows as an HTMX partial.
func SkillListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "category", []string{"category", "sort_order", "name", "proficiency"})
		var items []model.Skill
		query, pageResult := pagination.Paginate(db, &model.Skill{}, params, []string{"name", "category"})
		query.Find(&items)
		return c.Render("partials/skill_rows", fiber.Map{
			"Skills":     items,
			"Pagination": pageResult,
		})
	}
}

// SkillNewForm renders an empty skill form partial.
func SkillNewForm() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Render("partials/skill_form", fiber.Map{"Skill": model.Skill{}})
	}
}

// SkillEditForm renders a pre-filled skill form partial.
func SkillEditForm(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Skill
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		return c.Render("partials/skill_form", fiber.Map{"Skill": item})
	}
}

// SkillCreate handles creating a new skill.
func SkillCreate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Skill
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		if err := db.Create(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create")
		}
		setToast(c, "skill_created", "success")
		return c.Render("partials/skill_row", fiber.Map{"Skill": item})
	}
}

// SkillUpdate handles updating an existing skill.
func SkillUpdate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var item model.Skill
		if err := db.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		if err := c.Bind().Body(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		db.Save(&item)
		setToast(c, "skill_updated", "success")
		return c.Render("partials/skill_row", fiber.Map{"Skill": item})
	}
}

// SkillDelete handles deleting a skill.
func SkillDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := db.Delete(&model.Skill{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		setToast(c, "skill_deleted", "success")
		return c.SendString("")
	}
}
