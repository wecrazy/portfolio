// Package router registers all application routes.
package router

import (
	"my-portfolio/internal/handler"
	"my-portfolio/internal/handler/admin"
	"my-portfolio/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes wires up every route in the application.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// ── Public routes ──────────────────────────────────────────────
	app.Get("/", handler.PortfolioPage(db))
	app.Get("/resume", handler.ServeResumePDF(db))
	app.Get("/resume/download", handler.DownloadResumePDF(db))

	// OAuth
	app.Get("/auth/google", handler.GoogleLogin())
	app.Get("/auth/google/callback", handler.GoogleCallback(db))
	app.Get("/auth/github", handler.GitHubLogin())
	app.Get("/auth/github/callback", handler.GitHubCallback(db))
	app.Get("/auth/logout", handler.OAuthLogout())

	// Comments (public)
	app.Get("/comments", handler.GetComments(db))
	app.Post("/comments", middleware.OAuthAuth(), handler.PostComment(db))

	// Projects pagination (HTMX lazy load)
	app.Get("/projects", handler.ProjectsPage(db))

	// Blog
	app.Get("/blog/more", handler.BlogPostsPartial(db))
	app.Get("/blog", handler.BlogListPage(db))
	app.Get("/blog/:slug", handler.BlogPostPage(db))

	// Contact
	app.Post("/contact", handler.SubmitContact(db))

	// ── Admin auth (no middleware) ─────────────────────────────────
	app.Get("/admin/login", handler.AdminLoginPage())
	app.Post("/admin/login", handler.AdminLoginSubmit(db))
	app.Post("/admin/logout", handler.AdminLogout())

	// ── Admin routes (protected) ───────────────────────────────────
	adm := app.Group("/admin", middleware.AdminAuth(db))

	adm.Get("/", admin.Dashboard(db))

	// Owner / About
	adm.Get("/owner", admin.OwnerEditPage(db))
	adm.Put("/owner", admin.OwnerUpdate(db))
	adm.Post("/owner/upload-image", admin.OwnerUploadImage(db))
	adm.Post("/owner/upload-resume", admin.OwnerUploadResume(db))

	// Projects CRUD
	adm.Get("/projects", admin.ProjectListPage())
	adm.Get("/projects/list", admin.ProjectListPartial(db))
	adm.Get("/projects/new", admin.ProjectNewForm())
	adm.Get("/projects/:id/edit", admin.ProjectEditForm(db))
	adm.Post("/projects", admin.ProjectCreate(db))
	adm.Put("/projects/:id", admin.ProjectUpdate(db))
	adm.Delete("/projects/:id", admin.ProjectDelete(db))

	// Experience CRUD
	adm.Get("/experience", admin.ExperienceListPage())
	adm.Get("/experience/list", admin.ExperienceListPartial(db))
	adm.Get("/experience/new", admin.ExperienceNewForm())
	adm.Get("/experience/:id/edit", admin.ExperienceEditForm(db))
	adm.Post("/experience", admin.ExperienceCreate(db))
	adm.Put("/experience/:id", admin.ExperienceUpdate(db))
	adm.Delete("/experience/:id", admin.ExperienceDelete(db))

	// Skills CRUD
	adm.Get("/skills", admin.SkillListPage())
	adm.Get("/skills/list", admin.SkillListPartial(db))
	adm.Get("/skills/new", admin.SkillNewForm())
	adm.Get("/skills/:id/edit", admin.SkillEditForm(db))
	adm.Post("/skills", admin.SkillCreate(db))
	adm.Put("/skills/:id", admin.SkillUpdate(db))
	adm.Delete("/skills/:id", admin.SkillDelete(db))

	// Social Links CRUD
	adm.Get("/social-links", admin.SocialListPage())
	adm.Get("/social-links/list", admin.SocialListPartial(db))
	adm.Get("/social-links/new", admin.SocialNewForm())
	adm.Get("/social-links/:id/edit", admin.SocialEditForm(db))
	adm.Post("/social-links", admin.SocialCreate(db))
	adm.Put("/social-links/:id", admin.SocialUpdate(db))
	adm.Delete("/social-links/:id", admin.SocialDelete(db))

	// Uploads
	adm.Get("/uploads", admin.UploadListPage())
	adm.Get("/uploads/list", admin.UploadListPartial(db))
	adm.Post("/uploads", admin.UploadCreate(db))
	adm.Delete("/uploads/:id", admin.UploadDelete(db))

	// Comment moderation
	adm.Get("/comments", admin.CommentListPage())
	adm.Get("/comments/list", admin.CommentListPartial(db))
	adm.Put("/comments/:id/approve", admin.CommentApprove(db))
	adm.Put("/comments/:id/reject", admin.CommentReject(db))
	adm.Delete("/comments/:id", admin.CommentDelete(db))
	adm.Post("/comments/:id/reply", admin.CommentReply(db))

	// Contact messages
	adm.Get("/contacts", admin.ContactListPage())
	adm.Get("/contacts/list", admin.ContactListPartial(db))
	adm.Put("/contacts/:id/read", admin.ContactMarkRead(db))
	adm.Delete("/contacts/:id", admin.ContactDelete(db))

	// Tech Stack CRUD
	adm.Get("/tech-stacks", admin.TechStackListPage())
	adm.Get("/tech-stacks/list", admin.TechStackListPartial(db))
	adm.Get("/tech-stacks/new", admin.TechStackNewForm())
	adm.Get("/tech-stacks/:id/edit", admin.TechStackEditForm(db))
	adm.Post("/tech-stacks", admin.TechStackCreate(db))
	adm.Put("/tech-stacks/:id", admin.TechStackUpdate(db))
	adm.Delete("/tech-stacks/:id", admin.TechStackDelete(db))

	// Blog Posts CRUD
	adm.Get("/posts", admin.PostListPage())
	adm.Get("/posts/list", admin.PostListPartial(db))
	adm.Get("/posts/new", admin.PostNewForm())
	adm.Get("/posts/:id/edit", admin.PostEditForm(db))
	adm.Post("/posts", admin.PostCreate(db))
	adm.Put("/posts/:id", admin.PostUpdate(db))
	adm.Delete("/posts/:id", admin.PostDelete(db))
	adm.Post("/posts/upload-thumbnail", admin.PostUploadThumbnail(db))

	// Upcoming Items CRUD
	adm.Get("/upcoming", admin.UpcomingListPage())
	adm.Get("/upcoming/list", admin.UpcomingListPartial(db))
	adm.Get("/upcoming/new", admin.UpcomingNewForm())
	adm.Get("/upcoming/:id/edit", admin.UpcomingEditForm(db))
	adm.Post("/upcoming", admin.UpcomingCreate(db))
	adm.Put("/upcoming/:id", admin.UpcomingUpdate(db))
	adm.Delete("/upcoming/:id", admin.UpcomingDelete(db))
}
