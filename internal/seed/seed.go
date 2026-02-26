// Package seed provides initial data seeding for the database.
package seed

import (
	"log"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedIfNeeded creates the default admin user and an empty owner profile when
// they don't already exist. Safe to call on every startup.
func SeedIfNeeded(db *gorm.DB, cfg config.TypeMyPortfolio) {
	seedAdmin(db, cfg)
	seedOwner(db)
	seedDemoData(db)
}

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

func seedOwner(db *gorm.DB) {
	var count int64
	db.Model(&model.Owner{}).Count(&count)
	if count > 0 {
		return
	}

	owner := model.Owner{
		FullName: "John Doe",
		Title:    "Full-Stack Developer",
		Bio:      "Passionate full-stack developer with 5+ years of experience building modern web applications. I love turning complex problems into simple, beautiful, and intuitive solutions.",
		Email:    "john@example.com",
		Phone:    "+62 812 3456 7890",
		Location: "Jakarta, Indonesia",
	}
	if err := db.Create(&owner).Error; err != nil {
		log.Fatalf("Failed to seed owner profile: %v", err)
	}
	log.Println("Seeded default owner profile")
}

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

func seedProjects(db *gorm.DB) {
	projects := []model.Project{
		{Title: "E-Commerce Platform", Slug: "e-commerce-platform", Description: "A full-featured e-commerce platform with cart, checkout, and payment integration.", Tags: "Go,React,PostgreSQL,Stripe", Status: "published", SortOrder: 1, Featured: true, LiveURL: "https://example.com", RepoURL: "https://github.com/example/ecommerce"},
		{Title: "Task Management App", Slug: "task-management-app", Description: "Real-time task management application with team collaboration features.", Tags: "TypeScript,Next.js,Prisma,WebSocket", Status: "published", SortOrder: 2, Featured: true, RepoURL: "https://github.com/example/tasks"},
		{Title: "Weather Dashboard", Slug: "weather-dashboard", Description: "Beautiful weather dashboard with 7-day forecast, radar maps, and location search.", Tags: "Vue.js,OpenWeather API,Chart.js", Status: "published", SortOrder: 3, LiveURL: "https://example.com/weather"},
		{Title: "Blog Engine", Slug: "blog-engine", Description: "Markdown-powered blog engine with SEO optimization and RSS feed.", Tags: "Go,Fiber,SQLite,HTMX", Status: "published", SortOrder: 4, RepoURL: "https://github.com/example/blog"},
		{Title: "Chat Application", Slug: "chat-application", Description: "Real-time chat app with rooms, direct messages, and file sharing.", Tags: "Go,WebSocket,Redis,React", Status: "published", SortOrder: 5, Featured: true},
		{Title: "Portfolio Builder", Slug: "portfolio-builder", Description: "Drag-and-drop portfolio builder for developers with custom themes.", Tags: "Next.js,Tailwind,MongoDB", Status: "published", SortOrder: 6},
		{Title: "API Gateway", Slug: "api-gateway", Description: "High-performance API gateway with rate limiting, caching, and auth.", Tags: "Go,Redis,Docker,gRPC", Status: "published", SortOrder: 7, RepoURL: "https://github.com/example/gateway"},
		{Title: "Mobile Fitness App", Slug: "mobile-fitness-app", Description: "Cross-platform fitness tracking app with workout plans and progress charts.", Tags: "Flutter,Firebase,Dart", Status: "published", SortOrder: 8},
	}
	db.Create(&projects)
}

func seedExperiences(db *gorm.DB) {
	past1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	past2 := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	past2End := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	past3 := time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)
	past3End := time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC)
	edu1 := time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC)
	edu1End := time.Date(2019, 6, 30, 0, 0, 0, 0, time.UTC)

	experiences := []model.Experience{
		{Type: "Work", Title: "Senior Full-Stack Developer", Org: "Tech Corp", Location: "Jakarta, Indonesia", StartDate: past1, IsCurrent: true, Description: "Leading a team of 5 developers building microservices architecture. Implemented CI/CD pipelines that reduced deployment time by 60%.", SortOrder: 1},
		{Type: "Work", Title: "Full-Stack Developer", Org: "StartupXYZ", Location: "Bandung, Indonesia", StartDate: past2, EndDate: &past2End, Description: "Built and maintained the core product platform serving 50K+ users. Migrated legacy monolith to Go microservices.", SortOrder: 2},
		{Type: "Work", Title: "Junior Developer", Org: "WebAgency", Location: "Yogyakarta, Indonesia", StartDate: past3, EndDate: &past3End, Description: "Developed responsive web applications for various clients. Gained experience with React, Node.js, and PostgreSQL.", SortOrder: 3},
		{Type: "Education", Title: "Bachelor of Computer Science", Org: "Universitas Indonesia", Location: "Jakarta, Indonesia", StartDate: edu1, EndDate: &edu1End, Description: "Graduated with honors. Focus on software engineering and distributed systems.", SortOrder: 4},
	}
	db.Create(&experiences)
}

