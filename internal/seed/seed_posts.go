package seed

import (
	"log"
	"os"
	"time"

	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedDemoMediaFiles registers the pre-generated demo media assets (image, video,
// audio) into the uploaded_files table so the blog posts can reference them by URL.
// It is idempotent: if the records already exist it does nothing.
func seedDemoMediaFiles(db *gorm.DB) {
	type fileSpec struct {
		stored   string
		filePath string
		mime     string
		category string
	}
	specs := []fileSpec{
		{"demo_blog_cover.jpg", "uploads/images/demo_blog_cover.jpg", "image/jpeg", "images"},
		{"demo_video.mp4", "uploads/video/demo_video.mp4", "video/mp4", "video"},
		{"demo_audio.wav", "uploads/audio/demo_audio.wav", "audio/wav", "audio"},
	}
	for _, s := range specs {
		var c int64
		db.Model(&model.UploadedFile{}).Where("stored_name = ?", s.stored).Count(&c)
		if c > 0 {
			continue
		}
		var size int64
		if info, err := os.Stat(s.filePath); err == nil {
			size = info.Size()
		}
		rec := model.UploadedFile{
			OriginalName: s.stored,
			StoredName:   s.stored,
			FilePath:     s.filePath,
			MimeType:     s.mime,
			FileSize:     size,
			Category:     s.category,
		}
		if err := db.Create(&rec).Error; err != nil {
			log.Printf("Warning: could not seed media file record %s: %v", s.stored, err)
		}
	}
}

// seedPosts creates demo blog posts. It is intentionally kept in the seed package because it depends on internal/model.
func seedPosts(db *gorm.DB) error {
	pub1 := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	pub2 := time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC)
	pub3 := time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)
	pub4 := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	posts := []model.Post{
		{
			Title:   "Boosting Productivity While Coding",
			Slug:    "boosting-productivity-while-coding",
			Excerpt: "Practical strategies to stay focused, ship faster, and keep momentum while building software.",
			Content: "## Stay Focused with Small Goals\n\n" +
				"Break work into 25–50 minute focused sessions (Pomodoro), with short breaks to avoid burnout. Use `make` tasks to automate common flows (build, test, run).\n\n" +
				"## Reduce Context Switching\n\n" +
				"Keep your editor and terminal clean. Use short-lived branches and small commits. Prefer server-rendered HTML (HTMX) for UI work when it avoids complex client-side state.\n\n" +
				"## Automate Repetitive Work\n\n" +
				"Use `air` for hot reload during local Go development, `docker-compose` for databases, and scripts for migrations. Automation saves minutes that add up to hours.",
			Tags:        "Productivity,Workflow,Go,Tools",
			Status:      "published",
			SortOrder:   1,
			PublishedAt: &pub1,
		},
		{
			Title:   "Why Lubuntu Can Be Faster and Simpler Than Ubuntu",
			Slug:    "lubuntu-faster-simpler-than-ubuntu",
			Excerpt: "A pragmatic look at lightweight Linux desktops — why Lubuntu may be a better fit for older hardware or minimal setups.",
			Content: "Lubuntu focuses on being lightweight and fast. For many developers and low-resource machines it's a great alternative to Ubuntu.\n\n" +
				"- **Low RAM footprint:** Lubuntu's desktop and default apps use fewer resources.\n" +
				"- **Simplicity:** Fewer background services means less maintenance.\n\n" +
				"If you want to learn more, see the official site: https://lubuntu.me and compare with Ubuntu at https://ubuntu.com.\n\n" +
				"**Tip:** Try Lubuntu in a VM or live USB to measure improvements on your hardware before switching permanently.",
			Tags:        "Linux,Lubuntu,Ubuntu,Sysadmin",
			Status:      "published",
			SortOrder:   2,
			PublishedAt: &pub2,
		},
		{
			Title:   "Structuring Small Go Projects for Speed",
			Slug:    "structuring-small-go-projects",
			Excerpt: "A concise layout and conventions that keep small Go services productive and easy to maintain.",
			Content: "## Keep It Simple\n\n" +
				"Start with a clear `cmd/` entry, `internal/` for private packages, and small, focused packages. Avoid over-architecting — prefer composition over deep abstractions.\n\n" +
				"## Repeatable Local Dev\n\n" +
				"Provide `Makefile` targets for `make run`, `make test`, and `make db-reset`. Use environment-specific config files to keep local development fast.\n\n" +
				"## Fast Iteration\n\n" +
				"Hot-reload, concise templates, and small endpoints let you iterate quickly without losing clarity.",
			Tags:        "Go,Architecture,Productivity",
			Status:      "published",
			SortOrder:   3,
			PublishedAt: &pub3,
		},
		{
			Title:   "Rich Media in Blog Posts — Images, Video & Audio",
			Slug:    "rich-media-in-blog-posts",
			Excerpt: "This blog now supports embedded images, videos, and audio clips inside Markdown posts. Here's a quick tour of each media type.",
			Content: "## Embedded Images\n\n" +
				"Drop any uploaded image directly into your post using standard Markdown syntax:\n\n" +
				"```md\n![Alt text](/uploads/images/your-file.jpg)\n```\n\n" +
				"![Demo cover image](/uploads/images/demo_blog_cover.jpg)\n\n" +
				"---\n\n" +
				"## Embedded Video\n\n" +
				"Use a raw HTML `<video>` block for self-hosted videos uploaded through the admin panel:\n\n" +
				"<video controls style=\"width:100%;border-radius:0.75rem;margin:1rem 0\">\n" +
				"  <source src=\"/uploads/video/demo_video.mp4\" type=\"video/mp4\">\n" +
				"  Your browser does not support the video tag.\n" +
				"</video>\n\n" +
				"You can also embed YouTube or Vimeo via an `<iframe>`:\n\n" +
				"```html\n<iframe src=\"https://www.youtube.com/embed/dQw4w9WgXcQ\"\n  allowfullscreen style=\"width:100%;aspect-ratio:16/9\"></iframe>\n```\n\n" +
				"---\n\n" +
				"## Embedded Audio\n\n" +
				"Self-hosted audio files work just as smoothly:\n\n" +
				"<audio controls style=\"width:100%;margin:1rem 0\">\n" +
				"  <source src=\"/uploads/audio/demo_audio.wav\" type=\"audio/wav\">\n" +
				"  Your browser does not support audio playback.\n" +
				"</audio>\n\n" +
				"---\n\n" +
				"## How It Works\n\n" +
				"Posts are written in **Markdown with raw HTML support**, rendered by [goldmark](https://github.com/yuin/goldmark) on the server and sanitized by [bluemonday](https://github.com/microcosm-cc/bluemonday) before display.\n\n" +
				"Use the three upload buttons in the admin editor — **Insert Image**, **Insert Video**, **Insert Audio** — to upload files and insert the correct snippet automatically.",
			Tags:        "Blog,Media,Markdown,Tutorial",
			Status:      "published",
			SortOrder:   4,
			PublishedAt: &pub4,
		},
	}
	return db.Create(&posts).Error
}
