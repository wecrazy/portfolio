// Package seed provides initial data seeding for the database.
package seed

import (
	"log"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// deviconCDN is the base URL for fetching technology icons in seeds. It is used
// by both skills and tech stacks to ensure consistent icons. It is intentionally
// kept in the seed package because it is only relevant for seeding demo data and
// not used elsewhere in the app.
const deviconCDN = "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons"

// SeedIfNeeded creates the default admin user and an empty owner profile when
// they don't already exist. Safe to call on every startup.
func SeedIfNeeded(db *gorm.DB, cfg config.TypeMyPortfolio) {
	seedAdmin(db, cfg)
	seedOwner(db, cfg)
	seedDemoData(db)
}

// seedDemoData creates demo projects, experience, skills, etc if they don't already exist. Safe to call on every startup.
func seedDemoData(db *gorm.DB) {
	// Only seed projects/experience/skills/etc if projects table is empty.
	var count int64
	db.Model(&model.Project{}).Count(&count)
	if count == 0 {
		seedProjects(db)
		seedExperiences(db)
		seedSkills(db)
		seedSocialLinks(db)
		seedTechStacks(db)
		seedCertificates(db)
		// seedComments(db)
		log.Println("Seeded demo data")
	}

	// Seed demo media files (image, video, audio) used by the demo blog posts.
	seedDemoMediaFiles(db)

	// Seed posts independently so they work even on existing installs.
	var postCount int64
	if err := db.Model(&model.Post{}).Count(&postCount).Error; err != nil {
		log.Printf("ERROR counting posts: %v", err)
	}

	if postCount == 0 {
		if err := seedPosts(db); err != nil {
			log.Printf("ERROR seeding posts: %v", err)
		} else {
			log.Println("Seeded demo blog posts")
		}
	}

	// Seed upcoming items independently.
	var upcomingCount int64
	db.Model(&model.UpcomingItem{}).Count(&upcomingCount)
	if upcomingCount == 0 {
		seedUpcomingItems(db)
		log.Println("Seeded demo upcoming items")
	}
}