func seedSkills(db *gorm.DB) {
	skills := []model.Skill{
		{Name: "Go", Category: "Languages", IconClass: "devicon-go-original-wordmark", Proficiency: 95, SortOrder: 1},
		{Name: "TypeScript", Category: "Languages", IconClass: "devicon-typescript-plain", Proficiency: 90, SortOrder: 2},
		{Name: "Python", Category: "Languages", IconClass: "devicon-python-plain", Proficiency: 80, SortOrder: 3},
		{Name: "JavaScript", Category: "Languages", IconClass: "devicon-javascript-plain", Proficiency: 90, SortOrder: 4},
		{Name: "React", Category: "Frontend", IconClass: "devicon-react-original", Proficiency: 85, SortOrder: 1},
		{Name: "Vue.js", Category: "Frontend", IconClass: "devicon-vuejs-plain", Proficiency: 75, SortOrder: 2},
		{Name: "HTMX", Category: "Frontend", IconClass: "bi bi-lightning-charge-fill", Proficiency: 90, SortOrder: 3},
		{Name: "Tailwind CSS", Category: "Frontend", IconClass: "devicon-tailwindcss-original", Proficiency: 85, SortOrder: 4},
		{Name: "Fiber", Category: "Backend", IconClass: "bi bi-lightning-fill", Proficiency: 95, SortOrder: 1},
		{Name: "PostgreSQL", Category: "Backend", IconClass: "devicon-postgresql-plain", Proficiency: 85, SortOrder: 2},
		{Name: "Redis", Category: "Backend", IconClass: "devicon-redis-plain", Proficiency: 80, SortOrder: 3},
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", Proficiency: 85, SortOrder: 1},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", Proficiency: 80, SortOrder: 2},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", Proficiency: 90, SortOrder: 3},
	}
	db.Create(&skills)
}

func seedSocialLinks(db *gorm.DB) {
	links := []model.SocialLink{
		{Platform: "GitHub", URL: "https://github.com/johndoe", IconClass: "fab fa-github", Label: "GitHub", SortOrder: 1},
		{Platform: "LinkedIn", URL: "https://linkedin.com/in/johndoe", IconClass: "fab fa-linkedin", Label: "LinkedIn", SortOrder: 2},
		{Platform: "Instagram", URL: "https://instagram.com/johndoe", IconClass: "fab fa-instagram", Label: "Instagram", SortOrder: 3},
		{Platform: "Twitter", URL: "https://twitter.com/johndoe", IconClass: "fab fa-x-twitter", Label: "Twitter", SortOrder: 4},
	}
	db.Create(&links)
}

