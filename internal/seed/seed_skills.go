package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedSkills creates demo skills. It is intentionally kept in the seed package because it depends on internal/model.
func seedSkills(db *gorm.DB) {
	skills := []model.Skill{
		{Name: "Go", Category: "Languages", IconClass: "devicon-go-original-wordmark", IconURL: deviconCDN + "/go/go-original-wordmark.svg", Proficiency: 85, SortOrder: 1},
		{Name: "Python", Category: "Languages", IconClass: "devicon-python-plain", IconURL: deviconCDN + "/python/python-original.svg", Proficiency: 75, SortOrder: 2},
		{Name: "JavaScript", Category: "Languages", IconClass: "devicon-javascript-plain", IconURL: deviconCDN + "/javascript/javascript-plain.svg", Proficiency: 76, SortOrder: 3},
		{Name: "PHP", Category: "Languages", IconClass: "devicon-php-plain", IconURL: deviconCDN + "/php/php-plain.svg", Proficiency: 80, SortOrder: 4},
		{Name: "Rust", Category: "Languages", IconClass: "devicon-rust-original", IconURL: deviconCDN + "/rust/rust-original.svg", Proficiency: 5, SortOrder: 5},
		{Name: "C++", Category: "Languages", IconClass: "devicon-cplusplus-plain", IconURL: deviconCDN + "/cplusplus/cplusplus-original.svg", Proficiency: 50, SortOrder: 6},
		{Name: "Java", Category: "Languages", IconClass: "devicon-java-plain", IconURL: deviconCDN + "/java/java-original.svg", Proficiency: 30, SortOrder: 7},
		{Name: "HTMX", Category: "Frontend", IconClass: "bxf bx-bolt-circle", IconURL: "https://cdn.jsdelivr.net/gh/bigskysoftware/htmx@v2.0.4/www/static/img/htmx_logo.1.png", Proficiency: 90, SortOrder: 1},
		{Name: "Fiber", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg", Proficiency: 90, SortOrder: 1},
		{Name: "Gin", Category: "Backend", IconClass: "bxf bx-bolt-circle", IconURL: "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png", Proficiency: 90, SortOrder: 2},
		{Name: "Code Igniter", Category: "Backend", IconClass: "devicon-codeigniter-plain", IconURL: deviconCDN + "/codeigniter/codeigniter-plain.svg", Proficiency: 80, SortOrder: 3},
		{Name: "Docker", Category: "DevOps", IconClass: "devicon-docker-plain", IconURL: deviconCDN + "/docker/docker-original.svg", Proficiency: 76, SortOrder: 1},
		{Name: "Linux", Category: "DevOps", IconClass: "devicon-linux-plain", IconURL: deviconCDN + "/linux/linux-original.svg", Proficiency: 80, SortOrder: 2},
		{Name: "Git", Category: "DevOps", IconClass: "devicon-git-plain", IconURL: deviconCDN + "/git/git-original.svg", Proficiency: 85, SortOrder: 3},
		{Name: "Adobe Photoshop", Category: "Design", IconClass: "devicon-photoshop-plain", IconURL: deviconCDN + "/photoshop/photoshop-plain.svg", Proficiency: 50, SortOrder: 1},
		{Name: "Adobe After Effects", Category: "Design", IconClass: "devicon-aftereffects-plain", IconURL: deviconCDN + "/aftereffects/aftereffects-plain.svg", Proficiency: 30, SortOrder: 2},
		{Name: "Adobe Premiere Pro", Category: "Design", IconClass: "devicon-premierepro-plain", IconURL: deviconCDN + "/premierepro/premierepro-plain.svg", Proficiency: 25, SortOrder: 3},
	}
	db.Create(&skills)
}
