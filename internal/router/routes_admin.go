package router

import (
	"my-portfolio/internal/handler/admin"
	"my-portfolio/internal/middleware"
	"my-portfolio/internal/model"

	contribfgprof "github.com/gofiber/contrib/v3/fgprof"
	contribmonitor "github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// registerAdminRoutes wires all session-protected /admin/* routes.
// Every route in this group requires a valid admin session cookie.
func registerAdminRoutes(app *fiber.App, db *gorm.DB, rdb *redis.Client) {
	adm := app.Group("/admin", middleware.AdminAuth(db, rdb))

	// ── Dashboard ──────────────────────────────────────────────────
	adm.Get("/", admin.Dashboard(db))

	// ── Server monitor & profiler ──────────────────────────────────
	adm.Get("/monitor", contribmonitor.New(contribmonitor.Config{
		Title: "Portfolio · Server Monitor",
	}))
	// Full goroutine-aware profiler — GET /admin/debug/fgprof?seconds=10
	adm.Use(contribfgprof.New(contribfgprof.Config{Prefix: "/admin"}))

	// ── Owner / About ──────────────────────────────────────────────
	adm.Get("/owner", admin.OwnerEditPage(db))
	adm.Put("/owner", admin.OwnerUpdate(db))
	adm.Post("/owner/upload-image", admin.OwnerUploadImage(db))
	adm.Post("/owner/upload-resume", admin.OwnerUploadResume(db))

	// ── Projects ───────────────────────────────────────────────────
	adm.Get("/projects", admin.ProjectListPage())
	adm.Get("/projects/list", admin.ProjectListPartial(db))
	adm.Get("/projects/new", admin.ProjectNewForm())
	adm.Get("/projects/:id/edit", admin.ProjectEditForm(db))
	adm.Post("/projects", admin.ProjectCreate(db))
	adm.Put("/projects/:id", admin.ProjectUpdate(db))
	adm.Delete("/projects/:id", admin.ProjectDelete(db))

	// ── Experience ─────────────────────────────────────────────────
	adm.Get("/experience", admin.ExperienceListPage())
	adm.Get("/experience/list", admin.ExperienceListPartial(db))
	adm.Get("/experience/new", admin.ExperienceNewForm())
	adm.Get("/experience/:id/edit", admin.ExperienceEditForm(db))
	adm.Post("/experience", admin.ExperienceCreate(db))
	adm.Put("/experience/:id", admin.ExperienceUpdate(db))
	adm.Delete("/experience/:id", admin.ExperienceDelete(db))
	adm.Post("/experience/upload-image", admin.ExperienceUploadImage(db))

	// ── Skills ─────────────────────────────────────────────────────
	adm.Get("/skills", admin.SkillListPage())
	adm.Get("/skills/list", admin.SkillListPartial(db))
	adm.Get("/skills/new", admin.SkillNewForm())
	adm.Get("/skills/:id/edit", admin.SkillEditForm(db))
	adm.Post("/skills", admin.SkillCreate(db))
	adm.Put("/skills/:id", admin.SkillUpdate(db))
	adm.Delete("/skills/:id", admin.SkillDelete(db))

	// ── Social Links ────────────────────────────────────────────────
	adm.Get("/social-links", admin.SocialListPage())
	adm.Get("/social-links/list", admin.SocialListPartial(db))
	adm.Get("/social-links/new", admin.SocialNewForm())
	adm.Get("/social-links/:id/edit", admin.SocialEditForm(db))
	adm.Post("/social-links", admin.SocialCreate(db))
	adm.Put("/social-links/:id", admin.SocialUpdate(db))
	adm.Delete("/social-links/:id", admin.SocialDelete(db))

	// ── Uploads ────────────────────────────────────────────────────
	adm.Get("/uploads", admin.UploadListPage())
	adm.Get("/uploads/list", admin.UploadListPartial(db))
	adm.Post("/uploads", admin.UploadCreate(db))
	adm.Delete("/uploads/:id", admin.UploadDelete(db))

	// ── Comment moderation ─────────────────────────────────────────
	adm.Get("/comments", admin.CommentListPage())
	adm.Get("/comments/list", admin.CommentListPartial(db))
	adm.Put("/comments/:id/approve", admin.CommentApprove(db))
	adm.Put("/comments/:id/reject", admin.CommentReject(db))
	adm.Delete("/comments/:id", admin.CommentDelete(db))
	adm.Post("/comments/:id/reply", admin.CommentReply(db))

	// ── Contact messages ────────────────────────────────────────────
	adm.Get("/contacts", admin.ContactListPage())
	adm.Get("/contacts/list", admin.ContactListPartial(db))
	adm.Put("/contacts/:id/read", admin.ContactMarkRead(db))
	adm.Delete("/contacts/:id", admin.ContactDelete(db))

	// ── Tech Stacks ────────────────────────────────────────────────
	adm.Get("/tech-stacks", admin.TechStackListPage())
	adm.Get("/tech-stacks/list", admin.TechStackListPartial(db))
	adm.Get("/tech-stacks/new", admin.TechStackNewForm())
	adm.Get("/tech-stacks/:id/edit", admin.TechStackEditForm(db))
	adm.Post("/tech-stacks", admin.TechStackCreate(db))
	adm.Put("/tech-stacks/:id", admin.TechStackUpdate(db))
	adm.Delete("/tech-stacks/:id", admin.TechStackDelete(db))

	// ── Certificates ────────────────────────────────────────────────
	adm.Get("/certificates", admin.CertificateListPage())
	adm.Get("/certificates/list", admin.CertificateListPartial(db))
	adm.Get("/certificates/new", admin.CertificateForm(&model.Certificate{}, "/admin/certificates/new"))
	adm.Post("/certificates/new", admin.CreateCertificate(db))
	adm.Get("/certificates/edit/:id", func(c fiber.Ctx) error {
		var cert model.Certificate
		if err := db.First(&cert, c.Params("id")).Error; err != nil {
			return c.Status(404).SendString("Certificate not found")
		}
		// call the handler returned by CertificateForm with current context
		return admin.CertificateForm(&cert, "/admin/certificates/edit/"+c.Params("id"))(c)
	})
	adm.Post("/certificates/edit/:id", admin.EditCertificate(db))
	// legacy GET delete kept for compatibility, new HTMX uses DELETE
	adm.Delete("/certificates/:id", admin.DeleteCertificate(db))
	adm.Get("/certificates/delete/:id", admin.DeleteCertificate(db))

	// ── Blog Posts ─────────────────────────────────────────────────
	adm.Get("/posts", admin.PostListPage())
	adm.Get("/posts/list", admin.PostListPartial(db))
	adm.Get("/posts/new", admin.PostNewForm())
	adm.Get("/posts/:id/edit", admin.PostEditForm(db))
	adm.Post("/posts", admin.PostCreate(db))
	adm.Put("/posts/:id", admin.PostUpdate(db))
	adm.Delete("/posts/:id", admin.PostDelete(db))
	adm.Post("/posts/upload-thumbnail", admin.PostUploadThumbnail(db))
	adm.Post("/posts/upload-media", admin.PostUploadMedia(db))
	adm.Post("/posts/upload-video", admin.PostUploadVideo(db))
	adm.Post("/posts/upload-audio", admin.PostUploadAudio(db))

	// ── Upcoming Items ─────────────────────────────────────────────
	adm.Get("/upcoming", admin.UpcomingListPage())
	adm.Get("/upcoming/list", admin.UpcomingListPartial(db))
	adm.Get("/upcoming/new", admin.UpcomingNewForm())
	adm.Get("/upcoming/:id/edit", admin.UpcomingEditForm(db))
	adm.Post("/upcoming", admin.UpcomingCreate(db))
	adm.Put("/upcoming/:id", admin.UpcomingUpdate(db))
	adm.Delete("/upcoming/:id", admin.UpcomingDelete(db))
}
