package admin

import (
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ProjectListPage renders the projects admin page.
func ProjectListPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin/projects", fiber.Map{
			"Title": "Projects",
			"Admin": c.Locals("admin"),
		}, "layouts/admin_base")
	}
}

// ProjectListPartial returns the projects table rows as an HTMX partial.
func ProjectListPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projects []model.Project
		db.Order("sort_order ASC, created_at DESC").Find(&projects)
		return c.Render("partials/project_rows", fiber.Map{"Projects": projects})
	}
}

// ProjectNewForm renders an empty project form partial.
func ProjectNewForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("partials/project_form", fiber.Map{"Project": model.Project{}})
	}
}

// ProjectEditForm renders a pre-filled project form partial.
func ProjectEditForm(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var project model.Project
		if err := db.First(&project, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Project not found")
		}
		return c.Render("partials/project_form", fiber.Map{"Project": project})
	}
}

// ProjectCreate handles creating a new project.
func ProjectCreate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var project model.Project
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		project.Slug = slug.Make(project.Title)
		project.Featured = c.FormValue("featured") == "on"
		if err := db.Create(&project).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create project")
		}
		c.Set("HX-Trigger", `{"showToast":"Project created"}`)
		return c.Render("partials/project_row", fiber.Map{"Project": project})
	}
}

// ProjectUpdate handles updating an existing project.
func ProjectUpdate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var project model.Project
		if err := db.First(&project, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Project not found")
		}
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		project.Slug = slug.Make(project.Title)
		project.Featured = c.FormValue("featured") == "on"
		db.Save(&project)
		c.Set("HX-Trigger", `{"showToast":"Project updated"}`)
		return c.Render("partials/project_row", fiber.Map{"Project": project})
	}
}

// ProjectDelete handles deleting a project.
func ProjectDelete(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Delete(&model.Project{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		c.Set("HX-Trigger", `{"showToast":"Project deleted"}`)
		return c.SendString("")
	}
}
