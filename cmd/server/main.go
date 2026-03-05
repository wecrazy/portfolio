// Package main is the entry point for the portfolio server.
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3/middleware/static"

	"my-portfolio/internal/config"
	"my-portfolio/internal/database"
	"my-portfolio/internal/hub"
	"my-portfolio/internal/middleware"
	"my-portfolio/internal/model"
	"my-portfolio/internal/router"
	"my-portfolio/internal/seed"
	"my-portfolio/pkg/installer"

	contribzap "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/lumberjack.v2"
)

// serviceInstallerCfg is the installer.Config used by the service flags.
var serviceInstallerCfg = installer.Config{
	Name:        "my-portfolio",
	DisplayName: "My Portfolio Server",
	Description: "Personal portfolio web server (Fiber / SQLite)",
}

// handleServiceFlags checks os.Args for --install / --uninstall /
// --service-status and, if present, runs the matching installer action then
// exits.  Returns false when no service flag was found.
func handleServiceFlags() bool {
	argToAction := map[string]string{
		"--install":        "install",
		"--uninstall":      "uninstall",
		"--service-status": "status",
	}
	for _, arg := range os.Args[1:] {
		action, ok := argToAction[arg]
		if !ok {
			continue
		}
		if err := installer.Execute(action, serviceInstallerCfg); err != nil {
			log.Fatalf("installer: %v", err)
		}
		return true
	}
	return false
}

func main() {
	// 0. Handle service-management flags before anything else.
	//    Usage:
	//      sudo ./bin/my-portfolio --install
	//      sudo ./bin/my-portfolio --uninstall
	//           ./bin/my-portfolio --service-status
	if handleServiceFlags() {
		os.Exit(0)
	}

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

	// 3b. Init Redis (admin session store — survives Go server restarts).
	rdb := database.InitRedis(cfg)

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
	engine.AddFunc("currentYear", func() int { return time.Now().Year() })
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
		ErrorHandler: func(c fiber.Ctx, err error) error {
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

	// 7. Build structured logger: JSON file (rotated by lumberjack) + optional stdout.
	if err := os.MkdirAll(cfg.Log.Dir, 0o755); err != nil {
		log.Fatalf("Failed to create log directory %s: %v", cfg.Log.Dir, err)
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, cfg.Log.Filename),
		MaxSize:    cfg.Log.MaxSizeMB,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAgeDays,
		Compress:   cfg.Log.Compress,
	})

	logLevel := zapcore.InfoLevel
	if cfg.App.Debug {
		logLevel = zapcore.DebugLevel
	}

	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, fileSyncer, logLevel),
	}
	if cfg.Log.Stdout {
		consoleEnc := zap.NewDevelopmentEncoderConfig()
		consoleEnc.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEnc.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEnc),
			zapcore.AddSync(os.Stdout),
			logLevel,
		))
	}

	zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer zapLogger.Sync() //nolint:errcheck

	// 8. Global middleware.
	app.Use(recover.New())
	app.Use(contribzap.New(contribzap.Config{
		Logger: zapLogger,
		Fields: []string{
			"latency",
			"status",
			"method",
			"url",
			"ip",
			"host",
			"route",
			"ua",
			"referer",
			"queryParams",
			"bytesSent",
			"bytesReceived",
			"error",
		},
		SkipURIs: []string{"/livez", "/readyz"},
	}))
	app.Use(compress.New())
	app.Use(middleware.Security())

	// 9. Static files (cache 1 year; browser will revalidate on hard refresh).
	app.Get("/static*", static.New("./web/static", static.Config{
		Compress:      true,
		MaxAge:        31536000,
		CacheDuration: 24 * time.Hour,
	}))
	app.Get("/uploads*", static.New(cfg.App.UploadDir, static.Config{
		Compress:      true,
		MaxAge:        31536000,
		CacheDuration: 24 * time.Hour,
	}))

	// 10. WebSocket / broadcast hub.
	h := hub.New()

	// 11. Routes (loadshed, circuitbreaker, WS, hcaptcha are registered inside).
	router.RegisterRoutes(app, db, rdb, h)

	// 12. Shutdown hooks — broadcast event so clients show a toast.
	app.Hooks().OnPreShutdown(func() error {
		h.Broadcast(hub.Event{
			Type: hub.EventShutdown,
			Data: map[string]any{"message": "Server is restarting — please wait…"},
		})
		return nil
	})

	// 13. Start server in background, wait for Ctrl+C (SIGINT/SIGTERM).
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
