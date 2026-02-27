package admin

import (
	"fmt"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"
	"my-portfolio/pkg/pagination"

	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// PostListPage renders the posts admin page.
func PostListPage() fiber.Handler {
	return func(c fiber.Ctx) error {
		cfg := config.MyPortfolio.Get()
		return c.Render("admin/posts", fiber.Map{
			"Title":          "Blog Posts",
			"Admin":          c.Locals("admin"),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/admin_base")
	}
}

// PostListPartial returns the posts table rows as an HTMX partial.
func PostListPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		params := pagination.ParseParams(c, "sort_order", []string{"sort_order", "title", "status", "created_at"})
		var posts []model.Post
		query, pageResult := pagination.Paginate(db, &model.Post{}, params, []string{"title", "excerpt", "tags"})
		query.Find(&posts)
		return c.Render("partials/post_rows", fiber.Map{
			"Posts":      posts,
			"Pagination": pageResult,
		})
	}
}

// PostNewForm renders an empty post form partial.
func PostNewForm() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Render("partials/post_form", fiber.Map{"Post": model.Post{}})
	}
}

// PostEditForm renders a pre-filled post form partial.
func PostEditForm(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var post model.Post
		if err := db.Preload("ThumbnailFile").First(&post, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		return c.Render("partials/post_form", fiber.Map{"Post": post})
	}
}

// PostCreate handles creating a new post.
func PostCreate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var post model.Post
		if err := c.Bind().Body(&post); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		post.Slug = slug.Make(post.Title)
		if post.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
		if err := db.Create(&post).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create post")
		}
		setToast(c, "post_created", "success")
		return c.Render("partials/post_row", fiber.Map{"Post": post})
	}
}

// PostUpdate handles updating an existing post.
func PostUpdate(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var post model.Post
		if err := db.First(&post, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		prevStatus := post.Status
		if err := c.Bind().Body(&post); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
		}
		post.Slug = slug.Make(post.Title)
		if post.Status == "published" && prevStatus != "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
		db.Save(&post)
		setToast(c, "post_updated", "success")
		return c.Render("partials/post_row", fiber.Map{"Post": post})
	}
}

// PostDelete handles deleting a post.
func PostDelete(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := db.Delete(&model.Post{}, c.Params("id")).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete")
		}
		setToast(c, "post_deleted", "success")
		return c.SendString("")
	}
}

// PostUploadThumbnail handles thumbnail image upload for a post.
func PostUploadThumbnail(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("thumbnail")
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

		// Return preview image + OOB swap to update the hidden thumbnail_file_id input.
		html := fmt.Sprintf(
			`<img src="/uploads/images/%s" class="img-thumbnail mt-2" style="max-height:120px" alt="Thumbnail">
<input type="hidden" id="thumbnail_file_id" name="thumbnail_file_id" value="%d" hx-swap-oob="outerHTML:#thumbnail_file_id">`,
			uploaded.StoredName, uploaded.ID,
		)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	}
}
