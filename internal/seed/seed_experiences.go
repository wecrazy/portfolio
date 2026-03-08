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
	internStart := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	internEnd := time.Date(2021, time.February, 28, 0, 0, 0, 0, time.UTC)

	experiences := []model.Experience{
		{
			Type:        "Work",
			Title:       "Full Stack Programmer",
			Org:         "PT. Smartweb Indonesia Kreasi",
			Location:    "Indonesia",
			StartDate:   workStart,
			EndDate:     &workEnd,
			IsCurrent:   false,
			Description: "Worked as a Fullstack Developer building web applications, backend services, APIs, and chatbot automation using Golang. Responsible for designing and implementing systems based on business requirements, from architecture design to deployment. Developed customizable chatbot and automation systems integrated with WhatsApp (Whatsmeow), Telegram Bot, and Twilio. Also implemented reporting and data export services for internal teams and clients, tailored to their operational needs. Worked with Docker, Linux, and Git for environment management and version control, and used Postman for API testing and documentation to ensure reliable integrations.",
			SortOrder:   1,
			ImageURL:    "https://image-service-cdn.seek.com.au/ca1202d0d39bddb50cefc1b6b7a050c435669c2e/",
		},
		{
			Type:        "Work",
			Title:       "IT Internship",
			Org:         "Department of Agriculture and Plantation, Enrekang Regency",
			Location:    "Enrekang Regency, South Sulawesi, Indonesia",
			StartDate:   internStart,
			EndDate:     &internEnd,
			IsCurrent:   false,
			Description: "Participated in a university internship program at the Department of Agriculture and Plantation of Enrekang, assisting with distribution reporting for fertilizer distribution and recipients of agricultural materials and equipment. Contributed to developing a simple informational website and participated in technology outreach activities, introducing practical uses of technology for the department and local communities, particularly in the agriculture and plantation sectors.",
			SortOrder:   2,
			ImageURL:    "https://lh3.googleusercontent.com/gps-cs-s/AHVAwer1Kwym4neteWEvT0CoKF61AqkAY1Bb9wIctWsreBzON1j4gr0Z7Lr933heoZMis8WJUNxEEP8G0uOc4LzpAUy8EQyjFF1DGLrrQXseLkIH81mArM6uiQJ-Y4ZmdlnbEYD2GGA=s680-w680-h510-rw",
		},

		{
			Type:        "Education",
			Title:       "Bachelor of Informatics Engineering",
			Org:         "Universitas Kristen Indonesia Toraja",
			Location:    "Toraja, South Sulawesi, Indonesia",
			StartDate:   eduStart,
			EndDate:     &eduEnd,
			IsCurrent:   false,
			Description: "Graduated Cum Laude with a GPA of 3.89, earning a Bachelor’s degree in Informatics Engineering (S.Kom) from the Faculty of Engineering, Universitas Kristen Indonesia Toraja in 2022, and was recognized as the best graduate of the program.",
			SortOrder:   3,
			CertURL:     "https://drive.google.com/file/d/1A0hDQingopqunPV87WMY93O3JOWjg92I/view?usp=sharing",
			ImageURL:    "https://ukitoraja.ac.id/wp-content/uploads/2019/05/Logo-UKIT.png",
		},
	}
	db.Create(&experiences)
}
