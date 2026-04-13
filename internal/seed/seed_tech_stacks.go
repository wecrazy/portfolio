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
		{Name: "JavaScript", Category: "Language", IconClass: "devicon-javascript-plain", IconURL: deviconCDN + "/javascript/javascript-original.svg", Desc: "Frontend & backend scripting", SortOrder: 3},
		{Name: "TypeScript", Category: "Language", IconClass: "devicon-typescript-plain", IconURL: deviconCDN + "/typescript/typescript-original.svg", Desc: "Typed JavaScript for scalable apps", SortOrder: 4},
		{Name: "Python", Category: "Language", IconClass: "devicon-python-plain", IconURL: deviconCDN + "/python/python-original.svg", Desc: "Scripting & automation", SortOrder: 5},
		{Name: "Java", Category: "Language", IconClass: "devicon-java-plain", IconURL: deviconCDN + "/java/java-original.svg", Desc: "Enterprise & Android development", SortOrder: 6},
		{Name: "C++", Category: "Language", IconClass: "devicon-cplusplus-plain", IconURL: deviconCDN + "/cplusplus/cplusplus-original.svg", Desc: "System programming & performance-critical apps", SortOrder: 7},
		{Name: "Rust", Category: "Language", IconClass: "devicon-rust-plain", IconURL: deviconCDN + "/rust/rust-original.svg", Desc: "Systems programming with safety", SortOrder: 8},
		{Name: "Bash", Category: "Language", IconClass: "devicon-bash-plain", IconURL: deviconCDN + "/bash/bash-original.svg", Desc: "Shell scripting for automation", SortOrder: 9},

		// ===================== Backend =====================
		{Name: "Fiber", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg", Desc: "Express-inspired Go web framework", URL: "https://gofiber.io", SortOrder: 1},
		{Name: "Gin", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png", Desc: "Minimalist Go web framework", URL: "https://gin-gonic.com", SortOrder: 2},
		{Name: "CodeIgniter", Category: "Backend", IconClass: "devicon-codeigniter-plain", IconURL: deviconCDN + "/codeigniter/codeigniter-plain.svg", Desc: "Lightweight PHP framework", URL: "https://codeigniter.com", SortOrder: 3},
		{Name: "GORM", Category: "Backend", IconClass: "bxf bx-data", IconURL: "https://media.brand.dev/196dd3de-5580-4788-9d77-708172bc9cda.svg", Desc: "ORM library for Golang", SortOrder: 4},
		{Name: "Concurrency (Goroutines)", Category: "Backend", IconClass: "bx bx-refresh-ccw-alt", Desc: "Lightweight thread managed by the Go runtime", SortOrder: 5},

		// ===================== API & Communication =====================
		{Name: "REST API", Category: "API & Communication", IconClass: "bxf bx-webhook", Desc: "RESTful API design & development", SortOrder: 1},
		{Name: "GraphQL", Category: "API & Communication", IconClass: "bx bx-code-alt", IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d0/GraphQL_logo_%28horizontal%29.svg/1280px-GraphQL_logo_%28horizontal%29.svg.png", Desc: "Flexible query language for APIs", SortOrder: 2},
		{Name: "gRPC", Category: "API & Communication", IconClass: "bxf bx-transfer", IconURL: "https://grpc.io/img/logos/grpc-icon-color.png", Desc: "High-performance RPC framework", SortOrder: 3},
		{Name: "Protocol Buffers", Category: "API & Communication", IconClass: "bxf bx-code", Desc: "Efficient serialization (Protobuf)", SortOrder: 4},
		{Name: "Middleware", Category: "API & Communication", IconClass: "bx bx-gear", Desc: "Request pipeline & handler chaining", SortOrder: 5},
		{Name: "JWT", Category: "API & Communication", IconClass: "bx bx-key", Desc: "Token-based authentication & authorization", SortOrder: 6},
		{Name: "OAuth 2.0", Category: "API & Communication", IconClass: "bx bx-lock-open", Desc: "Secure third-party authorization framework", SortOrder: 7},
		{Name: "Rate Limiting", Category: "API & Communication", IconClass: "bx bx-tachometer", Desc: "API traffic control & abuse prevention", SortOrder: 8},

		// ===================== System Design =====================
		{Name: "Monolith", Category: "System Design", IconClass: "bx bx-git-pull-request-draft", Desc: "Single application architecture", SortOrder: 1},
		{Name: "Microservices", Category: "System Design", IconClass: "bx bx-layers-down-left", Desc: "Distributed system architecture", SortOrder: 2},
		{Name: "Monorepo", Category: "System Design", IconClass: "bx bx-git-repo-forked", Desc: "Single repository for multiple services", SortOrder: 3},

		// ===================== Frontend =====================
		{Name: "React.js", Category: "Frontend", IconClass: "devicon-react-original", IconURL: deviconCDN + "/react/react-original.svg", Desc: "Component-based UI library", SortOrder: 1},
		{Name: "Next.js", Category: "Frontend", IconClass: "devicon-nextjs-original", IconURL: deviconCDN + "/nextjs/nextjs-original.svg", Desc: "Fullstack React framework", SortOrder: 2},
		{Name: "HTMX", Category: "Frontend", IconClass: "bxf bx-code-curly", IconURL: "https://logo.svgcdn.com/logos/htmx.png", Desc: "Server-driven UI interactions", SortOrder: 3},
		{Name: "Vite", Category: "Frontend", IconClass: "bxf bx-bolt", IconURL: "https://vite.dev/assets/vite-light.t8GCa_VF.svg", Desc: "Fast frontend tooling", SortOrder: 4},
		{Name: "SPA", Category: "Frontend", IconClass: "bx bx-book-content", Desc: "Single Page Application architecture", SortOrder: 5},

		// ===================== Data & Storage =====================
		{Name: "MySQL", Category: "Data & Storage", IconClass: "devicon-mysql-plain", IconURL: deviconCDN + "/mysql/mysql-original.svg", Desc: "Relational database", SortOrder: 1},
		{Name: "PostgreSQL", Category: "Data & Storage", IconClass: "devicon-postgresql-plain", IconURL: deviconCDN + "/postgresql/postgresql-original.svg", Desc: "Advanced relational database", SortOrder: 2},
		{Name: "SQLite", Category: "Data & Storage", IconClass: "devicon-sqlite-plain", IconURL: deviconCDN + "/sqlite/sqlite-original.svg", Desc: "Embedded database", SortOrder: 3},
		{Name: "Redis", Category: "Data & Storage", IconClass: "devicon-redis-plain", IconURL: deviconCDN + "/redis/redis-original.svg", Desc: "Caching & in-memory store", SortOrder: 4},
		{Name: "MongoDB", Category: "Data & Storage", IconClass: "devicon-mongodb-plain", IconURL: deviconCDN + "/mongodb/mongodb-original.svg", Desc: "NoSQL document database", SortOrder: 5},

		// ===================== DevOps & Infrastructure =====================
		{Name: "Docker", Category: "DevOps & Infrastructure", IconClass: "devicon-docker-plain", IconURL: deviconCDN + "/docker/docker-original.svg", Desc: "Containerization platform", SortOrder: 1},
		{Name: "Podman", Category: "DevOps & Infrastructure", IconClass: "bxf bx-cube", IconURL: deviconCDN + "/podman/podman-original.svg", Desc: "Daemonless container engine", SortOrder: 2},
		{Name: "Linux", Category: "DevOps & Infrastructure", IconClass: "devicon-linux-plain", IconURL: deviconCDN + "/linux/linux-original.svg", Desc: "Server operating system", SortOrder: 3},
		{Name: "Nginx", Category: "DevOps & Infrastructure", IconClass: "devicon-nginx-original", IconURL: deviconCDN + "/nginx/nginx-original.svg", Desc: "Web server & reverse proxy", SortOrder: 4},
		{Name: "CI/CD", Category: "DevOps & Infrastructure", IconClass: "bxf bx-git-branch", Desc: "Automated build & deployment pipelines", SortOrder: 5},
		{Name: "Git", Category: "DevOps & Infrastructure", IconClass: "devicon-git-plain", IconURL: deviconCDN + "/git/git-original.svg", Desc: "Version control system", SortOrder: 6},
		{Name: "GitHub", Category: "DevOps & Infrastructure", IconClass: "devicon-github-original", IconURL: deviconCDN + "/github/github-original.svg", Desc: "Code hosting & collaboration", SortOrder: 7},
		{Name: "cPanel (Server Management)", Category: "DevOps & Infrastructure", IconClass: "bx bx-server", IconURL: deviconCDN + "/cpanel/cpanel-original.svg", Desc: "Managing hosting, domains, databases, and deployments", SortOrder: 8},

		// ===================== Observability =====================
		{Name: "Prometheus", Category: "Observability", IconClass: "bxf bx-line-chart", IconURL: "https://cdn.worldvectorlogo.com/logos/prometheus.svg", Desc: "Metrics monitoring & alerting", SortOrder: 1},
		{Name: "Grafana", Category: "Observability", IconClass: "bxf bx-bar-chart-alt-2", IconURL: "https://cdn.worldvectorlogo.com/logos/grafana.svg", Desc: "Dashboards & visualization", SortOrder: 2},
		{Name: "Loki", Category: "Observability", IconClass: "bxf bx-file", IconURL: "https://grafana.com/static/img/logos/logo-loki.svg", Desc: "Log aggregation system", SortOrder: 3},
		{Name: "Tempo", Category: "Observability", IconClass: "bxf bx-timer", IconURL: "https://grafana.com/static/assets/img/logos/grafana-tempo.svg", Desc: "Distributed tracing backend", SortOrder: 4},

		// ===================== Integration & Messaging =====================
		{Name: "RabbitMQ", Category: "Integration & Messaging", IconClass: "bx bx-transfer", IconURL: "https://www.rabbitmq.com/img/rabbitmq-logo-by-tanzu.svg", Desc: "Message broker for async systems", SortOrder: 1},
		{Name: "Whatsmeow", Category: "Integration & Messaging", IconClass: "bxf bxl-whatsapp", IconURL: "https://pkg.go.dev/static/shared/logo/go-blue.svg", Desc: "Go library for WhatsApp Web multidevice API", SortOrder: 2},
		{Name: "WhatsApp Desktop", Category: "Integration & Messaging", IconClass: "bxf bxl-whatsapp", IconURL: "https://static.whatsapp.net/rsrc.php/yY/r/_mMwO8HKa4V.svg", Desc: "WhatsApp's official desktop client (Windows)", SortOrder: 1},
		{Name: "WhatsApp Business API", Category: "Integration & Messaging", IconClass: "bxf bxl-whatsapp", IconURL: "https://cdn.worldvectorlogo.com/logos/whatsapp.svg", Desc: "Official WhatsApp API by Meta", SortOrder: 3},
		{Name: "Telegram Bot", Category: "Integration & Messaging", IconClass: "bxf bxl-telegram", IconURL: "https://cdn.worldvectorlogo.com/logos/telegram.svg", Desc: "Telegram bot development", SortOrder: 4},
		{Name: "Email (SMTP/IMAP)", Category: "Integration & Messaging", IconClass: "bxf bx-mail-send", IconURL: "https://cdn-icons-png.flaticon.com/512/732/732200.png", Desc: "Email protocols for communication", SortOrder: 5},
		{Name: "Twilio", Category: "Integration & Messaging", IconClass: "bxf bx-phone", IconURL: "https://cdn.worldvectorlogo.com/logos/twilio.svg", Desc: "Communication APIs", SortOrder: 6},
		{Name: "n8n", Category: "Integration & Messaging", IconClass: "bxf bx-cog", IconURL: "https://n8n.io/brandguidelines/logo-dark.svg", Desc: "Workflow automation", URL: "https://n8n.io", SortOrder: 7},
		{Name: "Odoo", Category: "Integration & Messaging", IconClass: "bxf bxl-odoo", IconURL: "https://cdn.worldvectorlogo.com/logos/odoo.svg", Desc: "ERP system", URL: "https://www.odoo.com", SortOrder: 8},

		// ===================== Testing & Performance =====================
		{Name: "k6", Category: "Testing & Performance", IconClass: "bxf bx-tachometer", IconURL: "https://grafana.com/media/docs/k6/GrafanaLogo_k6_orange_icon.svg", Desc: "Modern load testing tool by Grafana Labs", SortOrder: 1},
		{Name: "JMeter", Category: "Testing & Performance", IconClass: "bxf bx-bar-chart", IconURL: "https://jmeter.apache.org/images/logo.svg", Desc: "Performance testing by Apache", SortOrder: 2},
		{Name: "Testify", Category: "Testing & Performance", IconClass: "bx bx-medical-flask", Desc: "Testing toolkit for Go (assertions, mocking, suites)", SortOrder: 3},

		// ===================== Developer Tools =====================
		{Name: "Postman", Category: "Developer Tools", IconClass: "devicon-postman-plain", IconURL: deviconCDN + "/postman/postman-original.svg", Desc: "API testing & development tool", SortOrder: 1},
		{Name: "Swagger", Category: "Developer Tools", IconClass: "bxf bx-book", IconURL: "https://static1.smartbear.co/swagger/media/assets/images/swagger_logo.svg", Desc: "OpenAPI documentation tool", SortOrder: 2},
		{Name: "CLI Development", Category: "Developer Tools", IconClass: "bx bx-terminal", Desc: "Building command-line tools and developer utilities", SortOrder: 3},
		{Name: "TUI (Bubble Tea)", Category: "Developer Tools", IconClass: "bx bx-window", IconURL: "https://charm.land/bubbletea-light.41979931daa0fa73.webp", Desc: "Terminal UI development using Bubble Tea by Charm", SortOrder: 4},
	}

	db.Create(&stacks)
}
