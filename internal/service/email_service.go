// Package service provides business logic for email, OAuth, and file uploads.
package service

import (
	"crypto/tls"
	"fmt"
	"html"

	"my-portfolio/internal/config"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

// newHermes builds a branded hermes instance from live config.
func newHermes() hermes.Hermes {
	cfg := config.MyPortfolio.Get()
	return hermes.Hermes{
		Theme: new(hermes.Default),
		Product: hermes.Product{
			Name:      cfg.Owner.Name,
			Link:      cfg.App.BaseURL,
			Copyright: fmt.Sprintf("© %s — All rights reserved.", cfg.Owner.Name),
			TroubleText: "If the button above is not working, copy and paste the" +
				" URL below into your web browser.",
		},
	}
}

// SendContactEmail sends a contact form message to the portfolio owner via SMTP.
func SendContactEmail(name, email, subject, message string) error {
	cfg := config.MyPortfolio.Get()
	smtp := cfg.SMTP

	if subject == "" {
		subject = "(no subject)"
	}

	h := newHermes()

	esc := html.EscapeString

	hermesEmail := hermes.Email{
		Body: hermes.Body{
			Name: cfg.Owner.Name,
			Intros: []string{
				"You have received a new contact message via your portfolio.",
			},
			Dictionary: []hermes.Entry{
				{Key: "From", Value: esc(name)},
				{Key: "Email", Value: email},
				{Key: "Subject", Value: esc(subject)},
				{Key: "Message", Value: esc(message)},
			},
			Outros: []string{
				"To reply, simply respond to this email — the Reply-To address is set to the sender.",
				"This message was sent automatically from your portfolio contact form.",
			},
		},
	}

	htmlBody, err := h.GenerateHTML(hermesEmail)
	if err != nil {
		return fmt.Errorf("hermes generate HTML: %w", err)
	}
	textBody, err := h.GeneratePlainText(hermesEmail)
	if err != nil {
		return fmt.Errorf("hermes generate text: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtp.From)
	m.SetHeader("To", smtp.To)
	m.SetHeader("Reply-To", email)
	m.SetHeader("Subject", fmt.Sprintf("[Portfolio Contact] %s", subject))
	m.SetHeader("Auto-Submitted", "auto-generated")
	m.SetHeader("X-Auto-Response-Suppress", "OOF, AutoReply")
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	// Ensure proper TLS hostname verification for providers that require SNI.
	d.TLSConfig = &tls.Config{ServerName: smtp.Host}

	dialErr := d.DialAndSend(m)
	if dialErr != nil {
		// Surface SMTP host and user in the error for easier diagnosis (no password).
		return fmt.Errorf("smtp send to %s as %s: %w", smtp.Host, smtp.Username, dialErr)
	}

	return nil
}
