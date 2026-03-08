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

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"my-portfolio/internal/config"
	"my-portfolio/internal/database"
	"my-portfolio/internal/hub"
	"my-portfolio/internal/middleware"
	"my-portfolio/internal/router"
	"my-portfolio/internal/seed"
	"my-portfolio/pkg/find"
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
	// service flag handling is deliberately the only logic here so the
	// entrypoint remains trivial.  revive will no longer complain about
	// complexity because all of the "real" work is delegated to smaller
	// helper functions.
	if handleServiceFlags() {
		os.Exit(0)
	}

	if err := run(); err != nil {
		log.Fatalf("startup failed: %v", err)
	}
}

// run orchestrates all of the initialization steps.  it is still fairly
// linear, but each section has been extracted into its own helper so the
// resulting cyclomatic complexity is tiny.
func run() error {
	cfg := initConfig()
	db, rdb := initDatabases(cfg)
	seed.SeedIfNeeded(db, cfg)
	engine := initTemplates(cfg)
	logger := initLogger(cfg)
	app := initFiber(appConfig{cfg: cfg, engine: engine, logger: logger})

	h := hub.New()
	router.RegisterRoutes(app, db, rdb, h)
	registerShutdownHook(app, h)

	certFile, keyFile, err := findCertFiles(cfg)
	if err != nil {
		return err
	}

	startServer(app, cfg, certFile, keyFile)
	waitForSignal(app, cfg)

	config.MyPortfolio.Close()
	return nil
}

// The following helper types and functions were extracted from the original
// main body.  Each can more easily be tested and reviewed in isolation.

type appConfig struct {
	cfg    config.TypeMyPortfolio
	engine *html.Engine
	logger *zap.Logger
}

func initConfig() config.TypeMyPortfolio {
	config.MyPortfolio.MustInit("my-portfolio")
	go config.MyPortfolio.Watch()
	return config.MyPortfolio.Get()
}

func initDatabases(cfg config.TypeMyPortfolio) (*gorm.DB, *redis.Client) {
	db := database.InitSQLite(cfg)
	database.AutoMigrate(db)
	return db, database.InitRedis(cfg)
}

func initTemplates(cfg config.TypeMyPortfolio) *html.Engine {
	engine := html.New("./web/templates", ".html")
	engine.Reload(cfg.App.Debug)
	// helper functions moved into a separate function to keep this one short
	addTemplateFuncs(engine, cfg)
	return engine
}

func addTemplateFuncs(engine *html.Engine, cfg config.TypeMyPortfolio) {
	engine.AddFunc("upper", strings.ToUpper)
	engine.AddFunc("lower", strings.ToLower)
	engine.AddFunc("contains", strings.Contains)
	engine.AddFunc("hasSuffix", strings.HasSuffix)
	engine.AddFunc("split", strings.Split)
	engine.AddFunc("join", strings.Join)
	engine.AddFunc("safeHTML", func(s string) template.HTML { return template.HTML(s) })
	engine.AddFunc("formatDate", func(t time.Time) string { return t.Format("Jan 2006") })
	engine.AddFunc("formatDateTime", func(t time.Time) string { return t.Format("Jan 02, 2006 15:04") })
	engine.AddFunc("formatDateInput", func(t time.Time) string { return t.Format("2006-01-02") })
	engine.AddFunc("derefTime", func(t *time.Time) time.Time {
		if t == nil {
			return time.Time{}
		}
		return *t
	})
	engine.AddFunc("isZeroTime", func(t time.Time) bool { return t.IsZero() })
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
}

func initLogger(cfg config.TypeMyPortfolio) *zap.Logger {
	if err := os.MkdirAll(cfg.Log.Dir, 0o755); err != nil {
		log.Fatalf("Failed to create log directory %s: %v", cfg.Log.Dir, err)
	}
	level := determineLogLevel(cfg)
	cores := buildLoggerCores(cfg, level)
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.WithCaller(false),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return logger
}

func determineLogLevel(cfg config.TypeMyPortfolio) zapcore.Level {
	if cfg.App.Debug {
		return zapcore.DebugLevel
	}
	switch strings.ToUpper(cfg.App.LogLevel) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN", "WARNING":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		log.Fatalf("Unknown log level %s, defaulting to INFO", cfg.App.LogLevel)
		return zapcore.InfoLevel // unreachable
	}
}

