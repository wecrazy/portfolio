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
	eduStart := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC)

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
			Description: "Graduated Cum Laude from the Faculty of Engineering, Department of Informatics Engineering with a GPA of 3.89, earning recognition as the best graduate of the class",
			SortOrder:   2,
			CertURL:     "https://drive.google.com/file/d/1A0hDQingopqunPV87WMY93O3JOWjg92I/view?usp=drive_link",
			ImageURL:    "https://ukitoraja.ac.id/wp-content/uploads/2019/05/Logo-UKIT.png",
		},
	}
	db.Create(&experiences)
}
