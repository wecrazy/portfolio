package handler

import (
	"strconv"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// blogPageData loads shared data needed by blog pages (owner, social links for navbar).
func blogPageData(db *gorm.DB) (model.Owner, []model.SocialLink) {
	var owner model.Owner
	db.Preload("ProfileImage").Preload("ResumeFile").First(&owner)
	var socialLinks []model.SocialLink
	db.Order("sort_order ASC").Find(&socialLinks)
	return owner, socialLinks
}

// BlogListPage renders the public blog listing page.
func BlogListPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const pageSize = 6
		owner, socialLinks := blogPageData(db)

		var posts []model.Post
		var total int64
		db.Model(&model.Post{}).Where("status = ?", "published").Count(&total)
		db.Where("status = ?", "published").
			Preload("ThumbnailFile").
			Order("sort_order ASC, published_at DESC").
			Limit(pageSize).
			Find(&posts)

		hasMore := total > pageSize

		return c.Render("public/blog", fiber.Map{
			"Title":          "Blog",
			"Owner":          owner,
			"SocialLinks":    socialLinks,
			"Posts":          posts,
			"HasMorePosts":   hasMore,
			"NextPage":       2,
			"SupportedLangs": config.MyPortfolio.Get().I18n.SupportedLangs,
			"DefaultLang":    config.MyPortfolio.Get().I18n.DefaultLang,
		}, "layouts/public_base")
	}
}

// BlogPostsPartial returns the next batch of posts as an HTMX partial (load more).
func BlogPostsPartial(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const pageSize = 6
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * pageSize

		var posts []model.Post
		var total int64
		db.Model(&model.Post{}).Where("status = ?", "published").Count(&total)
		db.Where("status = ?", "published").
			Preload("ThumbnailFile").
			Order("sort_order ASC, published_at DESC").
			Offset(offset).Limit(pageSize).
			Find(&posts)

		hasMore := int64(offset+pageSize) < total

		return c.Render("partials/blog_cards", fiber.Map{
			"Posts":        posts,
			"HasMorePosts": hasMore,
			"NextPage":     page + 1,
		})
	}
}

// BlogPostPage renders a single blog post detail page.
func BlogPostPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner, socialLinks := blogPageData(db)

		var post model.Post
		if err := db.Preload("ThumbnailFile").
			Where("slug = ? AND status = ?", c.Params("slug"), "published").
			First(&post).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Post not found")
		}

		return c.Render("public/blog_post", fiber.Map{
			"Title":          post.Title,
			"Owner":          owner,
			"SocialLinks":    socialLinks,
			"Post":           post,
			"SupportedLangs": config.MyPortfolio.Get().I18n.SupportedLangs,
			"DefaultLang":    config.MyPortfolio.Get().I18n.DefaultLang,
		}, "layouts/public_base")
	}
}
