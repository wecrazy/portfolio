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

	csnaDesc := `* Led development of backend systems, web apps, and automation solutions with focus on scalability and reliability
* Designed and customized Odoo ERP modules (automation, reporting, dashboards, integrations) for internal and client use
* Acted as App Lead for Call Center, building system integrating CS dashboard with Odoo for task handling and communication flows
* Developed Call Center app with WhatsApp Desktop integration for direct chat/call with customers
* Supervised Technical Assistance division, creating validation systems and dashboards for technician visit data
* Built technician tracking dashboards (PHP, DataTables.js, Leaflet.js) with performance metrics and geolocation
* Engineered automated Warning Letter (SP) system with scheduled checks and multi-channel alerts (Email, WhatsApp, Telegram)
* Developed multi-language chatbot integrated with WhatsApp, Telegram, Twilio for ticket status, reports, and internal tools
* Built configurable chatbot dashboard for dynamic flows and automation without code changes
* Automated technician payslip generation from Excel templates with delivery via Email/WhatsApp
* Developed multi-app platform with role-based access using Go & Gin / Fiber.
* Used PostgreSQL, MySQL, MongoDB, SQLite, Redis (caching) based on system needs
* Implemented microservices communication using gRPC and Protocol Buffers
* Managed Linux servers, Nginx, and deployment environments
* Containerized apps using Podman/Docker
* Documented APIs with Swagger/OpenAPI & Postman
* Performed testing (k6, JMeter, Postman, Golang unit/integration)
* Implemented observability with Grafana, Prometheus, Loki, Tempo
* Built frontend apps using HTMX, React.js (Vite) and Next.js`

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
