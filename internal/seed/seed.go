// Package seed provides initial data seeding for the database.
package seed

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/fileutil"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// deviconCDN is the base URL for fetching technology icons in seeds. It is used by both skills and tech stacks to ensure consistent icons.  It is intentionally kept in the seed package because it is only relevant for seeding demo data and not used elsewhere in the app.
const deviconCDN = "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons"

// SeedIfNeeded creates the default admin user and an empty owner profile when
// they don't already exist. Safe to call on every startup.
func SeedIfNeeded(db *gorm.DB, cfg config.TypeMyPortfolio) {
	seedAdmin(db, cfg)
	seedOwner(db, cfg)
	seedDemoData(db)
}

// seedAdmin creates the default admin user if no admins exist. Safe to call on every startup.
func seedAdmin(db *gorm.DB, cfg config.TypeMyPortfolio) {
	var count int64
	db.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.DefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default admin password: %v", err)
	}

	admin := model.Admin{
		Username:     cfg.Admin.DefaultUsername,
		Email:        cfg.Admin.DefaultEmail,
		PasswordHash: string(hash),
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}
	log.Printf("Seeded default admin user: %s", admin.Username)
}

// linkImageRecord creates an UploadedFile DB record for an image already on disk.
// It is intentionally kept in the seed package because it depends on internal/model.
func linkImageRecord(db *gorm.DB, storedName, filePath string) *model.UploadedFile {
	ext := strings.ToLower(filepath.Ext(storedName))
	mimeType := fileutil.MimeByExt(ext)
	if mimeType == "application/octet-stream" {
		mimeType = "image/jpeg"
	}
	var size int64
	if info, err := os.Stat(filePath); err == nil {
		size = info.Size()
	}
	rec := &model.UploadedFile{
		OriginalName: storedName,
		StoredName:   storedName,
		FilePath:     filePath,
		MimeType:     mimeType,
		FileSize:     size,
		Category:     "images",
	}
	if err := db.Create(rec).Error; err != nil {
		log.Printf("Warning: failed to create image DB record: %v", err)
		return nil
	}
	return rec
}

// relinkUploadImage scans uploadDir for the newest allowed image and creates a
// DB record for it. Returns nil when no suitable file is found.
func relinkUploadImage(db *gorm.DB, uploadDir string, allowedExts map[string]bool) *model.UploadedFile {
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return nil
	}

	type candidate struct {
		name    string
		modTime time.Time
	}
	var candidates []candidate
	for _, e := range entries {
		if e.IsDir() || !allowedExts[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		if info, err := e.Info(); err == nil {
			candidates = append(candidates, candidate{name: e.Name(), modTime: info.ModTime()})
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].modTime.After(candidates[j].modTime)
	})
	chosen := candidates[0]
	rec := linkImageRecord(db, chosen.name, filepath.Join(uploadDir, chosen.name))
	if rec != nil {
		log.Printf("Re-linked existing profile image from uploads: %s", chosen.name)
	}
	return rec
}

// copyStaticImage copies the configured static profile image into uploadDir and
// creates a DB record for it. Returns nil when the source is absent or invalid.
func copyStaticImage(db *gorm.DB, cfg config.TypeMyPortfolio, uploadDir string, allowedExts map[string]bool) *model.UploadedFile {
	if cfg.Owner.ProfileImage == "" {
		return nil
	}
	srcPath := filepath.Join(cfg.App.StaticDir, strings.TrimPrefix(cfg.Owner.ProfileImage, "/"))
	if !fileutil.Exists(srcPath) {
		return nil
	}
	ext := strings.ToLower(filepath.Ext(srcPath))
	if !allowedExts[ext] {
		return nil
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil
	}
	storedName := uuid.New().String() + ext
	dstPath := filepath.Join(uploadDir, storedName)
	if err := fileutil.CopyFile(srcPath, dstPath); err != nil {
		log.Printf("Warning: could not copy static profile image: %v", err)
		return nil
	}
	rec := linkImageRecord(db, storedName, dstPath)
	if rec != nil {
		log.Printf("Copied static profile image to uploads: %s", storedName)
	}
	return rec
}

