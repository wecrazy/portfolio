package handler

import (
	"log"

	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SubmitContact processes the public contact form: saves message and sends email.
func SubmitContact(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		subject := c.FormValue("subject")
		message := c.FormValue("message")

		if name == "" || email == "" || message == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Name, email, and message are required")
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

		// Return success partial or redirect.
		if c.Get("HX-Request") == "true" {
			return c.SendString(`<div class="alert alert-success">Thank you! Your message has been sent.</div>`)
		}
		return c.Redirect("/#contact")
	}
}
