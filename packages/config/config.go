package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config contains all configuration parameters for the application.
type Config struct {
	Env                string
	Port               string
	DatabaseURL        string
	DatabaseReplicaURL string
	RedisURL           string
	UploadProvider     string
	UploadEndpoint     string
	UploadAccessKey    string
	UploadSecretKey    string
	UploadRegion       string
	UploadSSL          bool
	UploadPublicURL    string
	JWTSecret          string
	ModulesDir         string

	// Cache config
	CacheProvider string
	CacheURL      string

	// Queue config
	QueueProvider string
	QueueURL      string

	// EventBus config
	EventsProvider string
	EventsURL      string

	// Monitoring config
	MonitoringProvider string
	MonitoringDSN      string
	MonitoringEndpoint string

	// Payment config
	PaymentProvider  string
	PaymentAPIKey    string
	PaymentClientID  string
	PaymentMerchant  string
	PaymentSecretKey string
	PaymentSandbox   bool

	// Search config
	SearchProvider string
	SearchURL      string
	SearchAPIKey   string
}

// Load loads configuration from .env file (if present) and environment variables.
func Load() *Config {
	// Automatically find and load .env file at project root if it exists
	_ = godotenv.Load()

	return &Config{
		Env:                getEnv("ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:secret@localhost:5432/ecommerce?sslmode=disable"),
		DatabaseReplicaURL: os.Getenv("DATABASE_REPLICA_URL"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		UploadProvider:     getEnv("UPLOAD_PROVIDER", "minio"),
		UploadEndpoint:     getEnv("UPLOAD_ENDPOINT", "localhost:9000"),
		UploadAccessKey:    getEnv("UPLOAD_ACCESS_KEY", "minioadmin"),
		UploadSecretKey:    getEnv("UPLOAD_SECRET_KEY", "minioadmin"),
		UploadRegion:       getEnv("UPLOAD_REGION", ""),
		UploadSSL:          getEnv("UPLOAD_SSL", "false") == "true",
		UploadPublicURL:    getEnv("UPLOAD_PUBLIC_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production"),
		ModulesDir:         getEnv("MODULES_DIR", "modules"),

		// Cache config
		CacheProvider: getEnv("CACHE_PROVIDER", "inmemory"),
		CacheURL:      getEnv("CACHE_URL", "redis://localhost:6379"),

		// Queue config
		QueueProvider: getEnv("QUEUE_PROVIDER", "inmemory"),
		QueueURL:      getEnv("QUEUE_URL", ""),

		// EventBus config
		EventsProvider: getEnv("EVENTS_PROVIDER", "log"),
		EventsURL:      getEnv("EVENTS_URL", ""),

		// Monitoring config
		MonitoringProvider: getEnv("MONITORING_PROVIDER", "console"),
		MonitoringDSN:      getEnv("MONITORING_DSN", ""),
		MonitoringEndpoint: getEnv("MONITORING_ENDPOINT", ""),

		// Payment config
		PaymentProvider:  getEnv("PAYMENT_PROVIDER", "stripe"),
		PaymentAPIKey:    getEnv("PAYMENT_API_KEY", ""),
		PaymentClientID:  getEnv("PAYMENT_CLIENT_ID", ""),
		PaymentMerchant:  getEnv("PAYMENT_MERCHANT_ID", ""),
		PaymentSecretKey: getEnv("PAYMENT_SECRET_KEY", ""),
		PaymentSandbox:   getEnv("PAYMENT_SANDBOX", "true") == "true",

		// Search config
		SearchProvider: getEnv("SEARCH_PROVIDER", "inmemory"),
		SearchURL:      getEnv("SEARCH_URL", ""),
		SearchAPIKey:   getEnv("SEARCH_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