// seedOwner creates the default owner profile if it doesn't exist. Safe to call on every startup.
func seedOwner(db *gorm.DB, cfg config.TypeMyPortfolio) {
	var count int64
	db.Model(&model.Owner{}).Count(&count)
	if count > 0 {
		return
	}

	allowedExts := fileutil.AllowedExts(cfg.Upload.AllowedImageTypes)
	uploadDir := filepath.Join(cfg.App.UploadDir, "images")

	// Priority 1: re-link the newest image already in uploads/images/ (survives db-reset).
	imgProfile := relinkUploadImage(db, uploadDir, allowedExts)

	// Priority 2: fall back to the static file declared in config, copy it into uploads/images/.
	if imgProfile == nil {
		imgProfile = copyStaticImage(db, cfg, uploadDir, allowedExts)
	}

	// Resume: pick up any PDF already sitting in uploads/resume/ (survives db-reset).
	var resumeFile *model.UploadedFile
	resumeDir := filepath.Join(cfg.App.UploadDir, "resume")
	if entries, err := os.ReadDir(resumeDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.ToLower(filepath.Ext(name)) != ".pdf" {
				continue
			}
			fullPath := filepath.Join(resumeDir, name)
			var size int64
			if info, err2 := e.Info(); err2 == nil {
				size = info.Size()
			}
			rec := &model.UploadedFile{
				OriginalName: name,
				StoredName:   name,
				FilePath:     fullPath,
				MimeType:     "application/pdf",
				FileSize:     size,
				Category:     "resume",
			}
			if err2 := db.Create(rec).Error; err2 == nil {
				resumeFile = rec
			}
			break // use first PDF found
		}
	}

	owner := model.Owner{
		FullName:     cfg.Owner.Name,
		Title:        cfg.Owner.Title,
		Bio:          cfg.Owner.Bio,
		ProfileImage: imgProfile,
		ResumeFile:   resumeFile,
		Email:        cfg.Owner.Email,
		Phone:        cfg.Owner.Phone,
		Location:     cfg.Owner.Location,
	}
	if err := db.Create(&owner).Error; err != nil {
		log.Fatalf("Failed to seed owner profile: %v", err)
	}
	log.Println("Seeded default owner profile")
}

// seedDemoData creates demo projects, experience, skills, etc if they don't already exist. Safe to call on every startup.
func seedDemoData(db *gorm.DB) {
	// Only seed projects/experience/skills/etc if projects table is empty.
	var count int64
	db.Model(&model.Project{}).Count(&count)
	if count == 0 {
		seedProjects(db)
		seedExperiences(db)
		seedSkills(db)
		seedSocialLinks(db)
		seedTechStacks(db)
		seedComments(db)
		log.Println("Seeded demo data")
	}

	// Seed posts independently so they work even on existing installs.
	var postCount int64
	if err := db.Model(&model.Post{}).Count(&postCount).Error; err != nil {
		log.Printf("ERROR counting posts: %v", err)
	}
	log.Printf("Post count before seeding: %d", postCount)
	if postCount == 0 {
		if err := seedPosts(db); err != nil {
			log.Printf("ERROR seeding posts: %v", err)
		} else {
			log.Println("Seeded demo blog posts")
		}
	}

	// Seed upcoming items independently.
	var upcomingCount int64
	db.Model(&model.UpcomingItem{}).Count(&upcomingCount)
	if upcomingCount == 0 {
		seedUpcomingItems(db)
		log.Println("Seeded demo upcoming items")
	}
}

