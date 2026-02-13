package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
	LLM      LLMConfig
	Trending TrendingConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	ShutdownTimeout  time.Duration
}

type DatabaseConfig struct {
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	MaxConnections     int
	MaxIdleConnections int
	ConnectionLifetime time.Duration
}

type APIConfig struct {
	DefaultLimit int
	MaxLimit     int
	EnableCORS   bool
}

type LLMConfig struct {
	SummaryModel        string
	SummaryMaxTokens    int
	SummaryTemperature  float32
	IntentModel         string
	IntentMaxTokens     int
	IntentTemperature   float32
	OpenAIAPIKey        string
}

type TrendingConfig struct {
	CacheTTL                   time.Duration
	DefaultRadiusKm            int
	DefaultTimeWindowHours     int
	EventWeights               map[string]float64
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	// Set default values from config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Override with environment variables
	viper.SetEnvPrefix("SIMPLNEWS")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:             getIntEnv("SERVER_PORT", 8080),
			ReadTimeout:      getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout:     getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout:  getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:               getStringEnv("DATABASE_HOST", "localhost"),
			Port:               getIntEnv("DATABASE_PORT", 5432),
			Name:               getStringEnv("DATABASE_NAME", "simplnews"),
			User:               getStringEnv("DATABASE_USER", "simplnews_user"),
			Password:           getStringEnv("DATABASE_PASSWORD", "changeme123"),
			MaxConnections:     getIntEnv("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections: getIntEnv("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnectionLifetime: getDurationEnv("DB_CONNECTION_LIFETIME", 5*time.Minute),
		},
		API: APIConfig{
			DefaultLimit: getIntEnv("API_DEFAULT_LIMIT", 5),
			MaxLimit:     getIntEnv("API_MAX_LIMIT", 20),
			EnableCORS:   true,
		},
		LLM: LLMConfig{
			SummaryModel:       getStringEnv("LLM_SUMMARY_MODEL", "gpt-3.5-turbo-16k"),
			SummaryMaxTokens:   getIntEnv("LLM_SUMMARY_MAX_TOKENS", 150),
			SummaryTemperature: getFloat32Env("LLM_SUMMARY_TEMPERATURE", 0.3),
			IntentModel:        getStringEnv("LLM_INTENT_MODEL", "gpt-3.5-turbo-16k"),
			IntentMaxTokens:    getIntEnv("LLM_INTENT_MAX_TOKENS", 300),
			IntentTemperature:  getFloat32Env("LLM_INTENT_TEMPERATURE", 0.1),
			OpenAIAPIKey:       getStringEnv("OPENAI_API_KEY", ""),
		},
		Trending: TrendingConfig{
			CacheTTL:               getDurationEnv("TRENDING_CACHE_TTL", 5*time.Minute),
			DefaultRadiusKm:        getIntEnv("TRENDING_DEFAULT_RADIUS_KM", 100),
			DefaultTimeWindowHours: getIntEnv("TRENDING_DEFAULT_TIME_WINDOW_HOURS", 24),
			EventWeights: map[string]float64{
				"view":  1.0,
				"click": 2.0,
				"share": 3.0,
			},
		},
		Logging: LoggingConfig{
			Level:  getStringEnv("LOG_LEVEL", "info"),
			Format: getStringEnv("LOG_FORMAT", "json"),
			Output: getStringEnv("LOG_OUTPUT", "stdout"),
		},
	}

	if cfg.LLM.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if val := viper.GetInt(key); val != 0 {
		return val
	}
	return defaultValue
}

func getFloat32Env(key string, defaultValue float32) float32 {
	if val := viper.GetFloat64(key); val != 0 {
		return float32(val)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if val := viper.GetDuration(key); val != 0 {
		return val
	}
	return defaultValue
}
