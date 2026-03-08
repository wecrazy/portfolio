package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedSocialLinks creates demo social links. It is intentionally kept in the seed package because it depends on internal/model.
func seedSocialLinks(db *gorm.DB) {
	links := []model.SocialLink{
		{Platform: "github", URL: "https://github.com/wecrazy", IconClass: "bxl bx-github", Label: "GitHub", SortOrder: 1},
		// {Platform: "linkedin", URL: "https://id.linkedin.com/in/wegirandol-histara-littu-926219195", IconClass: "bxl bx-linkedin", Label: "LinkedIn", SortOrder: 2},
		{Platform: "jobstreet", URL: "https://id.jobstreet.com/id/profiles/wegirandol-histaralittu-SNWM16dvxC", IconClass: "bx bx-briefcase-alt-2", Label: "Jobstreet", SortOrder: 2},
		{Platform: "instagram", URL: "https://www.instagram.com/wecraz_y", IconClass: "bxl bx-instagram", Label: "Instagram", SortOrder: 3},
		{Platform: "facebook", URL: "https://www.facebook.com/wegil", IconClass: "bxl bx-facebook", Label: "Facebook", SortOrder: 4},
	}
	db.Create(&links)
}
