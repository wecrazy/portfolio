<h1 align="center">My Portfolio</h1>

<p align="center">
  A full-featured, self-hosted personal portfolio &amp; blog built with
  <a href="https://gofiber.io">Go Fiber</a>,
  server-rendered HTML templates, HTMX, SQLite, and Redis.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/Fiber-v3-00ACD7?logo=go&logoColor=white" alt="Fiber v3" />
  <img src="https://img.shields.io/badge/SQLite-WAL-003B57?logo=sqlite&logoColor=white" alt="SQLite" />
  <img src="https://img.shields.io/badge/Redis-sessions-DC382D?logo=redis&logoColor=white" alt="Redis" />
  <img src="https://img.shields.io/badge/HTMX-dynamic%20UI-3366CC" alt="HTMX" />
  <img src="https://img.shields.io/badge/License-MIT-green" alt="MIT License" />
</p>

<p align="center">
  <b>Live Demo →</b> <a href="https://wecrazy.my.id">wecrazy.my.id</a>
</p>

---

## Features

| Area | Details |
|---|---|
| **Portfolio** | Projects, experience, skills, tech stacks, social links, upcoming items |
| **Blog** | Markdown posts with thumbnail/media/video/audio uploads, slug-based URLs |
| **Comments** | OAuth (Google & GitHub) login, real-time WebSocket broadcast, admin moderation |
| **Contact** | Form with email notification (SMTP), hCaptcha spam protection |
| **Admin Panel** | Full CRUD dashboard for every entity, file/media upload manager, server monitor & profiler |
| **i18n** | Multi-language support (English / Indonesian out of the box) with client-side locale switching |
| **Security** | Rate limiting (per-route), load shedding (CPU-aware), circuit breaker, security headers |
| **Config** | YAML-based, environment-aware (`dev`/`prod`), hot-reload via fsnotify |
| **Deployment** | Installable as a systemd service (Linux) or Windows SCM service |

## Tech Stack

- **Language** — Go 1.26
- **Web Framework** — [Fiber v3](https://gofiber.io) + HTML template engine
- **Database** — SQLite (WAL mode, GORM)
- **Session Store** — Redis
- **Frontend** — Server-rendered HTML, HTMX, vanilla JS/CSS
- **Logging** — Zap (structured JSON) + Lumberjack (rotation)
- **Validation** — go-playground/validator
- **Markdown** — Goldmark + Bluemonday (sanitization)
- **OAuth** — Google & GitHub via `golang.org/x/oauth2`
- **Captcha** — hCaptcha
- **Email** — gomail (SMTP)
- **Linting** — Revive

## Project Structure

```
my-portfolio/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # YAML config loading, hot-reload, validation
│   ├── database/        # SQLite + Redis initialization
│   ├── handler/         # Public route handlers (portfolio, blog, auth, comments, contact)
│   │   └── admin/       # Admin panel handlers (CRUD for all entities)
│   ├── hub/             # WebSocket broadcast hub (real-time comments)
│   ├── i18n/            # Internationalization
│   ├── middleware/       # Auth guards, security headers
│   ├── model/           # GORM models
│   ├── router/          # Route registration (public, admin, auth, API)
│   ├── seed/            # Database seeding (admin user, demo data)
│   └── service/         # Email, OAuth, upload services
├── pkg/                 # Reusable packages (crypto, file utils, pagination, sanitize)
├── swagger/             # OpenAPI spec (JSON)
├── web/
│   ├── locales/         # i18n translation files (en.yaml, id.yaml)
│   ├── static/          # CSS, JS, images
│   └── templates/       # Go HTML templates (admin, public, layouts, partials)
├── uploads/             # User-uploaded files (images, resume, video, audio)
├── data/                # SQLite database files
├── Makefile             # Build, dev, test, service management commands
└── go.mod
```

## Prerequisites

- **Go** >= 1.26
- **Redis** running locally (default `localhost:6379`)
- **Make** (optional, but recommended)

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/wecrazy/my-portfolio.git
cd my-portfolio

# 2. Initialize (create directories + download deps)
make init

# 3. Copy and edit the config
cp internal/config/my-portfolio.dev.yaml internal/config/my-portfolio.prod.yaml
# Edit the .yaml with your own secrets, SMTP, OAuth credentials, etc.

# 4. Seed the database with default admin + demo data
make seed

# 5. Run the server
make run
# → http://localhost:6969
```

## Configuration

Configuration lives in `internal/config/my-portfolio.<env>.yaml`. The environment is auto-detected from the `ENV` variable (defaults to `dev`).

Key sections:

| Section | Purpose |
|---|---|
| `app` | Name, host, port, base URL, secret key, debug toggle |
| `database` | SQLite DSN, connection pool settings |
| `admin` | Default admin credentials, session TTL, cookie settings |
| `oauth` | Google & GitHub OAuth client ID/secret/redirect |
| `smtp` | Email server for contact form notifications |
| `upload` | Max file sizes, allowed MIME types |
| `rate_limit` | Per-route request limits (contact, comments) |
| `redis` | Address, password, DB number |
| `i18n` | Default language, supported languages |
| `hcaptcha` | Site key, secret, enable/disable toggle |
| `log` | Directory, rotation, compression, stdout toggle |
| `owner` | Your name, title, bio, profile image, contact info |

Config changes are **hot-reloaded** — no server restart needed.

## Make Commands

```
make help
```

| Command | Description |
|---|---|
| `make run` | Run the server (clears stale static cache first) |
| `make build` | Compile binary to `bin/` |
| `make dev` | Hot-reload dev server via [Air](https://github.com/air-verse/air) |
| `make seed` | Seed database with default admin + demo data |
| `make deps` | `go mod tidy` + `go mod download` |
| `make revive` | Run the Revive linter |
| `make test` | Run all tests |
| `make db-reset` | Delete SQLite DB and re-seed from scratch |
| `make clean` | Remove build artifacts |
| `make clean-static` | Delete stale `.fiber.br`/`.fiber.gz` cache files |
| `make dirs` | Create required upload/data directories |
| `make init` | Full project initialization (dirs + deps) |
| `make install-service` | Build and install as a systemd/SCM service |
| `make uninstall-service` | Stop and remove the OS service |
| `make service-status` | Show current service status |

## Deploy as a System Service

```bash
# Build + install (Linux — requires sudo)
make install-service

# Check status
make service-status

# Remove
make uninstall-service
```

The installer creates a systemd unit on Linux and a Windows SCM service on Windows.

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Portfolio home page |
| `GET` | `/blog` | Blog listing |
| `GET` | `/blog/:slug` | Single blog post |
| `GET` | `/projects` | Projects (HTMX partial) |
| `GET` | `/upcoming` | Upcoming items (HTMX partial) |
| `GET` | `/resume` | View resume PDF |
| `GET` | `/resume/download` | Download resume PDF |
| `POST` | `/contact` | Submit contact form |
| `GET/POST` | `/comments` | View / post comments (OAuth required) |
| `GET` | `/ws/comments` | WebSocket for real-time comments |
| `GET` | `/auth/google` | Google OAuth login |
| `GET` | `/auth/github` | GitHub OAuth login |
| `POST` | `/api/translate` | Machine-translate content |
| `GET` | `/lang/:code` | Locale JSON for client-side i18n |
| `GET` | `/admin/*` | Admin panel (session-protected) |
| `GET` | `/swagger/*` | Swagger UI (admin-only) |
| `GET` | `/livez` | Liveness probe |
| `GET` | `/readyz` | Readiness probe |

## License

[MIT](LICENSE) — Wegirandol Histara Littu
