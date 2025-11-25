package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig contains all application configuration
type AppConfig struct {
	App       App       `mapstructure:"app"`
	Server    Server    `mapstructure:"server"`
	Database  Database  `mapstructure:"database"`
	Redis     Redis     `mapstructure:"redis"`
	Logger    Logger    `mapstructure:"logger"`
	JWT       JWT       `mapstructure:"jwt"`
	Metrics   Metrics   `mapstructure:"metrics"`
	Tracing   Tracing   `mapstructure:"tracing"`
	RateLimit RateLimit `mapstructure:"rate_limit"`
	Messaging Messaging `mapstructure:"messaging"`
	External  External  `mapstructure:"external"`
	Feature   Feature   `mapstructure:"feature"`
}

type App struct {
	Name         string        `mapstructure:"name"`
	Version      string        `mapstructure:"version"`
	Environment  string        `mapstructure:"environment"`
	Debug        bool          `mapstructure:"debug"`
	Timezone     string        `mapstructure:"timezone"`
	GracefulStop time.Duration `mapstructure:"graceful_stop"`
}

type Server struct {
	HTTP ServerHTTP `mapstructure:"http"`
	GRPC ServerGRPC `mapstructure:"grpc"`
}

type ServerHTTP struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
	TLS            TLS           `mapstructure:"tls"`
	CORS           CORS          `mapstructure:"cors"`
}

type ServerGRPC struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type TLS struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type CORS struct {
	Enabled        bool     `mapstructure:"enabled"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

type Database struct {
	Primary DatabaseConnection `mapstructure:"primary"`
	Replica DatabaseConnection `mapstructure:"replica"`
}

type DatabaseConnection struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogLevel        string        `mapstructure:"log_level"`
}

type Redis struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type Logger struct {
	Level       string   `mapstructure:"level"`
	Encoding    string   `mapstructure:"encoding"`
	OutputPaths []string `mapstructure:"output_paths"`
	ErrorPaths  []string `mapstructure:"error_paths"`
	Development bool     `mapstructure:"development"`
}

type JWT struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      string        `mapstructure:"audience"`
}

type Metrics struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

type Tracing struct {
	Enabled     bool    `mapstructure:"enabled"`
	ServiceName string  `mapstructure:"service_name"`
	Endpoint    string  `mapstructure:"endpoint"`
	SampleRate  float64 `mapstructure:"sample_rate"`
}

type RateLimit struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

type Messaging struct {
	NATS NATSConfig `mapstructure:"nats"`
}

type NATSConfig struct {
	URL           string        `mapstructure:"url"`
	MaxReconnects int           `mapstructure:"max_reconnects"`
	ReconnectWait time.Duration `mapstructure:"reconnect_wait"`
	Timeout       time.Duration `mapstructure:"timeout"`
	DrainTimeout  time.Duration `mapstructure:"drain_timeout"`
}

type External struct {
	PaymentService PaymentServiceConfig `mapstructure:"payment_service"`
	EmailService   EmailServiceConfig   `mapstructure:"email_service"`
	FileStorage    FileStorageConfig    `mapstructure:"file_storage"`
}

type PaymentServiceConfig struct {
	BaseURL    string        `mapstructure:"base_url"`
	APIKey     string        `mapstructure:"api_key"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
}

type EmailServiceConfig struct {
	Provider string     `mapstructure:"provider"`
	SMTP     SMTPConfig `mapstructure:"smtp"`
	APIKey   string     `mapstructure:"api_key"`
	From     string     `mapstructure:"from"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	UseTLS   bool   `mapstructure:"tls"`
}

type FileStorageConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BucketName      string `mapstructure:"bucket_name"`
	CDNUrl          string `mapstructure:"cdn_url"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

type Feature struct {
	EnableSwagger     bool `mapstructure:"enable_swagger"`
	EnablePprof       bool `mapstructure:"enable_pprof"`
	EnableHealthCheck bool `mapstructure:"enable_health_check"`
	EnableMetrics     bool `mapstructure:"enable_metrics"`
}

// GetDSN returns database connection string
func (db *DatabaseConnection) GetDSN() string {
	switch db.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
			db.Host, db.Port, db.Username, db.Password, db.Database, db.SSLMode)
	default:
		return ""
	}
}

// LoadConfig loads configuration from various sources
func LoadConfig(configPath string) (*AppConfig, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.AddConfigPath("../configs")
	}

	// Set defaults
	setDefaults(v)

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "go-mvc-enterprise")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.timezone", "UTC")
	v.SetDefault("app.graceful_stop", "30s")

	// Server defaults
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.http.read_timeout", "30s")
	v.SetDefault("server.http.write_timeout", "30s")
	v.SetDefault("server.http.idle_timeout", "120s")
	v.SetDefault("server.http.max_header_bytes", 1048576)

	// CORS defaults
	v.SetDefault("server.http.cors.enabled", true)
	v.SetDefault("server.http.cors.allowed_origins", []string{"*"})
	v.SetDefault("server.http.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("server.http.cors.allowed_headers", []string{"*"})
	v.SetDefault("server.http.cors.max_age", 86400)

	// Database defaults
	v.SetDefault("database.primary.driver", "postgres")
	v.SetDefault("database.primary.host", "localhost")
	v.SetDefault("database.primary.port", 5432)
	v.SetDefault("database.primary.username", "postgres")
	v.SetDefault("database.primary.password", "postgres")
	v.SetDefault("database.primary.database", "gomvc")
	v.SetDefault("database.primary.ssl_mode", "disable")
	v.SetDefault("database.primary.max_open_conns", 25)
	v.SetDefault("database.primary.max_idle_conns", 25)
	v.SetDefault("database.primary.conn_max_lifetime", "5m")
	v.SetDefault("database.primary.conn_max_idle_time", "5m")
	v.SetDefault("database.primary.log_level", "warn")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")
	v.SetDefault("redis.idle_timeout", "5m")

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.encoding", "json")
	v.SetDefault("logger.output_paths", []string{"stdout"})
	v.SetDefault("logger.error_paths", []string{"stderr"})
	v.SetDefault("logger.development", false)

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.access_expiry", "15m")
	v.SetDefault("jwt.refresh_expiry", "7d")
	v.SetDefault("jwt.issuer", "go-mvc-enterprise")
	v.SetDefault("jwt.audience", "go-mvc-enterprise")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9090)

	// Tracing defaults
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.service_name", "go-mvc-enterprise")
	v.SetDefault("tracing.sample_rate", 0.1)

	// Rate limiting defaults
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.rps", 100)
	v.SetDefault("rate_limit.burst", 200)

	// Feature flags defaults
	v.SetDefault("feature.enable_swagger", true)
	v.SetDefault("feature.enable_pprof", false)
	v.SetDefault("feature.enable_health_check", true)
	v.SetDefault("feature.enable_metrics", true)
}

// validateConfig validates the configuration
func validateConfig(config *AppConfig) error {
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if config.Server.HTTP.Port <= 0 || config.Server.HTTP.Port > 65535 {
		return fmt.Errorf("server.http.port must be between 1 and 65535")
	}

	if config.JWT.Secret == "" && config.App.Environment == "production" {
		return fmt.Errorf("jwt.secret is required in production")
	}

	return nil
}
