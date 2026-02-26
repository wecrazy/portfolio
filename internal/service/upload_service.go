package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"github.com/google/uuid"
)

// ProcessUpload validates and stores an uploaded file, returning the model (unsaved).
func ProcessUpload(fileHeader *multipart.FileHeader, category string) (*model.UploadedFile, error) {
	cfg := config.MyPortfolio.Get()

	// Determine size limit and allowed types.
	var maxSize int64
	var allowedTypes []string
	switch category {
	case "images":
		maxSize = cfg.Upload.MaxImageSize * 1024 * 1024
		allowedTypes = cfg.Upload.AllowedImageTypes
	case "resume":
		maxSize = cfg.Upload.MaxResumeSize * 1024 * 1024
		allowedTypes = cfg.Upload.AllowedResumeTypes
	default:
		maxSize = cfg.Upload.MaxImageSize * 1024 * 1024
		allowedTypes = cfg.Upload.AllowedImageTypes
	}

	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("file too large (max %d MB)", maxSize/(1024*1024))
	}

	// Read content type from header.
	contentType := fileHeader.Header.Get("Content-Type")
	allowed := false
	for _, t := range allowedTypes {
		if strings.EqualFold(contentType, t) {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("file type %s not allowed", contentType)
	}

	// Generate unique stored name.
	ext := filepath.Ext(fileHeader.Filename)
	storedName := uuid.New().String() + ext
	targetDir := filepath.Join(cfg.App.UploadDir, category)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	filePath := filepath.Join(targetDir, storedName)

	return &model.UploadedFile{
		OriginalName: fileHeader.Filename,
		StoredName:   storedName,
		FilePath:     filePath,
		MimeType:     contentType,
		FileSize:     fileHeader.Size,
		Category:     category,
	}, nil
}
