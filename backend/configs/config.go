package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logging  LoggingConfig
	CORS     CORSConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Env  string
	Name string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
	Path  string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Enable environment variables
	viper.AutomaticEnv()

	// Set defaults and env bindings
	setDefaults()

	// Determine environment (from env vars or default)
	env := viper.GetString("app.env")
	if env == "" {
		env = "development"
	}

	// Load environment-specific config file: config.<env>.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.SetConfigName(fmt.Sprintf("config.%s", env))

	// Try to read config file (optional)
	_ = viper.ReadInConfig()

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Handle CORS_ALLOWED_ORIGINS from environment variable
	// Viper reads it as string, but we need []string
	// Check if CORS_ALLOWED_ORIGINS env var is set and parse it if it's a string
	if viper.IsSet("cors.allowed_origins") {
		corsOriginsRaw := viper.Get("cors.allowed_origins")
		if corsOriginsStr, ok := corsOriginsRaw.(string); ok {
			// It's a string from env var, parse it
			if corsOriginsStr == "*" {
				cfg.CORS.AllowedOrigins = []string{"*"}
			} else {
				// Split by comma and trim spaces
				origins := strings.Split(corsOriginsStr, ",")
				cfg.CORS.AllowedOrigins = make([]string, 0, len(origins))
				for _, origin := range origins {
					trimmed := strings.TrimSpace(origin)
					if trimmed != "" {
						cfg.CORS.AllowedOrigins = append(cfg.CORS.AllowedOrigins, trimmed)
					}
				}
			}
		}
		// If it's already a slice (from config file), cfg.CORS.AllowedOrigins is already set correctly
	}

	return cfg, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.name", "lexigo")

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "10s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database", "lexigo")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_conns", 25)
	viper.SetDefault("database.min_conns", 5)
	viper.SetDefault("database.max_conn_lifetime", "5m")
	viper.SetDefault("database.max_conn_idle_time", "1m")

	// JWT defaults
	viper.SetDefault("jwt.secret", "change-me-in-production")
	viper.SetDefault("jwt.expiration", "24h")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.path", "/var/log/language-coach")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:3300"})

	// Environment variable mappings
	// Viper automatically maps environment variables, but we need to set up the key replacer
	// Since viper.NewReplacer doesn't exist in newer versions, we'll handle it differently
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("server.port", "APP_PORT")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.database", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("logging.level", "LOG_LEVEL")
	viper.BindEnv("logging.path", "LOG_PATH")
	viper.BindEnv("cors.allowed_origins", "CORS_ALLOWED_ORIGINS")
}
