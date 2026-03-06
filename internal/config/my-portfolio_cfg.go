package config

// MyPortfolio is the singleton config accessor for the portfolio app.
// Follows the exact same pattern as service-platform's ServicePlatform variable.
var MyPortfolio = &configs[TypeMyPortfolio]{}

// TypeMyPortfolio represents the structure of my-portfolio.<env>.yaml
type TypeMyPortfolio struct {
	App struct {
		Name            string `yaml:"name" validate:"required"`
		Version         string `yaml:"version" validate:"required"`
		Host            string `yaml:"host" validate:"required"`
		Port            int    `yaml:"port" validate:"required"`
		BaseURL         string `yaml:"base_url" validate:"required"`
		Debug           bool   `yaml:"debug"`
		LogLevel        string `yaml:"log_level" validate:"required"`
		SecretKey       string `yaml:"secret_key" validate:"required"`
		StaticDir       string `yaml:"static_dir" validate:"required"`
		UploadDir       string `yaml:"upload_dir" validate:"required"`
		ShutdownTimeout int    `yaml:"shutdown_timeout" validate:"required"`
	} `yaml:"app" validate:"required"`

	Database struct {
		DSN             string `yaml:"dsn" validate:"required"`
		MaxIdleConns    int    `yaml:"max_idle_conns" validate:"required"`
		MaxOpenConns    int    `yaml:"max_open_conns" validate:"required"`
		ConnMaxLifetime int    `yaml:"conn_max_lifetime" validate:"required"`
		LogLevel        string `yaml:"log_level" validate:"required"`
	} `yaml:"database" validate:"required"`

	Admin struct {
		DefaultUsername string `yaml:"default_username" validate:"required"`
		DefaultPassword string `yaml:"default_password" validate:"required"`
		DefaultEmail    string `yaml:"default_email" validate:"required"`
		SessionTTL      int    `yaml:"session_ttl" validate:"required"`
		CookieName      string `yaml:"cookie_name" validate:"required"`
		CookieSecure    bool   `yaml:"cookie_secure"`
		CookieDomain    string `yaml:"cookie_domain"`
	} `yaml:"admin" validate:"required"`

	OAuth struct {
		Google struct {
			ClientID     string `yaml:"client_id" validate:"required"`
			ClientSecret string `yaml:"client_secret" validate:"required"`
			RedirectURL  string `yaml:"redirect_url" validate:"required"`
		} `yaml:"google" validate:"required"`
		GitHub struct {
			ClientID     string `yaml:"client_id" validate:"required"`
			ClientSecret string `yaml:"client_secret" validate:"required"`
			RedirectURL  string `yaml:"redirect_url" validate:"required"`
		} `yaml:"github" validate:"required"`
	} `yaml:"oauth" validate:"required"`

	SMTP struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required"`
		Username string `yaml:"username" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		From     string `yaml:"from" validate:"required"`
		To       string `yaml:"to" validate:"required"`
	} `yaml:"smtp" validate:"required"`

	Upload struct {
		MaxImageSize       int64    `yaml:"max_image_size" validate:"required"`
		MaxResumeSize      int64    `yaml:"max_resume_size" validate:"required"`
		MaxVideoSize       int64    `yaml:"max_video_size"`
		MaxAudioSize       int64    `yaml:"max_audio_size"`
		AllowedImageTypes  []string `yaml:"allowed_image_types" validate:"required"`
		AllowedResumeTypes []string `yaml:"allowed_resume_types" validate:"required"`
		AllowedVideoTypes  []string `yaml:"allowed_video_types"`
		AllowedAudioTypes  []string `yaml:"allowed_audio_types"`
	} `yaml:"upload" validate:"required"`

	RateLimit struct {
		Enabled     bool `yaml:"enabled"`
		ContactForm int  `yaml:"contact_form" validate:"required"`
		Comments    int  `yaml:"comments" validate:"required"`
	} `yaml:"rate_limit" validate:"required"`

	Redis struct {
		Addr     string `yaml:"addr" validate:"required"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis" validate:"required"`

	I18n struct {
		DefaultLang    string   `yaml:"default_lang" validate:"required"`
		SupportedLangs []string `yaml:"supported_langs" validate:"required,min=1"`
	} `yaml:"i18n" validate:"required"`

	TableNames struct {
		Admins          string `yaml:"admins"`
		Owners          string `yaml:"owners"`
		Projects        string `yaml:"projects"`
		Experiences     string `yaml:"experiences"`
		Skills          string `yaml:"skills"`
		SocialLinks     string `yaml:"social_links"`
		UploadedFiles   string `yaml:"uploaded_files"`
		OAuthUsers      string `yaml:"oauth_users"`
		Comments        string `yaml:"comments"`
		ContactMessages string `yaml:"contact_messages"`
		TechStacks      string `yaml:"tech_stacks"`
		Posts           string `yaml:"posts"`
		UpcomingItems   string `yaml:"upcoming_items"`
	} `yaml:"table_names"`

	HCaptcha struct {
		SiteKey string `yaml:"site_key"`
		Secret  string `yaml:"secret"`
		// Enabled controls whether hcaptcha is verified on admin login.
		// Set to false in dev to skip verification.
		Enabled bool `yaml:"enabled"`
	} `yaml:"hcaptcha"`

	Log struct {
		Dir        string `yaml:"dir" validate:"required"`
		Filename   string `yaml:"filename" validate:"required"`
		MaxSizeMB  int    `yaml:"max_size_mb" validate:"required"`
		MaxBackups int    `yaml:"max_backups" validate:"required"`
		MaxAgeDays int    `yaml:"max_age_days" validate:"required"`
		Compress   bool   `yaml:"compress"`
		Stdout     bool   `yaml:"stdout"`
	} `yaml:"log" validate:"required"`

	Owner struct {
		Name         string `yaml:"name" validate:"required"`
		Title        string `yaml:"title" validate:"required"`
		Tagline      string `yaml:"tagline"`
		Bio          string `yaml:"bio" validate:"required"`
		ProfileImage string `yaml:"profile_image" validate:"required"`
		Email        string `yaml:"email" validate:"required,email"`
		Phone        string `yaml:"phone" validate:"required"`
		Location     string `yaml:"location" validate:"required"`
	} `yaml:"owner" validate:"required"`

	Cert struct {
		CertFile string `yaml:"cert_file" validate:"required"`
		KeyFile  string `yaml:"key_file" validate:"required"`
	} `yaml:"cert" validate:"required"`
}