func seedTechStacks(db *gorm.DB) {
	stacks := []model.TechStack{
		{Name: "Go", Category: "Language", IconClass: "devicon-go-original-wordmark", Desc: "Primary backend language", SortOrder: 1},
		{Name: "TypeScript", Category: "Language", IconClass: "devicon-typescript-plain", Desc: "Frontend & full-stack", SortOrder: 2},
		{Name: "Python", Category: "Language", IconClass: "devicon-python-plain", Desc: "Scripting & automation", SortOrder: 3},
		{Name: "Fiber", Category: "Framework", IconClass: "bi bi-lightning-fill", Desc: "Express-inspired Go web framework", URL: "https://gofiber.io", SortOrder: 1},
		{Name: "React", Category: "Framework", IconClass: "devicon-react-original", Desc: "UI component library", SortOrder: 2},
		{Name: "Next.js", Category: "Framework", IconClass: "devicon-nextjs-plain", Desc: "React meta-framework", SortOrder: 3},
		{Name: "PostgreSQL", Category: "Database", IconClass: "devicon-postgresql-plain", Desc: "Primary relational database", SortOrder: 1},
		{Name: "SQLite", Category: "Database", IconClass: "devicon-sqlite-plain", Desc: "Embedded database", SortOrder: 2},
		{Name: "Redis", Category: "Database", IconClass: "devicon-redis-plain", Desc: "Caching & sessions", SortOrder: 3},
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", Desc: "Containerization", SortOrder: 1},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", Desc: "Version control", SortOrder: 2},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", Desc: "Server OS", SortOrder: 3},
	}
	db.Create(&stacks)
}

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

func seedPosts(db *gorm.DB) error {
	pub1 := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	pub2 := time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC)
	pub3 := time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)

	posts := []model.Post{
		{
			Title:       "Getting Started with Go and Fiber",
			Slug:        "getting-started-with-go-and-fiber",
			Excerpt:     "A beginner-friendly guide to building fast web apps with Go and the Fiber framework — from zero to your first REST API.",
			Content:     "## Why Go + Fiber?\n\nGo is one of the fastest-growing languages for backend development, and Fiber is an Express.js-inspired framework that makes it incredibly easy to get started.\n\n## Setting Up\n\nFirst, initialize your Go module:\n\n```bash\ngo mod init myapp\ngo get github.com/gofiber/fiber/v2\n```\n\n## Your First Route\n\n```go\napp := fiber.New()\napp.Get(\"/\", func(c *fiber.Ctx) error {\n    return c.SendString(\"Hello, World!\")\n})\napp.Listen(\":3000\")\n```\n\n## What's Next?\n\nFrom here you can add middleware, connect a database with GORM, and build a full REST API. Go's performance and simplicity make it a great choice for modern backends.",
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

func seedUpcomingItems(db *gorm.DB) {
	items := []model.UpcomingItem{
		{
			Title:       "Open Source CLI Tool",
			Description: "A developer productivity CLI written in Go — automates repetitive project scaffolding tasks and integrates with popular APIs.",
			Type:        "project",
			Status:      "in-progress",
			IconClass:   "bi bi-terminal-fill",
			SortOrder:   1,
			IsVisible:   true,
		},
		{
			Title:       "Mobile Companion App",
			Description: "A cross-platform mobile app built with Flutter to complement the portfolio. Includes push notifications and offline support.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bi bi-phone-fill",
			SortOrder:   2,
			IsVisible:   true,
		},
		{
			Title:       "GraphQL API Service",
			Description: "Re-building the backend API layer with GraphQL on top of Go — type-safe, self-documenting, and ready for federation.",
			Type:        "project",
			Status:      "coming-soon",
			IconClass:   "bi bi-braces",
			SortOrder:   3,
			IsVisible:   true,
		},
		{
			Title:       "Tech Blog officially launches",
			Description: "This blog is going live with a dedicated series on Go, HTMX, and building side projects in public. Subscribe to get notified.",
			Type:        "announcement",
			Status:      "coming-soon",
			IconClass:   "bi bi-megaphone-fill",
			SortOrder:   4,
			IsVisible:   true,
		},
	}
	db.Create(&items)
}
