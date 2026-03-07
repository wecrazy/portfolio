package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"my-portfolio/internal/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// simple in-memory cache for fetched certificates.  Entries expire after
// a short period (2h) to avoid holding large files forever.  This is not
// distributed and resets on process restart, which is acceptable for a demo.
var certCache = struct {
	sync.Mutex
	m map[string]cacheEntry
}{m: make(map[string]cacheEntry)}

// cacheEntry holds the downloaded PDF bytes and its expiration time.
type cacheEntry struct {
	data []byte
	exp  time.Time
}

// StartCacheCleaner launches a goroutine to remove expired cache entries.
func StartCacheCleaner() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			certCache.Lock()
			now := time.Now()
			for k, ent := range certCache.m {
				if now.After(ent.exp) {
					delete(certCache.m, k)
				}
			}
			certCache.Unlock()
		}
	}()
}

// resolveSource determines the URL to fetch from either a raw "url" query
// parameter or a numeric "id" that maps to a certificate record.
func resolveSource(db *gorm.DB, c fiber.Ctx) (string, error) {
	raw := c.Query("url")
	if idstr := c.Query("id"); idstr != "" {
		var cert model.Certificate
		if err := db.Preload("File").First(&cert, idstr).Error; err != nil {
			return "", fiber.ErrNotFound
		}
		if cert.File != nil && cert.File.StoredName != "" {
			raw = "/uploads/" + cert.File.StoredName
		} else if cert.CertURL != "" {
			raw = cert.CertURL
		} else {
			return "", fiber.ErrNotFound
		}
	}
	if raw == "" {
		return "", fiber.ErrBadRequest
	}
	if strings.Contains(raw, "drive.google.com") {
		raw = convertDriveURL(raw)
	}
	return raw, nil
}

func convertDriveURL(raw string) string {
	parts := strings.Split(raw, "/")
	for i, p := range parts {
		if p == "d" && i+1 < len(parts) {
			id := parts[i+1]
			return "https://drive.google.com/uc?export=download&id=" + id
		}
	}
	return raw
}

func fetchAndCache(raw string) ([]byte, error) {
	certCache.Lock()
	if ent, ok := certCache.m[raw]; ok && time.Now().Before(ent.exp) {
		certCache.Unlock()
		return ent.data, nil
	}
	certCache.Unlock()

	resp, err := http.Get(raw)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch failed: %d", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	certCache.Lock()
	certCache.m[raw] = cacheEntry{data: buf, exp: time.Now().Add(2 * time.Hour)}
	certCache.Unlock()
	return buf, nil
}

// maybeFileCache looks for a cached file on disk and returns its bytes if still
// fresh (younger than two hours).
func maybeFileCache(raw string, c fiber.Ctx) ([]byte, bool) {
	cacheDir := "data/certcache"
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, false
	}
	fname := ""
	if idstr := c.Query("id"); idstr != "" {
		fname = idstr + ".pdf"
	} else {
		fname = fmt.Sprintf("%x.pdf", md5.Sum([]byte(raw)))
	}
	path := filepath.Join(cacheDir, fname)
	info, err := os.Stat(path)
	if err == nil && time.Since(info.ModTime()) < 2*time.Hour {
		data, err := ioutil.ReadFile(path)
		if err == nil {
			return data, true
		}
	}
	return nil, false
}

// CertPreview returns a handler that proxies the resolved source URL and
// ensures an Access-Control-Allow-Origin header for embedding.  Complexity is
// low because the heavy lifting is delegated to helpers above.
func CertPreview(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		raw, err := resolveSource(db, c)
		if err != nil {
			return err
		}

		if data, ok := maybeFileCache(raw, c); ok {
			c.Set("Access-Control-Allow-Origin", "*")
			c.Type("pdf")
			return c.Send(data)
		}

		buf, err := fetchAndCache(raw)
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, err.Error())
		}
		c.Set("Access-Control-Allow-Origin", "*")
		c.Type("pdf")
		return c.Send(buf)
	}
}
