// Package main is the entry point for the portfolio server.
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"my-portfolio/internal/config"
	"my-portfolio/internal/database"
	"my-portfolio/internal/middleware"
	"my-portfolio/internal/model"
	"my-portfolio/internal/router"
	"my-portfolio/internal/seed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// 1. Load config (auto-detects ENV=dev|prod, hot-reloads on change).
	config.MyPortfolio.MustInit("my-portfolio")
	go config.MyPortfolio.Watch()

	cfg := config.MyPortfolio.Get()

	// 2. Wire table names from config into models.
	model.SetTableNames(model.TableNameConfig{
		Admins:          cfg.TableNames.Admins,
		Owners:          cfg.TableNames.Owners,
		Projects:        cfg.TableNames.Projects,
		Experiences:     cfg.TableNames.Experiences,
		Skills:          cfg.TableNames.Skills,
		SocialLinks:     cfg.TableNames.SocialLinks,
		UploadedFiles:   cfg.TableNames.UploadedFiles,
		OAuthUsers:      cfg.TableNames.OAuthUsers,
		Comments:        cfg.TableNames.Comments,
		ContactMessages: cfg.TableNames.ContactMessages,
		TechStacks:      cfg.TableNames.TechStacks,
		Posts:           cfg.TableNames.Posts,
		UpcomingItems:   cfg.TableNames.UpcomingItems,
	})

	// 3. Init database.
	db := database.InitSQLite(cfg)
	database.AutoMigrate(db)

	// 4. Seed defaults.
	seed.SeedIfNeeded(db, cfg)

	// 5. Template engine.
	engine := html.New("./web/templates", ".html")
	engine.Reload(cfg.App.Debug)
	engine.AddFunc("upper", strings.ToUpper)
	engine.AddFunc("lower", strings.ToLower)
	engine.AddFunc("contains", strings.Contains)
	engine.AddFunc("split", strings.Split)
	engine.AddFunc("join", strings.Join)
	engine.AddFunc("safeHTML", func(s string) template.HTML {
		return template.HTML(s)
	})
	engine.AddFunc("formatDate", func(t time.Time) string {
		return t.Format("Jan 2006")
	})
	engine.AddFunc("formatDateTime", func(t time.Time) string {
		return t.Format("Jan 02, 2006 15:04")
	})
	engine.AddFunc("formatDateInput", func(t time.Time) string {
		return t.Format("2006-01-02")
	})
	engine.AddFunc("derefTime", func(t *time.Time) time.Time {
		if t == nil {
			return time.Time{}
		}
		return *t
	})
	engine.AddFunc("isZeroTime", func(t time.Time) bool {
		return t.IsZero()
	})
	engine.AddFunc("derefUint", func(p *uint) uint {
		if p == nil {
			return 0
		}
		return *p
	})
	engine.AddFunc("seq", func(n int) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i
		}
		return s
	})
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("sub", func(a, b int) int { return a - b })
	engine.AddFunc("appVersion", func() string { return cfg.App.Version })
	engine.AddFunc("humanSize", func(size int64) string {
		if size < 1024 {
			return fmt.Sprintf("%d B", size)
		} else if size < 1024*1024 {
			return fmt.Sprintf("%.1f KB", float64(size)/1024)
		}
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	})

	// 6. Fiber app.
	app := fiber.New(fiber.Config{
		AppName:   fmt.Sprintf("%s v%s", cfg.App.Name, cfg.App.Version),
		Views:     engine,
		BodyLimit: int(cfg.Upload.MaxResumeSize) * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).Render("public/error", fiber.Map{
				"Code":    code,
				"Message": err.Error(),
			}, "layouts/public_base")
		},
	})

	// 7. Global middleware.
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(compress.New())
	app.Use(middleware.Security())

	// 8. Static files (cache 1 year; browser will revalidate on hard refresh).
	app.Static("/static", "./web/static", fiber.Static{
		Compress:      true,
		MaxAge:        31536000,
		CacheDuration: 24 * time.Hour,
	})
	app.Static("/uploads", cfg.App.UploadDir, fiber.Static{
		Compress:      true,
		MaxAge:        31536000,
		CacheDuration: 24 * time.Hour,
	})

	// 9. Routes.
	router.RegisterRoutes(app, db)

	// 10. Start server in background, wait for Ctrl+C (SIGINT/SIGTERM).
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
		log.Printf("Starting %s on http://%s (Ctrl+C to stop)", cfg.App.Name, addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownTimeoutInt := cfg.App.ShutdownTimeout
	if shutdownTimeoutInt <= 0 {
		shutdownTimeoutInt = 10
	}

	log.Printf("Shutting down gracefully (%ds timeout)...\n", shutdownTimeoutInt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeoutInt)*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
	config.MyPortfolio.Close()
	log.Println("Server stopped")
}
