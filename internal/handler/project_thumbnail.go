package handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/asseturl"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
)

const (
	projectThumbCacheDir    = "data/thumbcache/projects"
	projectThumbCacheTTL    = 24 * time.Hour
	projectThumbCacheHeader = "public, max-age=86400, stale-while-revalidate=604800"
	projectThumbMaxBytes    = 20 << 20
	projectThumbMaxWidth    = 1280
	projectThumbMaxHeight   = 880
)

var projectThumbHTTPClient = &http.Client{Timeout: 15 * time.Second}

// ProjectThumbnail proxies and caches remote project thumbnails so the browser
// does not need to hit slow or unstable third-party image hosts directly.
func ProjectThumbnail(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var project model.Project
		if err := db.Select("id", "thumbnail_url").First(&project, c.Params("id")).Error; err != nil {
			return serveDefaultProjectThumbnail(c)
		}

		src := asseturl.NormalizeExternalImageURL(project.ThumbnailURL)
		if !isAllowedRemoteImageURL(src) {
			return serveDefaultProjectThumbnail(c)
		}

		if cachePath, ok := findFreshProjectThumbnailCache(src); ok {
			return sendProjectThumbnailFile(c, cachePath)
		}

		data, contentType, err := fetchProjectThumbnail(src)
		if err != nil {
			return serveDefaultProjectThumbnail(c)
		}

		optimized, optimizedType, ext, err := optimizeProjectThumbnail(data, contentType)
		if err != nil {
			return serveDefaultProjectThumbnail(c)
		}

		if cachePath, err := storeProjectThumbnailCache(src, optimized, ext); err == nil {
			return sendProjectThumbnailFile(c, cachePath)
		}

		c.Set("Cache-Control", projectThumbCacheHeader)
		c.Set("Content-Type", optimizedType)
		return c.Send(optimized)
	}
}

func serveDefaultProjectThumbnail(c fiber.Ctx) error {
	c.Set("Cache-Control", projectThumbCacheHeader)
	return c.SendFile("./web/static/img/default-project.jpg")
}

func sendProjectThumbnailFile(c fiber.Ctx, cachePath string) error {
	c.Set("Cache-Control", projectThumbCacheHeader)
	return c.SendFile(cachePath)
}

func isAllowedRemoteImageURL(raw string) bool {
	if raw == "" {
		return false
	}

	u, err := url.Parse(raw)
	if err != nil {
		return false
	}

	scheme := strings.ToLower(u.Scheme)
	if (scheme != "http" && scheme != "https") || u.Host == "" {
		return false
	}

	host := strings.ToLower(u.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return false
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return false
		}
	}

	return true
}

func fetchProjectThumbnail(raw string) ([]byte, string, error) {
	req, err := http.NewRequest(http.MethodGet, raw, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "my-portfolio-thumbnail-proxy/1.0")

	resp, err := projectThumbHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("thumbnail fetch failed: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, projectThumbMaxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(data) > projectThumbMaxBytes {
		return nil, "", fmt.Errorf("thumbnail exceeds %d bytes", projectThumbMaxBytes)
	}

	contentType := normalizeContentType(resp.Header.Get("Content-Type"), data)
	if !strings.HasPrefix(contentType, "image/") && !isSVGContent(data) {
		return nil, "", fmt.Errorf("unsupported thumbnail type: %s", contentType)
	}

	return data, contentType, nil
}

func optimizeProjectThumbnail(data []byte, hintedType string) ([]byte, string, string, error) {
	contentType := normalizeContentType(hintedType, data)
	if isSVGContent(data) || contentType == "image/svg+xml" {
		return data, "image/svg+xml", ".svg", nil
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		ext := extensionForContentType(contentType)
		if ext == "" {
			return nil, "", "", err
		}
		return data, contentType, ext, nil
	}

	resized := resizeProjectThumbnail(img)
	format = strings.ToLower(format)

	var buf bytes.Buffer
	var optimizedType string
	var ext string

	switch format {
	case "jpeg", "jpg", "webp":
		if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 82}); err != nil {
			return nil, "", "", err
		}
		optimizedType = "image/jpeg"
		ext = ".jpg"
	case "png", "gif":
		if err := png.Encode(&buf, resized); err != nil {
			return nil, "", "", err
		}
		optimizedType = "image/png"
		ext = ".png"
	default:
		ext = extensionForContentType(contentType)
		if ext == "" {
			return nil, "", "", fmt.Errorf("unsupported decoded format: %s", format)
		}
		return data, contentType, ext, nil
	}

	optimized := buf.Bytes()
	if len(optimized) >= len(data) {
		ext = extensionForContentType(contentType)
		if ext == "" {
			ext = ".img"
		}
		return data, contentType, ext, nil
	}

	return optimized, optimizedType, ext, nil
}

func resizeProjectThumbnail(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return src
	}
	if width <= projectThumbMaxWidth && height <= projectThumbMaxHeight {
		return src
	}

	widthRatio := float64(projectThumbMaxWidth) / float64(width)
	heightRatio := float64(projectThumbMaxHeight) / float64(height)
	scale := widthRatio
	if heightRatio < scale {
		scale = heightRatio
	}

	dstWidth := int(float64(width) * scale)
	dstHeight := int(float64(height) * scale)
	if dstWidth < 1 {
		dstWidth = 1
	}
	if dstHeight < 1 {
		dstHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

func normalizeContentType(headerValue string, data []byte) string {
	if headerValue != "" {
		if mediaType, _, err := mime.ParseMediaType(headerValue); err == nil && mediaType != "" {
			return strings.ToLower(mediaType)
		}
	}
	return strings.ToLower(http.DetectContentType(data))
}

func isSVGContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	sample := strings.TrimSpace(string(data[:min(len(data), 512)]))
	return strings.HasPrefix(sample, "<svg") || (strings.HasPrefix(sample, "<?xml") && strings.Contains(sample, "<svg"))
}

func extensionForContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ""
	}
}

func findFreshProjectThumbnailCache(raw string) (string, bool) {
	hash := projectThumbnailHash(raw)
	pattern := filepath.Join(projectThumbCacheDir, hash+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", false
	}

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if time.Since(info.ModTime()) <= projectThumbCacheTTL {
			return match, true
		}
	}

	return "", false
}

func storeProjectThumbnailCache(raw string, data []byte, ext string) (string, error) {
	if err := os.MkdirAll(projectThumbCacheDir, 0o755); err != nil {
		return "", err
	}

	hash := projectThumbnailHash(raw)
	pattern := filepath.Join(projectThumbCacheDir, hash+".*")
	if matches, err := filepath.Glob(pattern); err == nil {
		for _, match := range matches {
			_ = os.Remove(match)
		}
	}

	cachePath := filepath.Join(projectThumbCacheDir, hash+ext)
	if err := os.WriteFile(cachePath, data, 0o644); err != nil {
		return "", err
	}

	return cachePath, nil
}

func projectThumbnailHash(raw string) string {
	sum := sha1.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
