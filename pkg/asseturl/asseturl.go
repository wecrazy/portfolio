package asseturl

import (
	"net/url"
	"strings"
)

// NormalizeExternalImageURL fixes a few common URL shapes so the application
// can treat them as direct image sources.
func NormalizeExternalImageURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return raw
	}

	host := strings.ToLower(u.Host)
	switch host {
	case "github.com", "www.github.com":
		u.Path = strings.Replace(u.Path, "/blob/", "/raw/", 1)
	case "raw.githubusercontent.com":
		u.Path = strings.Replace(u.Path, "/refs/heads/", "/", 1)
		u.Path = strings.Replace(u.Path, "/refs/tags/", "/", 1)
	}

	return u.String()
}
