// Package service provides business logic for email, OAuth, and file uploads.
package service

import (
	"fmt"

	"my-portfolio/internal/config"

	"gopkg.in/gomail.v2"
)

// SendContactEmail sends a contact form message to the portfolio owner via SMTP.
func SendContactEmail(name, email, subject, message string) error {
	cfg := config.MyPortfolio.Get().SMTP

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", cfg.To)
	m.SetHeader("Reply-To", email)
	m.SetHeader("Subject", fmt.Sprintf("[Portfolio Contact] %s", subject))
	m.SetBody("text/html", buildContactEmailHTML(name, email, subject, message))

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}

func buildContactEmailHTML(name, email, subject, body string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
    <h2 style="color: #667eea;">New Contact Message</h2>
    <table style="width: 100%%; border-collapse: collapse;">
        <tr><td style="padding: 8px; font-weight: bold;">From:</td><td style="padding: 8px;">%s</td></tr>
        <tr><td style="padding: 8px; font-weight: bold;">Email:</td><td style="padding: 8px;">%s</td></tr>
        <tr><td style="padding: 8px; font-weight: bold;">Subject:</td><td style="padding: 8px;">%s</td></tr>
    </table>
    <hr style="border: 1px solid #eee; margin: 20px 0;">
    <div style="padding: 10px; background: #f9f9f9; border-radius: 8px;">
        <p>%s</p>
    </div>
    <p style="color: #999; font-size: 12px; margin-top: 30px;">Sent from your portfolio contact form</p>
</body>
</html>`, name, email, subject, body)
}
