package seed

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"
	"my-portfolio/pkg/fileutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// linkImageRecord creates an UploadedFile DB record for an image already on disk.
// It is intentionally kept in the seed package because it depends on internal/model.
func linkImageRecord(db *gorm.DB, storedName, filePath string) *model.UploadedFile {
	ext := strings.ToLower(filepath.Ext(storedName))
	mimeType := fileutil.MimeByExt(ext)
	if mimeType == "application/octet-stream" {
		mimeType = "image/jpeg"
	}
	var size int64
	if info, err := os.Stat(filePath); err == nil {
		size = info.Size()
	}
	rec := &model.UploadedFile{
		OriginalName: storedName,
		StoredName:   storedName,
		FilePath:     filePath,
		MimeType:     mimeType,
		FileSize:     size,
		Category:     "images",
	}
	if err := db.Create(rec).Error; err != nil {
		log.Printf("Warning: failed to create image DB record: %v", err)
		return nil
	}
	return rec
}

// relinkUploadImage scans uploadDir for the newest allowed image and creates a
// DB record for it. Returns nil when no suitable file is found.
func relinkUploadImage(db *gorm.DB, uploadDir string, allowedExts map[string]bool) *model.UploadedFile {
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return nil
	}

	type candidate struct {
		name    string
		modTime time.Time
	}
	var candidates []candidate
	for _, e := range entries {
		if e.IsDir() || !allowedExts[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		if info, err := e.Info(); err == nil {
			candidates = append(candidates, candidate{name: e.Name(), modTime: info.ModTime()})
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].modTime.After(candidates[j].modTime)
	})
	chosen := candidates[0]
	rec := linkImageRecord(db, chosen.name, filepath.Join(uploadDir, chosen.name))
	if rec != nil {
		log.Printf("Re-linked existing profile image from uploads: %s", chosen.name)
	}
	return rec
}

// copyStaticImage copies the configured static profile image into uploadDir and
// creates a DB record for it. Returns nil when the source is absent or invalid.
func copyStaticImage(db *gorm.DB, cfg config.TypeMyPortfolio, uploadDir string, allowedExts map[string]bool) *model.UploadedFile {
	if cfg.Owner.ProfileImage == "" {
		return nil
	}
	srcPath := filepath.Join(cfg.App.StaticDir, strings.TrimPrefix(cfg.Owner.ProfileImage, "/"))
	if !fileutil.Exists(srcPath) {
		return nil
	}
	ext := strings.ToLower(filepath.Ext(srcPath))
	if !allowedExts[ext] {
		return nil
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil
	}
	storedName := uuid.New().String() + ext
	dstPath := filepath.Join(uploadDir, storedName)
	if err := fileutil.CopyFile(srcPath, dstPath); err != nil {
		log.Printf("Warning: could not copy static profile image: %v", err)
		return nil
	}
	rec := linkImageRecord(db, storedName, dstPath)
	if rec != nil {
		log.Printf("Copied static profile image to uploads: %s", storedName)
	}
	return rec
}

// seedOwner creates the default owner profile if it doesn't exist. Safe to call on every startup.
func seedOwner(db *gorm.DB, cfg config.TypeMyPortfolio) {
	var count int64
	db.Model(&model.Owner{}).Count(&count)
	if count > 0 {
		return
	}

	allowedExts := fileutil.AllowedExts(cfg.Upload.AllowedImageTypes)
	uploadDir := filepath.Join(cfg.App.UploadDir, "images")

	// Priority 1: re-link the newest image already in uploads/images/ (survives db-reset).
	imgProfile := relinkUploadImage(db, uploadDir, allowedExts)

	// Priority 2: fall back to the static file declared in config, copy it into uploads/images/.
	if imgProfile == nil {
		imgProfile = copyStaticImage(db, cfg, uploadDir, allowedExts)
	}

	// Resume: pick up any PDF already sitting in uploads/resume/ (survives db-reset).
	var resumeFile *model.UploadedFile
	resumeDir := filepath.Join(cfg.App.UploadDir, "resume")
	if entries, err := os.ReadDir(resumeDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.ToLower(filepath.Ext(name)) != ".pdf" {
				continue
			}
			fullPath := filepath.Join(resumeDir, name)
			var size int64
			if info, err2 := e.Info(); err2 == nil {
				size = info.Size()
			}
			rec := &model.UploadedFile{
				OriginalName: name,
				StoredName:   name,
				FilePath:     fullPath,
				MimeType:     "application/pdf",
				FileSize:     size,
				Category:     "resume",
			}
			if err2 := db.Create(rec).Error; err2 == nil {
				resumeFile = rec
			}
			break // use first PDF found
		}
	}

	owner := model.Owner{
		FullName:     cfg.Owner.Name,
		Title:        cfg.Owner.Title,
		Tagline:      cfg.Owner.Tagline,
		Bio:          cfg.Owner.Bio,
		ProfileImage: imgProfile,
		ResumeFile:   resumeFile,
		Email:        cfg.Owner.Email,
		Phone:        cfg.Owner.Phone,
		Location:     cfg.Owner.Location,
	}
	if err := db.Create(&owner).Error; err != nil {
		log.Fatalf("Failed to seed owner profile: %v", err)
	}
	log.Println("Seeded default owner profile")
}