// seedExperiences creates demo work and education experiences. It is intentionally kept in the seed package because it depends on internal/model.
func seedExperiences(db *gorm.DB) {
	workStart := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	workEnd := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	eduStart := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)

	experiences := []model.Experience{
		{
			Type:        "Work",
			Title:       "Full Stack Programmer",
			Org:         "PT. Cyber Smart Network Asia",
			Location:    "Indonesia",
			StartDate:   workStart,
			EndDate:     &workEnd,
			IsCurrent:   false,
			Description: "Worked as a Full Stack Programmer responsible for designing, developing, and maintaining web-based applications. Built and integrated backend APIs alongside dynamic frontend interfaces to deliver full-cycle software solutions.",
			SortOrder:   1,
			ImageURL:    "https://media.licdn.com/dms/image/v2/C4E0BAQHs0YUQvojhmA/company-logo_200_200/company-logo_200_200/0/1631316640234?e=2147483647&v=beta&t=RSKIYgOB-a3FNYeu3zNumK6iu5Laijgr410euHTSuWA",
		},
		{
			Type:        "Education",
			Title:       "Bachelor of Informatics Engineering",
			Org:         "Universitas Kristen Indonesia Toraja",
			Location:    "Toraja, South Sulawesi, Indonesia",
			StartDate:   eduStart,
			EndDate:     &eduEnd,
			IsCurrent:   false,
			Description: "Graduated Cum Laude from the Faculty of Engineering, Department of Informatics Engineering with a GPA of 3.89, earning recognition as the best graduate of the class.",
			SortOrder:   2,
			ImageURL:    "https://ukitoraja.ac.id/wp-content/uploads/2019/05/Logo-UKIT.png",
		},
	}
	db.Create(&experiences)
}

// seedProjects creates demo projects. It is intentionally kept in the seed package because it depends on internal/model.
func seedProjects(db *gorm.DB) {
	projects := []model.Project{
		{Title: "E-Commerce Platform", Slug: "e-commerce-platform", Description: "A full-featured e-commerce platform with cart, checkout, and payment integration.", Tags: "Go,React,PostgreSQL,Stripe", Status: "published", SortOrder: 1, Featured: true, LiveURL: "https://example.com", RepoURL: "https://github.com/example/ecommerce", ThumbnailURL: "https://placehold.co/600x400/6366f1/ffffff?text=E-Commerce"},
		{Title: "Task Management App", Slug: "task-management-app", Description: "Real-time task management application with team collaboration features.", Tags: "TypeScript,Next.js,Prisma,WebSocket", Status: "published", SortOrder: 2, Featured: true, RepoURL: "https://github.com/example/tasks", ThumbnailURL: "https://placehold.co/600x400/8b5cf6/ffffff?text=Task+Manager"},
		{Title: "Weather Dashboard", Slug: "weather-dashboard", Description: "Beautiful weather dashboard with 7-day forecast, radar maps, and location search.", Tags: "Vue.js,OpenWeather API,Chart.js", Status: "published", SortOrder: 3, LiveURL: "https://example.com/weather", ThumbnailURL: "https://placehold.co/600x400/0ea5e9/ffffff?text=Weather+Dashboard"},
		{Title: "Blog Engine", Slug: "blog-engine", Description: "Markdown-powered blog engine with SEO optimization and RSS feed.", Tags: "Go,Fiber,SQLite,HTMX", Status: "published", SortOrder: 4, RepoURL: "https://github.com/example/blog", ThumbnailURL: "https://placehold.co/600x400/10b981/ffffff?text=Blog+Engine"},
		{Title: "Chat Application", Slug: "chat-application", Description: "Real-time chat app with rooms, direct messages, and file sharing.", Tags: "Go,WebSocket,Redis,React", Status: "published", SortOrder: 5, Featured: true, ThumbnailURL: "https://placehold.co/600x400/f59e0b/ffffff?text=Chat+App"},
		{Title: "Portfolio Builder", Slug: "portfolio-builder", Description: "Drag-and-drop portfolio builder for developers with custom themes.", Tags: "Next.js,Tailwind,MongoDB", Status: "published", SortOrder: 6, ThumbnailURL: "https://placehold.co/600x400/ec4899/ffffff?text=Portfolio+Builder"},
		{Title: "API Gateway", Slug: "api-gateway", Description: "High-performance API gateway with rate limiting, caching, and auth.", Tags: "Go,Redis,Docker,gRPC", Status: "published", SortOrder: 7, RepoURL: "https://github.com/example/gateway", ThumbnailURL: "https://placehold.co/600x400/ef4444/ffffff?text=API+Gateway"},
		{Title: "Mobile Fitness App", Slug: "mobile-fitness-app", Description: "Cross-platform fitness tracking app with workout plans and progress charts.", Tags: "Flutter,Firebase,Dart", Status: "published", SortOrder: 8, ThumbnailURL: "https://placehold.co/600x400/14b8a6/ffffff?text=Fitness+App"},
	}
	db.Create(&projects)
}

