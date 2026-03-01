package middleware

import "github.com/gofiber/fiber/v3"

// Security adds hardened security headers to all responses.
//
// Headers applied:
//   - X-Content-Type-Options    — prevents MIME-sniffing
//   - X-Frame-Options           — blocks clickjacking (legacy browsers)
//   - X-XSS-Protection          — XSS auditor hint (legacy browsers)
//   - Referrer-Policy           — controls Referer header leakage
//   - Content-Security-Policy   — restricts resource origins
//   - Strict-Transport-Security — enforce HTTPS (HSTS)
//   - Permissions-Policy        — restricts browser feature access
//   - X-DNS-Prefetch-Control    — disables speculative DNS lookups
func Security() fiber.Handler {
	// Content-Security-Policy: allow known CDNs used by the portfolio and
	// admin UI (Bootstrap, Boxicons, HTMX, Alpine.js, highlight.js, hCaptcha).
	// 'unsafe-inline' is required by HTMX hx-* event handlers and inline styles.
	const csp = "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline' " +
		"cdn.jsdelivr.net unpkg.com newassets.hcaptcha.com js.hcaptcha.com; " +
		"style-src 'self' 'unsafe-inline' " +
		"cdn.jsdelivr.net fonts.googleapis.com cdn.boxicons.com unpkg.com; " +
		"font-src 'self' fonts.gstatic.com cdn.boxicons.com data:; " +
		"img-src 'self' data: blob: https:; " +
		"connect-src 'self' ws: wss: cdn.jsdelivr.net hcaptcha.com newassets.hcaptcha.com; " +
		"frame-src 'self' newassets.hcaptcha.com hcaptcha.com www.youtube.com youtube.com player.vimeo.com; " +
		"frame-ancestors 'self'; " +
		"object-src 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self';"

	return func(c fiber.Ctx) error {
		// Legacy hardening
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Modern hardening
		c.Set("Content-Security-Policy", csp)
		// max-age=1 year; includeSubDomains. preload omitted intentionally so
		// the site is not locked into the HSTS preload list prematurely.
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

		return c.Next()
	}
}
