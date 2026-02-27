package handler

import (
	"html"
	"log"
	"regexp"
	"strings"

	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Field length limits for the contact form.
const (
	maxContactName    = 100
	maxContactEmail   = 254
	maxContactSubject = 200
	maxContactMessage = 5000
)

// Simple RFC-5321-compliant-ish email regex.
var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// SubmitContact processes the public contact form: saves message and sends email.
func SubmitContact(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		name := strings.TrimSpace(c.FormValue("name"))
		email := strings.TrimSpace(c.FormValue("email"))
		subject := strings.TrimSpace(c.FormValue("subject"))
		message := strings.TrimSpace(c.FormValue("message"))

		// Required fields
		if name == "" || email == "" || message == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Name, email, and message are required")
		}

		// Length limits
		if len(name) > maxContactName {
			return c.Status(fiber.StatusBadRequest).SendString("Name is too long (max 100 characters)")
		}
		if len(email) > maxContactEmail {
			return c.Status(fiber.StatusBadRequest).SendString("Email address is too long")
		}
		if len(subject) > maxContactSubject {
			return c.Status(fiber.StatusBadRequest).SendString("Subject is too long (max 200 characters)")
		}
		if len(message) > maxContactMessage {
			return c.Status(fiber.StatusBadRequest).SendString("Message is too long (max 5000 characters)")
		}

		// Email format validation
		if !emailRE.MatchString(email) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid email address")
		}

		// Save to DB regardless of email delivery.
		contact := model.ContactMessage{
			Name:    name,
			Email:   email,
			Subject: subject,
			Message: message,
		}
		db.Create(&contact)

		// Attempt to send email.
		if err := service.SendContactEmail(name, email, subject, message); err != nil {
			log.Printf("Failed to send contact email: %v", err)
		}

		// Return success partial for HTMX or redirect for plain form.
		if c.Get("HX-Request") == "true" {
			msg, _ := appI18n.T.Localize(c, "contact_sent")
			if msg == "" {
				msg = "Your message has been sent successfully."
			}
			// Escape the i18n string to prevent XSS even if a key is misconfigured.
			return c.SendString(`<div class="alert alert-success">` + html.EscapeString(msg) + `</div>`)
		}
		return c.Redirect().To("/#contact")
	}
}
