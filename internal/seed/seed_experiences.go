package seed

import (
	"time"

	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedExperiences creates demo work and education experiences. It is intentionally kept in the seed package because it depends on internal/model.
func seedExperiences(db *gorm.DB) {
	workStart := time.Date(2023, time.February, 13, 0, 0, 0, 0, time.UTC)
	workEnd := time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC)

	eduStart := time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2022, time.September, 23, 0, 0, 0, 0, time.UTC)

	internStart := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	internEnd := time.Date(2021, time.February, 28, 0, 0, 0, 0, time.UTC)

	csnaDesc := `* Led development of scalable backend systems, web applications, and automation platforms, improving operational efficiency and system reliability
* Designed and customized Odoo ERP modules (automation, reporting, dashboards, integrations) for internal and client workflows

Technical Leadership
* Technical Assistance Lead (Feb 2025 – Apr 2026): Built validation middleware and TA dashboard for managed service workflows, including AI-assisted photo validation and Odoo-based verification. Enabled error/pending tracking and performance visibility
* Call Center Lead (Apr 2025 – Apr 2026): Led development of Call Center platform integrated with Odoo ERP, enabling task distribution, automated customer interaction (WhatsApp/call), and guided agent workflows. Delivered reporting and continuous improvements

Key Contributions
* Built multi-language chatbot platforms (WhatsApp, Telegram, Twilio, Whatsmeow) with configurable automation flows
* Engineered automation systems including Warning Letter (SP) and technician payslip generation with multi-channel delivery
* Developed multi-application platform with role-based access using Go (Gin/Fiber)

Tech Stack
* Languages: Go, PHP, JavaScript/TypeScript, Python
* Backend & Architecture: Gin, Fiber, GORM, REST, gRPC, Microservices, Monolith, Monorepo
* Data: PostgreSQL, MySQL, MongoDB, Redis, SQLite
* Infra & DevOps: Docker/Podman, Linux, Nginx, CI/CD
* Integration: Odoo, WhatsApp (Business API/Desktop), Telegram, Twilio, n8n, RabbitMQ, Email
* Observability: Grafana, Prometheus, Loki, Tempo
* Testing & Tools: Swagger/OpenAPI, Postman, k6, JMeter, Testify
* Frontend: HTMX, React (Vite), Next.js`

	distanEnrekangDesc := `* Supported the preparation and management of distribution reports for fertilizers and agricultural aid
* Assisted in organizing and validating operational data
* Contributed to the development of a simple informational website - Participated in technology outreach initiatives
* Collaborated to improve data-driven workflows and reporting`

	eduDesc := `Graduated Cum Laude with a GPA of 3.89, earning a Bachelor’s degree in Informatics Engineering (S.Kom) from the Faculty of Engineering, Universitas Kristen Indonesia Toraja in 2022, and was recognized as the best graduate of the program.`

	experiences := []model.Experience{
		{
			Type:        "Work",
			Title:       "Fullstack Developer",
			Org:         "PT. Cyber Smart Network Asia",
			Location:    "East Jakarta, Special Capital Region of Jakarta, Indonesia",
			StartDate:   workStart,
			EndDate:     &workEnd,
			IsCurrent:   false,
			Description: csnaDesc,
			SortOrder:   1,
			ImageURL:    "https://images.glints.com/unsafe/glints-dashboard.oss-ap-southeast-1-internal.aliyuncs.com/company-logo/76aa4598c9297e7d44202c32545eaeca.jpg",
		},
		{
			Type:        "Work",
			Title:       "IT Internship",
			Org:         "Department of Agriculture and Plantation, Enrekang Regency",
			Location:    "Enrekang Regency, South Sulawesi, Indonesia",
			StartDate:   internStart,
			EndDate:     &internEnd,
			IsCurrent:   false,
			Description: distanEnrekangDesc,
			SortOrder:   2,
			ImageURL:    "https://distanbun.wordpress.com/wp-content/uploads/2011/07/tanji.jpg",
		},

		{
			Type:        "Education",
			Title:       "Bachelor of Informatics Engineering",
			Org:         "Universitas Kristen Indonesia Toraja",
			Location:    "Tana Toraja, South Sulawesi, Indonesia",
			StartDate:   eduStart,
			EndDate:     &eduEnd,
			IsCurrent:   false,
			Description: eduDesc,
			SortOrder:   3,
			CertURL:     "https://drive.google.com/file/d/1A0hDQingopqunPV87WMY93O3JOWjg92I/view?usp=sharing",
			ImageURL:    "https://ukitoraja.ac.id/wp-content/uploads/2019/05/Logo-UKIT.png",
		},
	}
	db.Create(&experiences)
}
