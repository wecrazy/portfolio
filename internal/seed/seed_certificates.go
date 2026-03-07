package seed

import (
	"time"

	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedCertificates inserts a handful of sample certificates into the database
// so that the public certificate grid renders something during development.
// The function is idempotent; it does nothing if any certificates already
// exist.
func seedCertificates(db *gorm.DB) {
	var count int64
	db.Model(&model.Certificate{}).Count(&count)
	if count > 0 {
		return
	}

	certs := []model.Certificate{
		{
			Title:     "ISO 27001",
			Issuer:    "AKSESI Certification",
			IssueDate: time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC),
			CertURL:   "https://drive.google.com/file/d/1LEiesBoYKBSJqZGreBazhkafhhRpXG17/view?usp=sharing",
			IsVisible: true,
			SortOrder: 1,
		},
		{
			Title:     "Best Graduate",
			Issuer:    "Universitas Kristen Indonesia Toraja",
			IssueDate: time.Date(2022, time.September, 2, 0, 0, 0, 0, time.UTC),
			CertURL:   "https://drive.google.com/file/d/1A0hDQingopqunPV87WMY93O3JOWjg92I/view?usp=sharing",
			IsVisible: true,
			SortOrder: 2,
		},
	}

	db.Create(&certs)
}
