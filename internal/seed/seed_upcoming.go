package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedUpcomingItems creates demo upcoming projects and announcements. It is intentionally kept in the seed package because it depends on internal/model.
func seedUpcomingItems(db *gorm.DB) {
	items := []model.UpcomingItem{
		{
			Title:       "Fleet Tracking System",
			Description: "A real-time fleet tracking system using GPS, MQTT and Go. It provides live location updates, route history, and geofencing alerts for efficient fleet management.  Its workflow image: /uploads/images/fleet-tracking.jpeg",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-phone",
			SortOrder:   1,
			IsVisible:   true,
		},
		{
			Title:       "API for Toraja Dictionary",
			Description: "Creating a RESTful API for the Toraja dictionary to provide easy access to Toraja language resources. The API will support search, retrieval, and management of dictionary entries.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-book",
			SortOrder:   2,
			IsVisible:   true,
		},
		{
			Title:       "Self-Hosted Conversation Dashboard",
			Description: "Developing a self-hosted conversation dashboard integrated with a chat bot. The dashboard will allow users to manage and analyze conversations, while the chat bot will provide automated responses and insights.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-discussion",
			SortOrder:   3,
			IsVisible:   true,
		},
		{
			Title:       "AI Agent on Low Spec Devices",
			Description: "Developing an AI agent that can run on low-spec PCs or mobile devices using PicoClaw (https://github.com/sipeed/picoclaw). The agent will provide basic AI functionalities while optimizing performance for resource-constrained environments.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-robot",
			SortOrder:   4,
			IsVisible:   true,
		},
		{
			Title:       "Self Guided Product Demo",
			Description: "Building a fully automated product demonstration that runs without human intervention. The demo uses a pre-recorded script and video, allowing users to experience and interact with product features as if in a live session, but everything is handled automatically.",
			Type:        "project",
			Status:      "planned",
			IconClass:   "bxf bx-store",
			SortOrder:   5,
			IsVisible:   true,
		},
	}
	db.Create(&items)
}
