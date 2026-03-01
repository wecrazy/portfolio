package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"my-portfolio/internal/config"
	appI18n "my-portfolio/internal/i18n"
	"my-portfolio/internal/model"
	"my-portfolio/internal/service"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
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

// hxToast sends an HX-Trigger header that fires a contactToast event, which
// portfolio.html picks up and forwards to window.showToast().
func hxToast(c fiber.Ctx, msg, kind string) {
	type payload struct {
		Msg  string `json:"msg"`
		Type string `json:"type"`
	}
	type trigger struct {
		ContactToast payload `json:"contactToast"`
	}
	b, _ := json.Marshal(trigger{ContactToast: payload{Msg: msg, Type: kind}})
	c.Set("HX-Trigger", string(b))
}

// hcaptchaClient is a shared HTTP client with a reasonable timeout for
// calls to the hCaptcha siteverify API.
var hcaptchaClient = &http.Client{Timeout: 10 * time.Second}

// verifyHCaptcha calls the hCaptcha site-verify API and returns true when the
// token is valid. It is called inline from SubmitContact so the handler controls
// the full response — no middleware involved.
func verifyHCaptcha(secret, token, remoteIP string) bool {
	if secret == "" || token == "" {
		zap.L().Warn("hcaptcha: skipped — secret or token is empty")
		return false
	}
	vals := url.Values{
		"secret":   {secret},
		"response": {token},
	}
	if remoteIP != "" {
		vals.Set("remoteip", remoteIP)
	}
	resp, err := hcaptchaClient.PostForm("https://api.hcaptcha.com/siteverify", vals)
	if err != nil {
		zap.L().Error("hcaptcha: siteverify request failed", zap.Error(err))
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		zap.L().Error("hcaptcha: failed to parse response", zap.ByteString("body", body), zap.Error(err))
		return false
	}
	if !result.Success {
		zap.L().Warn("hcaptcha: verification failed", zap.Strings("error-codes", result.ErrorCodes))
	}
	return result.Success
}

// loc resolves a translation key for the current request language,
// falling back to the provided default when the key is missing.
func loc(c fiber.Ctx, key, fallback string) string {
	msg, _ := appI18n.T.Localize(c, key)
	if msg == "" {
		return fallback
	}
	return msg
}

// contactForm holds the parsed contact form values.
type contactForm struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// parseContactForm reads and trims all contact fields from the request.
func parseContactForm(c fiber.Ctx) contactForm {
	return contactForm{
		Name:    strings.TrimSpace(c.FormValue("name")),
		Email:   strings.TrimSpace(c.FormValue("email")),
		Subject: strings.TrimSpace(c.FormValue("subject")),
		Message: strings.TrimSpace(c.FormValue("message")),
	}
}

// validateContactForm checks required fields, length limits, and email format.
// Returns the i18n key + English fallback for the first violation, or empty strings when valid.
func validateContactForm(f contactForm) (key, fallback string) {
	if f.Name == "" || f.Email == "" || f.Message == "" {
		return "contact_required", "Name, email, and message are required."
	}
	if len(f.Name) > maxContactName {
		return "contact_name_too_long", "Name is too long (max 100 characters)."
	}
	if len(f.Email) > maxContactEmail {
		return "contact_email_too_long", "Email address is too long."
	}
	if len(f.Subject) > maxContactSubject {
		return "contact_subject_too_long", "Subject is too long (max 200 characters)."
	}
	if len(f.Message) > maxContactMessage {
		return "contact_message_too_long", "Message is too long (max 5000 characters)."
	}
	if !emailRE.MatchString(f.Email) {
		return "invalid_email", "Please enter a valid email address."
	}
	return "", ""
}

// SubmitContact processes the public contact form: saves message and sends email.
func SubmitContact(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Honeypot: hidden from real users, auto-filled by bots. Silently succeed.
		if c.FormValue("website") != "" {
			if c.Get("HX-Request") == "true" {
				hxToast(c, loc(c, "contact_sent", "Your message has been sent successfully."), "success")
				return c.Status(fiber.StatusOK).SendString("")
			}
			return c.Redirect().To("/#contact")
		}

		// hCaptcha: verify when enabled in config.
		cfg := config.MyPortfolio.Get()
		if cfg.HCaptcha.Enabled && !verifyHCaptcha(cfg.HCaptcha.Secret, c.FormValue("h-captcha-response"), c.IP()) {
			if c.Get("HX-Request") == "true" {
				hxToast(c, loc(c, "captcha_failed", "Captcha verification failed. Please try again."), "danger")
				return c.Status(fiber.StatusBadRequest).SendString("")
			}
			return c.Redirect().To("/#contact")
		}

		sendErr := func(key, fallback string) error {
			if c.Get("HX-Request") == "true" {
				hxToast(c, loc(c, key, fallback), "danger")
				return c.Status(fiber.StatusBadRequest).SendString("")
			}
			return c.Redirect().To("/#contact")
		}

		f := parseContactForm(c)
		if key, fb := validateContactForm(f); key != "" {
			return sendErr(key, fb)
		}

		db.Create(&model.ContactMessage{Name: f.Name, Email: f.Email, Subject: f.Subject, Message: f.Message})

		// Fire-and-forget: message is already saved; don't block the response on SMTP.
		go func() {
			if err := service.SendContactEmail(f.Name, f.Email, f.Subject, f.Message); err != nil {
				zap.L().Error("failed to send contact email",
					zap.String("to", f.Email),
					zap.String("subject", f.Subject),
					zap.Error(err),
				)
			}
		}()

		if c.Get("HX-Request") == "true" {
			hxToast(c, loc(c, "contact_sent", "Your message has been sent successfully."), "success")
			return c.SendString("") // empty body — clears #contact-result
		}
		return c.Redirect().To("/#contact")
	}
}
