package handler

import (
	"bytes"
	"html/template"
	"regexp"
	"strconv"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"gorm.io/gorm"
)

// markdownRenderer parses GFM (tables, strikethrough, etc.) and allows raw HTML blocks.
var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
)

// htmlPolicy builds a bluemonday policy that allows rich blog content including
// images, audio, video, and YouTube/Vimeo iframes written in raw Markdown HTML blocks.
var htmlPolicy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	// Media elements
	p.AllowElements("video", "audio", "source", "figure", "figcaption")
	p.AllowAttrs("controls", "autoplay", "loop", "muted", "preload", "width", "height", "style", "class").OnElements("video", "audio")
	p.AllowAttrs("src", "type").OnElements("source")
	// Iframes only from YouTube / Vimeo
	p.AllowElements("iframe")
	p.AllowAttrs("width", "height", "frameborder", "allow", "allowfullscreen", "title", "style", "class").OnElements("iframe")
	p.AllowAttrs("src").Matching(
		regexp.MustCompile(`^https?://(www\.)?(youtube\.com|youtu\.be|player\.vimeo\.com)/`),
	).OnElements("iframe")
	// Allow style on any element (for alignment etc.)
	p.AllowAttrs("style", "class").Globally()
	return p
}()

// renderMarkdown converts Markdown (with raw HTML blocks) to sanitized HTML.
func renderMarkdown(content string) template.HTML {
	var buf bytes.Buffer
	if err := markdownRenderer.Convert([]byte(content), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(content))
	}
	return template.HTML(htmlPolicy.SanitizeBytes(buf.Bytes()))
}

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
	return func(c fiber.Ctx) error {
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

		cfg := config.MyPortfolio.Get()
		ogImage := cfg.App.BaseURL + "/static/img/favicon.svg"
		if owner.ProfileImage != nil {
			ogImage = cfg.App.BaseURL + "/uploads/images/" + owner.ProfileImage.StoredName
		}

		return c.Render("public/blog", fiber.Map{
			"Title":          "Blog — " + owner.FullName,
			"BaseURL":        cfg.App.BaseURL,
			"OGImage":        ogImage,
			"OGDescription":  "Tech blog covering software development, projects and ideas by " + owner.FullName + ".",
			"Owner":          owner,
			"SocialLinks":    socialLinks,
			"Posts":          posts,
			"HasMorePosts":   hasMore,
			"NextPage":       2,
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/public_base")
	}
}

// BlogPostsPartial returns the next batch of posts as an HTMX partial (load more).
func BlogPostsPartial(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
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
	return func(c fiber.Ctx) error {
		owner, socialLinks := blogPageData(db)

		var post model.Post
		if err := db.Preload("ThumbnailFile").
			Where("slug = ? AND status = ?", c.Params("slug"), "published").
			First(&post).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Post not found")
		}

		cfg := config.MyPortfolio.Get()
		ogImage := cfg.App.BaseURL + "/static/img/favicon.svg"
		if post.ThumbnailFile != nil {
			ogImage = cfg.App.BaseURL + "/uploads/images/" + post.ThumbnailFile.StoredName
		} else if owner.ProfileImage != nil {
			ogImage = cfg.App.BaseURL + "/uploads/images/" + owner.ProfileImage.StoredName
		}
		ogDesc := post.Excerpt
		if len(ogDesc) > 160 {
			ogDesc = ogDesc[:157] + "..."
		}
		if ogDesc == "" {
			ogDesc = post.Title
		}

		return c.Render("public/blog_post", fiber.Map{
			"Title":          post.Title,
			"BaseURL":        cfg.App.BaseURL,
			"OGImage":        ogImage,
			"OGDescription":  ogDesc,
			"Owner":          owner,
			"SocialLinks":    socialLinks,
			"Post":           post,
			"PostHTML":       renderMarkdown(post.Content),
			"SupportedLangs": cfg.I18n.SupportedLangs,
			"DefaultLang":    cfg.I18n.DefaultLang,
		}, "layouts/public_base")
	}
}
