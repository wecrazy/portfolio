package model

// TableNameConfig holds configurable database table names for all models.
type TableNameConfig struct {
	Admins          string
	Owners          string
	Projects        string
	Experiences     string
	Skills          string
	SocialLinks     string
	UploadedFiles   string
	OAuthUsers      string
	Comments        string
	ContactMessages string
	TechStacks      string
	Posts           string
	UpcomingItems   string
}

// tables holds the active table name configuration.
// Defaults are provided so the app works with zero config.
var tables = TableNameConfig{
	Admins:          "admins",
	Owners:          "owners",
	Projects:        "projects",
	Experiences:     "experiences",
	Skills:          "skills",
	SocialLinks:     "social_links",
	UploadedFiles:   "uploaded_files",
	OAuthUsers:      "oauth_users",
	Comments:        "comments",
	ContactMessages: "contact_messages",
	TechStacks:      "tech_stacks",
	Posts:           "posts",
	UpcomingItems:   "upcoming_items",
}

// SetTableNames overrides default table names with non-empty values from cfg.
func SetTableNames(cfg TableNameConfig) {
	if cfg.Admins != "" {
		tables.Admins = cfg.Admins
	}
	if cfg.Owners != "" {
		tables.Owners = cfg.Owners
	}
	if cfg.Projects != "" {
		tables.Projects = cfg.Projects
	}
	if cfg.Experiences != "" {
		tables.Experiences = cfg.Experiences
	}
	if cfg.Skills != "" {
		tables.Skills = cfg.Skills
	}
	if cfg.SocialLinks != "" {
		tables.SocialLinks = cfg.SocialLinks
	}
	if cfg.UploadedFiles != "" {
		tables.UploadedFiles = cfg.UploadedFiles
	}
	if cfg.OAuthUsers != "" {
		tables.OAuthUsers = cfg.OAuthUsers
	}
	if cfg.Comments != "" {
		tables.Comments = cfg.Comments
	}
	if cfg.ContactMessages != "" {
		tables.ContactMessages = cfg.ContactMessages
	}
	if cfg.TechStacks != "" {
		tables.TechStacks = cfg.TechStacks
	}
	if cfg.Posts != "" {
		tables.Posts = cfg.Posts
	}
	if cfg.UpcomingItems != "" {
		tables.UpcomingItems = cfg.UpcomingItems
	}
}
