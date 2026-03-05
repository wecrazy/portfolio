package seed

import (
	"my-portfolio/internal/model"

	"gorm.io/gorm"
)

// seedComments creates demo comments and replies. It is intentionally kept in the seed package because it depends on internal/model.
func seedComments(db *gorm.DB) {
	users := []model.OAuthUser{
		{Provider: "github", ProviderID: "demo-1", Email: "alice@example.com", DisplayName: "Alice Chen", AvatarURL: "https://i.pravatar.cc/150?u=alice"},
		{Provider: "google", ProviderID: "demo-2", Email: "bob@example.com", DisplayName: "Bob Smith", AvatarURL: "https://i.pravatar.cc/150?u=bob"},
		{Provider: "github", ProviderID: "demo-3", Email: "carlos@example.com", DisplayName: "Carlos Mendez", AvatarURL: "https://i.pravatar.cc/150?u=carlos"},
		{Provider: "google", ProviderID: "demo-4", Email: "diana@example.com", DisplayName: "Diana Park", AvatarURL: "https://i.pravatar.cc/150?u=diana"},
		{Provider: "github", ProviderID: "demo-5", Email: "evan@example.com", DisplayName: "Evan Torres", AvatarURL: "https://i.pravatar.cc/150?u=evan"},
		{Provider: "google", ProviderID: "demo-6", Email: "fiona@example.com", DisplayName: "Fiona Lim", AvatarURL: "https://i.pravatar.cc/150?u=fiona"},
		{Provider: "github", ProviderID: "demo-7", Email: "george@example.com", DisplayName: "George Nakamura", AvatarURL: "https://i.pravatar.cc/150?u=george"},
		{Provider: "google", ProviderID: "demo-8", Email: "hana@example.com", DisplayName: "Hana Yılmaz", AvatarURL: "https://i.pravatar.cc/150?u=hana"},
		{Provider: "system", ProviderID: "owner", Email: "john@example.com", DisplayName: "John Doe (Owner)"},
	}
	db.Create(&users)
	owner := users[8]

	// 25 top-level comments spread across 8 visitors — enough for 3 pages of 10
	// and to demonstrate the DOM-windowing "Showing a window" indicator.
	comments := []model.Comment{
		{OAuthUserID: users[0].ID, Body: "Great portfolio! I love the clean design and the tech stack section is really informative.", IsApproved: true},
		{OAuthUserID: users[1].ID, Body: "Impressive project list. The e-commerce platform looks really solid. Would love to see a demo!", IsApproved: true},
		{OAuthUserID: users[2].ID, Body: "The Go + Fiber combo is underrated — fast to build and blazing fast at runtime. Keep it up!", IsApproved: true},
		{OAuthUserID: users[3].ID, Body: "I really appreciate the attention to UI/UX here. The dark theme with the glass-morphism cards feels premium.", IsApproved: true},
		{OAuthUserID: users[4].ID, Body: "Your experience section showed up perfectly on mobile for me. How did you handle the timeline on small screens?", IsApproved: true},
		{OAuthUserID: users[5].ID, Body: "Love that you used HTMX instead of a heavy SPA framework. The page feels snappy even on a slow connection.", IsApproved: true},
		{OAuthUserID: users[6].ID, Body: "The blog with Markdown + embedded video support is a nice touch. Most dev portfolios skip that entirely.", IsApproved: true},
		{OAuthUserID: users[7].ID, Body: "Cum Laude with a 3.89 GPA — that's impressive! Did you specialize in web or systems programming during your degree?", IsApproved: true},
		{OAuthUserID: users[0].ID, Body: "Just read the HTMX blog post — the 'Show More' pagination pattern explanation was crystal clear. Bookmarked!", IsApproved: true},
		{OAuthUserID: users[1].ID, Body: "The Go Fiber post convinced me to finally try it for my next side project. Any tips for production deployment?", IsApproved: true},
		{OAuthUserID: users[2].ID, Body: "Your Docker + Nginx + SQLite stack is exactly what I've been looking for for small-scale personal projects.", IsApproved: true},
		{OAuthUserID: users[3].ID, Body: "Loving the i18n support! Switching between EN and ID in real time without a page reload — clean.", IsApproved: true},
		{OAuthUserID: users[4].ID, Body: "Is the backend entirely Go? I couldn't find a Node.js or Python dependency anywhere. Impressive!", IsApproved: true},
		{OAuthUserID: users[5].ID, Body: "The skills section with the animated progress bars is satisfying to watch. Small but impactful detail.", IsApproved: true},
		{OAuthUserID: users[6].ID, Body: "I noticed the PDF resume opens inside the site — that's a much better UX than triggering a browser download.", IsApproved: true},
		{OAuthUserID: users[7].ID, Body: "The 'What's Next' section is a great idea. Keeps visitors in the loop without needing a separate blog post.", IsApproved: true},
		{OAuthUserID: users[0].ID, Body: "How long did this portfolio take you to build end-to-end? The feature set is surprisingly complete.", IsApproved: true},
		{OAuthUserID: users[1].ID, Body: "The comment section with OAuth login (Google + GitHub) is a thoughtful touch. Encourages real interaction.", IsApproved: true},
		{OAuthUserID: users[2].ID, Body: "Real-time comment notifications via WebSocket — that's a level of detail most portfolios completely skip.", IsApproved: true},
		{OAuthUserID: users[3].ID, Body: "The profile image hover animation is a fun little Easter egg. Adds personality without being distracting.", IsApproved: true},
		{OAuthUserID: users[4].ID, Body: "Great use of AOS animations — they add polish without slowing anything down. Performance score still looks solid.", IsApproved: true},
		{OAuthUserID: users[5].ID, Body: "The contact form works great! I tested it and received a confirmation almost instantly. Very responsive.", IsApproved: true},
		{OAuthUserID: users[6].ID, Body: "The tech stack grid with icons from devicons is a nice visual touch. Way better than a plain text list.", IsApproved: true},
		{OAuthUserID: users[7].ID, Body: "Are you planning to open-source this portfolio template? I'd love to use it as a starting point!", IsApproved: true},
		{OAuthUserID: users[0].ID, Body: "The dark/light theme toggle respects the system preference by default — that's the UX best practice most devs skip.", IsApproved: true},
	}
	db.Create(&comments)

	// Owner replies to a handful of comments
	replies := []model.Comment{
		{OAuthUserID: owner.ID, ParentID: &comments[0].ID, Body: "Thank you Alice! Glad you like the design.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: owner.ID, ParentID: &comments[1].ID, Body: "Thanks Bob! I'll add a live demo link soon.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: owner.ID, ParentID: &comments[4].ID, Body: "Used CSS flexbox + a bit of media-query magic for the timeline. It collapses to a single column on mobile.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: owner.ID, ParentID: &comments[9].ID, Body: "For production I use a systemd service + Nginx reverse proxy. The Makefile has an --install flag that sets it all up.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: owner.ID, ParentID: &comments[16].ID, Body: "Roughly 3 weeks of evenings and weekends. Iterating fast with HTMX + Go templates helped a lot.", IsOwnerReply: true, IsApproved: true},
		{OAuthUserID: owner.ID, ParentID: &comments[23].ID, Body: "That's the plan! Will clean it up and publish it on GitHub once I've polished a few more rough edges.", IsOwnerReply: true, IsApproved: true},
	}
	db.Create(&replies)
}