func buildLoggerCores(cfg config.TypeMyPortfolio, level zapcore.Level) []zapcore.Core {
	fileEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, cfg.Log.Filename),
		MaxSize:    cfg.Log.MaxSizeMB,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAgeDays,
		Compress:   cfg.Log.Compress,
	})
	cores := []zapcore.Core{zapcore.NewCore(fileEnc, fileSync, level)}
	if cfg.Log.Stdout {
		enc := zap.NewDevelopmentEncoderConfig()
		enc.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(enc), zapcore.AddSync(os.Stdout), level))
	}
	return cores
}

func initFiber(cfg appConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:   fmt.Sprintf("%s v%s", cfg.cfg.App.Name, cfg.cfg.App.Version),
		Views:     cfg.engine,
		BodyLimit: int(cfg.cfg.Upload.MaxResumeSize) * 1024 * 1024,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).Render("public/error", fiber.Map{
				"Code":       code,
				"Message":    err.Error(),
				"RetryAfter": 0,
			})
		},
	})
	// middleware
	app.Use(recover.New())
	app.Use(contribzap.New(contribzap.Config{Logger: cfg.logger, Fields: []string{"latency", "status", "method", "url", "ip", "host", "route", "ua", "referer", "queryParams", "bytesSent", "bytesReceived", "error"}, SkipURIs: []string{"/livez", "/readyz"}}))
	app.Use(compress.New())
	app.Use(middleware.Security())
	// static
	app.Get("/static*", static.New("./web/static", static.Config{Compress: true, MaxAge: 31536000, CacheDuration: 24 * time.Hour}))
	app.Get("/uploads*", static.New(cfg.cfg.App.UploadDir, static.Config{Compress: true, MaxAge: 31536000, CacheDuration: 24 * time.Hour}))
	return app
}

func registerShutdownHook(app *fiber.App, h *hub.Hub) {
	app.Hooks().OnPreShutdown(func() error {
		h.Broadcast(hub.Event{Type: hub.EventShutdown, Data: map[string]any{"message": "Server is restarting — please wait…"}})
		return nil
	})
}

func findCertFiles(cfg config.TypeMyPortfolio) (string, string, error) {
	if cfg.Cert.CertFile == "" || cfg.Cert.KeyFile == "" {
		return "", "", fmt.Errorf("cannot find cert and key file")
	}
	certFile, err := find.ValidDirectory([]string{cfg.Cert.CertFile, filepath.Join("..", cfg.Cert.CertFile), filepath.Join("..", "..", cfg.Cert.CertFile), filepath.Join("..", "..", "..", cfg.Cert.CertFile)})
	if err != nil {
		return "", "", err
	}
	keyFile, err := find.ValidDirectory([]string{cfg.Cert.KeyFile, filepath.Join("..", cfg.Cert.KeyFile), filepath.Join("..", "..", cfg.Cert.KeyFile), filepath.Join("..", "..", "..", cfg.Cert.KeyFile)})
	if err != nil {
		return "", "", err
	}
	if certFile == "" || keyFile == "" {
		return "", "", fmt.Errorf("certificate or key file not found")
	}
	return certFile, keyFile, nil
}

func startServer(app *fiber.App, cfg config.TypeMyPortfolio, certFile, keyFile string) {
	go func() {
		cfgModePath := cfg.App.ConfigModePath
		confYAMLPath, err := find.ValidDirectory([]string{
			cfgModePath,
			filepath.Join("..", cfgModePath),
			filepath.Join("..", "..", cfgModePath),
			filepath.Join("..", "..", "..", cfgModePath),
		})
		if err != nil {
			log.Fatalf("cannot find conf.yaml: %v", err)
		}

		var confDevMode bool = true
		if prod, err := config.FileContainsProd(confYAMLPath); err == nil && prod {
			confDevMode = false
		}

		addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
		log.Printf("Starting %s on http://%s (Ctrl+C to stop)\n", cfg.App.Name, addr)
		tlsCfg := fiber.ListenConfig{CertFile: certFile, CertKeyFile: keyFile}

		if confDevMode {
			log.Println("Running in `dev` mode (no TLS)")
			if err := app.Listen(addr); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		} else {
			log.Println("Running in `prod` mode (TLS enabled)")
			if err := app.Listen(addr, tlsCfg); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}

	}()
}

func waitForSignal(app *fiber.App, cfg config.TypeMyPortfolio) {
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
}