// seedSkills creates demo skills. It is intentionally kept in the seed package because it depends on internal/model.
func seedSkills(db *gorm.DB) {
	skills := []model.Skill{
		{Name: "Go", Category: "Languages", IconClass: "devicon-go-original-wordmark", IconURL: deviconCDN + "/go/go-original-wordmark.svg", Proficiency: 85, SortOrder: 1},
		{Name: "Python", Category: "Languages", IconClass: "devicon-python-plain", IconURL: deviconCDN + "/python/python-original.svg", Proficiency: 75, SortOrder: 2},
		{Name: "JavaScript", Category: "Languages", IconClass: "devicon-javascript-plain", IconURL: deviconCDN + "/javascript/javascript-plain.svg", Proficiency: 76, SortOrder: 3},
		{Name: "PHP", Category: "Languages", IconClass: "devicon-php-plain", IconURL: deviconCDN + "/php/php-plain.svg", Proficiency: 80, SortOrder: 4},
		{Name: "Rust", Category: "Languages", IconClass: "devicon-rust-original", IconURL: deviconCDN + "/rust/rust-original.svg", Proficiency: 5, SortOrder: 5},
		{Name: "C++", Category: "Languages", IconClass: "devicon-cplusplus-plain", IconURL: deviconCDN + "/cplusplus/cplusplus-original.svg", Proficiency: 50, SortOrder: 6},
		{Name: "Java", Category: "Languages", IconClass: "devicon-java-plain", IconURL: deviconCDN + "/java/java-original.svg", Proficiency: 30, SortOrder: 7},
		{Name: "HTMX", Category: "Frontend", IconClass: "bxf bx-bolt-circle", IconURL: "https://cdn.jsdelivr.net/gh/bigskysoftware/htmx@v2.0.4/www/static/img/htmx_logo.1.png", Proficiency: 90, SortOrder: 1},
		{Name: "Fiber", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg", Proficiency: 90, SortOrder: 1},
		{Name: "Gin", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png", Proficiency: 80, SortOrder: 2},
		{Name: "Code Igniter", Category: "Backend", IconClass: "devicon-codeigniter-plain", IconURL: deviconCDN + "/codeigniter/codeigniter-plain.svg", Proficiency: 80, SortOrder: 3},
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", IconURL: deviconCDN + "/docker/docker-original.svg", Proficiency: 76, SortOrder: 1},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", IconURL: deviconCDN + "/linux/linux-original.svg", Proficiency: 80, SortOrder: 2},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", IconURL: deviconCDN + "/git/git-original.svg", Proficiency: 85, SortOrder: 3},
		{Name: "Adobe Photoshop", Category: "Design", IconClass: "devicon-photoshop-plain", IconURL: deviconCDN + "/photoshop/photoshop-plain.svg", Proficiency: 50, SortOrder: 1},
		{Name: "Adobe After Effects", Category: "Design", IconClass: "devicon-aftereffects-plain", IconURL: deviconCDN + "/aftereffects/aftereffects-plain.svg", Proficiency: 30, SortOrder: 2},
		{Name: "Adobe Premiere Pro", Category: "Design", IconClass: "devicon-premierepro-plain", IconURL: deviconCDN + "/premierepro/premierepro-plain.svg", Proficiency: 25, SortOrder: 3},
	}
	db.Create(&skills)
}

// seedSocialLinks creates demo social links. It is intentionally kept in the seed package because it depends on internal/model.
func seedSocialLinks(db *gorm.DB) {
	links := []model.SocialLink{
		{Platform: "GitHub", URL: "https://github.com/wecrazy", IconClass: "bxl bx-github", Label: "GitHub", SortOrder: 1},
		{Platform: "LinkedIn", URL: "https://id.linkedin.com/in/wegirandol-histara-littu-926219195", IconClass: "bxl bx-linkedin", Label: "LinkedIn", SortOrder: 2},
		{Platform: "Instagram", URL: "https://www.instagram.com/wecraz_y", IconClass: "bxl bx-instagram", Label: "Instagram", SortOrder: 3},
		{Platform: "Facebook", URL: "https://www.facebook.com/wegil", IconClass: "bxl bx-facebook", Label: "Facebook", SortOrder: 4},
	}
	db.Create(&links)
}

// seedTechStacks creates demo tech stacks. It is intentionally kept in the seed package because it depends on internal/model.
func seedTechStacks(db *gorm.DB) {
	stacks := []model.TechStack{
		{Name: "Go", Category: "Language", IconClass: "devicon-go-original-wordmark", IconURL: deviconCDN + "/go/go-original-wordmark.svg", Desc: "Primary backend language", SortOrder: 1},
		{Name: "PHP", Category: "Language", IconClass: "devicon-php-plain", IconURL: deviconCDN + "/php/php-plain.svg", Desc: "Web & backend scripting", SortOrder: 2},
		{Name: "JavaScript", Category: "Language", IconClass: "devicon-javascript-plain", IconURL: deviconCDN + "/javascript/javascript-plain.svg", Desc: "Frontend & backend scripting", SortOrder: 3},
		{Name: "Python", Category: "Language", IconClass: "devicon-python-plain", IconURL: deviconCDN + "/python/python-original.svg", Desc: "Scripting & automation", SortOrder: 4},
		{Name: "Rust", Category: "Language", IconClass: "devicon-rust-original", IconURL: deviconCDN + "/rust/rust-original.svg", Desc: "Systems programming & performance-critical code", SortOrder: 5},
		{Name: "C++", Category: "Language", IconClass: "devicon-cplusplus-plain", IconURL: deviconCDN + "/cplusplus/cplusplus-original.svg", Desc: "Legacy systems & performance optimization", SortOrder: 6},
		{Name: "Java", Category: "Language", IconClass: "devicon-java-plain", IconURL: deviconCDN + "/java/java-original.svg", Desc: "Enterprise applications & Android development", SortOrder: 7},
		{Name: "Fiber", Category: "Framework", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg", Desc: "Express-inspired Go web framework", URL: "https://gofiber.io", SortOrder: 1},
		{Name: "Gin", Category: "Framework", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png", Desc: "Minimalist Go web framework", URL: "https://gin-gonic.com", SortOrder: 2},
		{Name: "Code Igniter", Category: "Framework", IconClass: "devicon-codeigniter-plain", IconURL: deviconCDN + "/codeigniter/codeigniter-plain.svg", Desc: "Lightweight PHP framework", URL: "https://codeigniter.com", SortOrder: 3},
		{Name: "MySQL", Category: "Database", IconClass: "devicon-mysql-plain", IconURL: deviconCDN + "/mysql/mysql-original.svg", Desc: "Relational database", SortOrder: 1},
		{Name: "SQLite", Category: "Database", IconClass: "devicon-sqlite-plain", IconURL: deviconCDN + "/sqlite/sqlite-original.svg", Desc: "Embedded database", SortOrder: 2},
		{Name: "PostgreSQL", Category: "Database", IconClass: "devicon-postgresql-plain", IconURL: deviconCDN + "/postgresql/postgresql-original.svg", Desc: "Primary relational database", SortOrder: 3},
		{Name: "Redis", Category: "Database", IconClass: "devicon-redis-plain", IconURL: deviconCDN + "/redis/redis-original.svg", Desc: "Caching & sessions", SortOrder: 4},
		{Name: "MongoDB", Category: "Database", IconClass: "devicon-mongodb-plain", IconURL: deviconCDN + "/mongodb/mongodb-original.svg", Desc: "NoSQL document database", SortOrder: 5},
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", IconURL: deviconCDN + "/docker/docker-original.svg", Desc: "Containerization", SortOrder: 1},
		{Name: "Podman", Category: "DevOps", IconClass: "bxf bx-cube", IconURL: "https://cdn.jsdelivr.net/gh/containers/podman@main/logo/podman-logo-source.svg", IconURLDark: "https://cdn.jsdelivr.net/gh/containers/podman@main/logo/podman-logo-source.svg", Desc: "Alternative container engine", URL: "https://podman.io", SortOrder: 2},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", IconURL: deviconCDN + "/git/git-original.svg", Desc: "Version control", SortOrder: 3},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", IconURL: deviconCDN + "/linux/linux-original.svg", Desc: "Server OS", SortOrder: 4},
		{Name: "Nginx", Category: "DevOps", IconClass: "devicon-nginx-original", IconURL: deviconCDN + "/nginx/nginx-original.svg", Desc: "Web server & reverse proxy", SortOrder: 5},
		{Name: "Grafana", Category: "DevOps", IconClass: "bxf bx-line-chart", IconURL: "https://cdn.worldvectorlogo.com/logos/grafana.svg", Desc: "Monitoring & observability", URL: "https://grafana.com", SortOrder: 6},
		{Name: "n8n", Category: "Other Tools", IconClass: "bxf bx-cog", IconURL: "https://cdn.jsdelivr.net/gh/n8n-io/n8n@master/assets/n8n-logo.png", Desc: "Workflow automation tool", URL: "https://n8n.io", SortOrder: 1},
		{Name: "Postman", Category: "Other Tools", IconClass: "devicon-postman-plain", IconURL: deviconCDN + "/postman/postman-original.svg", Desc: "API documentation & testing tool", URL: "https://www.postman.com", SortOrder: 2},
		{Name: "ODOO", Category: "Other Tools", IconClass: "bxf bxl-odoo", IconURL: "https://cdn.worldvectorlogo.com/logos/odoo.svg", Desc: "ERP software for business management", URL: "https://www.odoo.com", SortOrder: 3},
	}
	db.Create(&stacks)
}

// seedComments creates demo comments and replies. It is intentionally kept in the seed package because it depends on internal/model.
func seedComments(db *gorm.DB) {
	users := []model.OAuthUser{
		{Provider: "github", ProviderID: "demo-1", Email: "alice@example.com", DisplayName: "Alice Chen", AvatarURL: "https://i.pravatar.cc/150?u=alice"},
		{Provider: "google", ProviderID: "demo-2", Email: "bob@example.com", DisplayName: "Bob Smith", AvatarURL: "https://i.pravatar.cc/150?u=bob"},
		{Provider: "system", ProviderID: "owner", Email: "john@example.com", DisplayName: "John Doe (Owner)"},
	}
	db.Create(&users)

	comments := []model.Comment{
		{OAuthUserID: users[0].ID, Body: "Great portfolio! I love the clean design and the tech stack section is really informative.", IsApproved: true},
		{OAuthUserID: users[1].ID, Body: "Impressive project list. The e-commerce platform looks really solid. Would love to see a demo!", IsApproved: true},
	}
	db.Create(&comments)

	replies := []model.Comment{
		{OAuthUserID: users[2].ID, ParentID: &comments[0].ID, Body: "Thank you Alice! Glad you like the design.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: users[2].ID, ParentID: &comments[1].ID, Body: "Thanks Bob! I'll add a live demo link soon.", IsOwnerReply: true, IsApproved: true},
	}
	db.Create(&replies)
}

// seedPosts creates demo blog posts. It is intentionally kept in the seed package because it depends on internal/model.
func seedPosts(db *gorm.DB) error {
	pub1 := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	pub2 := time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC)
	pub3 := time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)

	posts := []model.Post{
		{
			Title:       "Getting Started with Go and Fiber",
			Slug:        "getting-started-with-go-and-fiber",
			Excerpt:     "A beginner-friendly guide to building fast web apps with Go and the Fiber framework — from zero to your first REST API.",
			Content:     "## Why Go + Fiber?\n\nGo is one of the fastest-growing languages for backend development, and Fiber is an Express.js-inspired framework that makes it incredibly easy to get started.\n\n## Setting Up\n\nFirst, initialize your Go module:\n\n```bash\ngo mod init myapp\ngo get github.com/gofiber/fiber/v2\n```\n\n## Your First Route\n\n```go\napp := fiber.New()\napp.Get(\"/\", func(c fiber.Ctx) error {\n    return c.SendString(\"Hello, World!\")\n})\napp.Listen(\":3000\")\n```\n\n## What's Next?\n\nFrom here you can add middleware, connect a database with GORM, and build a full REST API. Go's performance and simplicity make it a great choice for modern backends.",
			Tags:        "Go,Fiber,Tutorial,Backend",
			Status:      "published",
			SortOrder:   1,
			PublishedAt: &pub1,
		},
		{
			Title:       "Building Dynamic UIs with HTMX — No JavaScript Framework Needed",
			Slug:        "building-dynamic-uis-with-htmx",
			Excerpt:     "How HTMX lets you add real-time interactivity to your pages with just HTML attributes, keeping things simple and fast.",
			Content:     "## What is HTMX?\n\nHTMX is a small (~14KB) JavaScript library that gives you access to AJAX, WebSockets, and server-sent events directly from HTML attributes — no framework needed.\n\n## A Simple Example\n\nLoad content without a full page refresh:\n\n```html\n<button hx-get=\"/api/data\" hx-target=\"#result\" hx-swap=\"innerHTML\">\n    Load Data\n</button>\n<div id=\"result\"></div>\n```\n\n## Why I Love It\n\nWith HTMX, I removed 80% of the custom JavaScript from this portfolio. The server renders HTML partials and HTMX swaps them in. It pairs perfectly with Go templates.\n\n## Pagination with HTMX\n\nThe \"Show More\" pattern is trivial — just return the next page of cards from the server and append them. No state management, no client-side routing.",
			Tags:        "HTMX,Frontend,HTML,Go",
			Status:      "published",
			SortOrder:   2,
			PublishedAt: &pub2,
		},
		{
			Title:       "My Development Workflow in 2025",
			Slug:        "my-development-workflow-2025",
			Excerpt:     "The tools, habits, and mindset behind how I build software day-to-day — from editor setup to deployment.",
			Content:     "## Editor\n\nI use **Neovim** with LSP for Go and TypeScript. It's fast, keyboard-driven, and highly customizable.\n\n## Version Control\n\nEvery project lives in Git. I follow conventional commits and keep branches small and focused.\n\n## Local Development\n\n- **Air** for hot-reload in Go projects\n- **Docker Compose** for local databases\n- **Make** for common commands (`make run`, `make build`, `make test`)\n\n## Deployment\n\nMost of my projects ship as a single Go binary behind **Nginx** on a Linux VPS. SQLite handles persistence for smaller apps; PostgreSQL for anything with real traffic.\n\n## Mindset\n\nShip early. Iterate fast. Keep dependencies minimal. The best code is the code you don't have to write.",
			Tags:        "Workflow,Tooling,Go,Productivity",
			Status:      "published",
			SortOrder:   3,
			PublishedAt: &pub3,
		},
	}
	return db.Create(&posts).Error
}

// seedUpcomingItems creates demo upcoming projects and announcements. It is intentionally kept in the seed package because it depends on internal/model.
func seedUpcomingItems(db *gorm.DB) {
	items := []model.UpcomingItem{
		{
			Title:       "Open Source CLI Tool",
			Description: "A developer productivity CLI written in Go — automates repetitive project scaffolding tasks and integrates with popular APIs.",
			Type:        "project",
			Status:      "in-progress",
			IconClass:   "bxf bx-terminal",
			SortOrder:   1,
			IsVisible:   true,
		},
		{
			Title:       "Mobile Companion App",
			Description: "A cross-platform mobile app built with Flutter to complement the portfolio. Includes push notifications and offline support.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-phone",
			SortOrder:   2,
			IsVisible:   true,
		},
		{
			Title:       "GraphQL API Service",
			Description: "Re-building the backend API layer with GraphQL on top of Go — type-safe, self-documenting, and ready for federation.",
			Type:        "project",
			Status:      "coming-soon",
			IconClass:   "bx bx-code-curly",
			SortOrder:   3,
			IsVisible:   true,
		},
		{
			Title:       "Tech Blog officially launches",
			Description: "This blog is going live with a dedicated series on Go, HTMX, and building side projects in public. Subscribe to get notified.",
			Type:        "announcement",
			Status:      "coming-soon",
			IconClass:   "bxf bx-megaphone",
			SortOrder:   4,
			IsVisible:   true,
		},
	}
	db.Create(&items)
}
