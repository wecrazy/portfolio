package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedTechStacks creates demo tech stacks. It is intentionally kept in the seed package because it depends on internal/model.
func seedTechStacks(db *gorm.DB) {
	stacks := []model.TechStack{

		// ===================== Languages =====================
		{Name: "Go", Category: "Language", IconClass: "devicon-go-original-wordmark", IconURL: deviconCDN + "/go/go-original-wordmark.svg", Desc: "Primary backend language", SortOrder: 1},
		{Name: "PHP", Category: "Language", IconClass: "devicon-php-plain", IconURL: deviconCDN + "/php/php-original.svg", Desc: "Web & backend scripting", SortOrder: 2},
		{Name: "JavaScript", Category: "Language", IconClass: "devicon-javascript-plain", IconURL: deviconCDN + "/javascript/javascript-plain.svg", Desc: "Frontend & backend scripting", SortOrder: 3},
		{Name: "TypeScript", Category: "Language", IconClass: "devicon-typescript-plain", IconURL: deviconCDN + "/typescript/typescript-original.svg", Desc: "Typed JavaScript for scalable apps", SortOrder: 4},
		{Name: "Python", Category: "Language", IconClass: "devicon-python-plain", IconURL: deviconCDN + "/python/python-original.svg", Desc: "Scripting & automation", SortOrder: 5},
		{Name: "Rust", Category: "Language", IconClass: "devicon-rust-original", IconURL: deviconCDN + "/rust/rust-original.svg", Desc: "Systems programming", SortOrder: 6},
		{Name: "C++", Category: "Language", IconClass: "devicon-cplusplus-plain", IconURL: deviconCDN + "/cplusplus/cplusplus-original.svg", Desc: "Performance & low-level systems", SortOrder: 7},
		{Name: "Java", Category: "Language", IconClass: "devicon-java-plain", IconURL: deviconCDN + "/java/java-original.svg", Desc: "Enterprise applications", SortOrder: 8},

		// ===================== Backend & Framework =====================
		{Name: "Fiber", Category: "Framework", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg", Desc: "Express-inspired Go framework", URL: "https://gofiber.io", SortOrder: 1},
		{Name: "Gin", Category: "Framework", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png", Desc: "Minimalist Go framework", URL: "https://gin-gonic.com", SortOrder: 2},
		{Name: "CodeIgniter", Category: "Framework", IconClass: "devicon-codeigniter-plain", IconURL: deviconCDN + "/codeigniter/codeigniter-plain.svg", Desc: "Lightweight PHP framework", URL: "https://codeigniter.com", SortOrder: 3},
		{Name: "GORM", Category: "Framework", IconClass: "bxf bx-data", IconURL: "https://media.brand.dev/196dd3de-5580-4788-9d77-708172bc9cda.svg", Desc: "ORM library for Golang", SortOrder: 4},

		// ===================== Architecture =====================
		{Name: "REST API", Category: "Architecture", IconClass: "bxf bx-network-chart", IconURL: "https://latenode.com/_next/image?url=https%3A%2F%2Fblog-static.latenode.com%2Flatenode-strapi-blog%2Fhow_to_reduce_api_latency_in_integrations_featured_89ab3e0e7e.jpg&w=384&q=75", Desc: "RESTful API design & development", SortOrder: 1},
		{Name: "gRPC", Category: "Architecture", IconClass: "bxf bx-transfer", IconURL: "https://grpc.io/img/logos/grpc-icon-color.png", Desc: "High-performance RPC framework", SortOrder: 2},
		{Name: "Protocol Buffers", Category: "Architecture", IconClass: "bxf bx-code-block", IconURL: "https://dz2cdn1.dzone.com/storage/temp/15915735-1653483592719.png", Desc: "Efficient serialization (Protobuf)", SortOrder: 3},
		{Name: "Monolith", Category: "Architecture", IconClass: "bx bx-git-pull-request-draft", Desc: "Single application architecture", SortOrder: 4},
		{Name: "Monorepo", Category: "Architecture", IconClass: "bx bx-git-repo-forked", Desc: "Single repository for multiple services", SortOrder: 5},
		{Name: "Microservices", Category: "Architecture", IconClass: "bx bx-cube", Desc: "Distributed system architecture", SortOrder: 6},
		{Name: "Middleware", Category: "Architecture", IconClass: "bx bx-gear", Desc: "Request pipeline handling", SortOrder: 7},
		{Name: "GraphQL", Category: "Architecture", IconClass: "bx bx-code-alt", IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d0/GraphQL_logo_%28horizontal%29.svg/1280px-GraphQL_logo_%28horizontal%29.svg.png", Desc: "Flexible query language for APIs", SortOrder: 8},

		// ===================== Frontend =====================
		{Name: "React.js", Category: "Frontend", IconClass: "devicon-react-original", IconURL: deviconCDN + "/react/react-original.svg", Desc: "Component-based UI library", SortOrder: 1},
		{Name: "Next.js", Category: "Frontend", IconClass: "devicon-nextjs-original", IconURL: deviconCDN + "/nextjs/nextjs-original.svg", Desc: "Fullstack React framework", SortOrder: 2},
		{Name: "HTMX", Category: "Frontend", IconClass: "bxf bx-code-curly", IconURL: "https://logo.svgcdn.com/logos/htmx.png", Desc: "Server-driven UI interactions", SortOrder: 3},
		{Name: "Vite", Category: "Frontend", IconClass: "bxf bx-bolt", IconURL: "https://vite.dev/assets/vite-light.t8GCa_VF.svg", Desc: "Fast frontend tooling", SortOrder: 4},
		{Name: "Single Page Application (SPA)", Category: "Frontend", IconClass: "bx bx-book-content", Desc: "Client-side rendered web application architecture", SortOrder: 5},

		// ===================== Databases =====================
		{Name: "MySQL", Category: "Database", IconClass: "devicon-mysql-plain", IconURL: deviconCDN + "/mysql/mysql-original.svg", Desc: "Relational database", SortOrder: 1},
		{Name: "SQLite", Category: "Database", IconClass: "devicon-sqlite-plain", IconURL: deviconCDN + "/sqlite/sqlite-original.svg", Desc: "Embedded database", SortOrder: 2},
		{Name: "PostgreSQL", Category: "Database", IconClass: "devicon-postgresql-plain", IconURL: deviconCDN + "/postgresql/postgresql-original.svg", Desc: "Advanced relational database", SortOrder: 3},
		{Name: "Redis", Category: "Database", IconClass: "devicon-redis-plain", IconURL: deviconCDN + "/redis/redis-original.svg", Desc: "Caching & in-memory store", SortOrder: 4},
		{Name: "MongoDB", Category: "Database", IconClass: "devicon-mongodb-plain", IconURL: deviconCDN + "/mongodb/mongodb-original.svg", Desc: "NoSQL document database", SortOrder: 5},

		// ===================== DevOps =====================
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", IconURL: deviconCDN + "/docker/docker-original.svg", Desc: "Containerization platform", SortOrder: 1},
		{Name: "Podman", Category: "DevOps", IconClass: "bxf bx-cube", IconURL: "https://cdn.jsdelivr.net/gh/containers/podman@main/logo/podman-logo-source.svg", Desc: "Daemonless container engine", SortOrder: 2},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", IconURL: deviconCDN + "/git/git-original.svg", Desc: "Version control", SortOrder: 3},
		{Name: "GitHub", Category: "DevOps", IconClass: "devicon-github-original", IconURL: deviconCDN + "/github/github-original.svg", Desc: "Code hosting platform", SortOrder: 4},
		{Name: "CI/CD", Category: "DevOps", IconClass: "bxf bx-git-branch", Desc: "Automation pipelines", SortOrder: 5},
		{Name: "RabbitMQ", Category: "DevOps", IconClass: "bx bx-transfer", IconURL: "https://www.rabbitmq.com/img/rabbitmq-logo-by-tanzu.svg", Desc: "Message broker for async systems", SortOrder: 6},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", IconURL: deviconCDN + "/linux/linux-original.svg", Desc: "Server OS", SortOrder: 7},
		{Name: "Nginx", Category: "DevOps", IconClass: "devicon-nginx-original", IconURL: deviconCDN + "/nginx/nginx-original.svg", Desc: "Web server & reverse proxy", SortOrder: 8},

		// ===================== Observability =====================
		{Name: "Prometheus", Category: "Observability", IconClass: "bxf bx-line-chart", IconURL: "https://cdn.worldvectorlogo.com/logos/prometheus.svg", Desc: "Metrics monitoring", SortOrder: 1},
		{Name: "Grafana", Category: "Observability", IconClass: "bxf bx-bar-chart-alt-2", IconURL: "https://cdn.worldvectorlogo.com/logos/grafana.svg", Desc: "Dashboards & visualization", SortOrder: 2},
		{Name: "Loki", Category: "Observability", IconClass: "bxf bx-file", IconURL: "https://grafana.com/static/img/logos/logo-loki.svg", Desc: "Log aggregation", SortOrder: 3},
		{Name: "Tempo", Category: "Observability", IconClass: "bxf bx-timer", IconURL: "https://grafana.com/static/assets/img/logos/grafana-tempo.svg", Desc: "Distributed tracing", SortOrder: 4},

		// ===================== Integration =====================
		{Name: "Whatsmeow", Category: "Integration", IconClass: "bxf bxl-whatsapp", IconURL: "https://pkg.go.dev/static/shared/logo/go-blue.svg", Desc: "whatsmeow is a Go library for the WhatsApp web multidevice API", SortOrder: 0},
		{Name: "WhatsApp Desktop", Category: "Integration", IconClass: "bxf bxl-whatsapp", IconURL: "https://static.whatsapp.net/rsrc.php/yY/r/_mMwO8HKa4V.svg", Desc: "WhatsApp's official desktop client (Windows)", SortOrder: 1},
		{Name: "WhatsApp Business API", Category: "Integration", IconClass: "bxf bxl-whatsapp", IconURL: "https://cdn.worldvectorlogo.com/logos/whatsapp.svg", Desc: "WhatsApp's official business API by Meta", SortOrder: 2},
		{Name: "Email (SMTP/IMAP)", Category: "Integration", IconClass: "bxf bx-mail-send", IconURL: "https://cdn-icons-png.flaticon.com/512/732/732200.png", Desc: "Email protocols for sending and receiving emails", SortOrder: 3},
		{Name: "Telegram Bot", Category: "Integration", IconClass: "bxf bxl-telegram", IconURL: "https://cdn.worldvectorlogo.com/logos/telegram.svg", Desc: "Telegram bot development", SortOrder: 4},
		{Name: "Twilio", Category: "Integration", IconClass: "bxf bx-phone", IconURL: "https://cdn.worldvectorlogo.com/logos/twilio.svg", Desc: "Communication APIs", SortOrder: 5},
		{Name: "n8n", Category: "Integration", IconClass: "bxf bx-cog", IconURL: "https://cdn.jsdelivr.net/gh/n8n-io/n8n@master/assets/n8n-logo.png", Desc: "Workflow automation", URL: "https://n8n.io", SortOrder: 6},
		{Name: "Odoo", Category: "Integration", IconClass: "bxf bxl-odoo", IconURL: "https://cdn.worldvectorlogo.com/logos/odoo.svg", Desc: "ERP system", URL: "https://www.odoo.com", SortOrder: 7},

		// ===================== Tools =====================
		{Name: "Postman", Category: "Tools", IconClass: "devicon-postman-plain", IconURL: deviconCDN + "/postman/postman-original.svg", Desc: "API testing tool", SortOrder: 1},
		{Name: "Swagger", Category: "Tools", IconClass: "bxf bx-book", IconURL: "https://static1.smartbear.co/swagger/media/assets/images/swagger_logo.svg", Desc: "OpenAPI documentation", SortOrder: 2},
		{Name: "k6", Category: "Tools", IconClass: "bxf bx-tachometer", IconURL: "https://grafana.com/media/docs/k6/GrafanaLogo_k6_orange_icon.svg", Desc: "Load testing", SortOrder: 3},
		{Name: "JMeter", Category: "Tools", IconClass: "bxf bx-bar-chart", IconURL: "https://jmeter.apache.org/images/logo.svg", Desc: "Performance testing by Apache JMeter", SortOrder: 4},
	}

	db.Create(&stacks)
}
