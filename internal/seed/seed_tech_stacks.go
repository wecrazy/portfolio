package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

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
