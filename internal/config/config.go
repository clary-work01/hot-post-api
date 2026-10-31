package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   []byte
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// 所有環境變數只在 config 讀一次，其他地方只拿 cfg，也只在這裡集中驗證
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file, using environment variables", "error", err)
	}

	cfg := Config{
		AppEnv:      getEnv("APP_ENV", "dev"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   []byte(os.Getenv("JWT_SECRET")),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	var errs []error

	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if len(c.JWTSecret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET must be at least 32 bytes"))
	}

	switch c.AppEnv {
	case "dev", "production":
	default:
		// %q: 把字串加上雙引號輸出
		errs = append(errs, fmt.Errorf("APP_ENV must be dev or production, got %q", c.AppEnv))
	}

	// 一次回報所有錯，不用改一個重啟一次
	return errors.Join(errs...)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
