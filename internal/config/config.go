package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

const LevelTrace = slog.Level(12)

type config struct {
	Port      string
	Username  string
	Password  string
	Host      string
	DBPort    string
	DBName    string
	SSLMode   string
	JWTSecret string
	Stage     string
}

func (c *config) BuildDBConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.DBPort, c.DBName, c.SSLMode)
}

func Init() config {
	return config{
		Port:     getEnv("PORT", "8080"),
		Username: getEnv("PG_USERNAME", "postgres"),
		Password: getEnv("PG_PASSWORD", "password"),
		Host:     getEnv("PG_HOST", "0.0.0.0"),
		DBPort:   getEnv("PG_PORT", "5432"),
		DBName:   getEnv("PG_DATABASE", "testDB"),
		SSLMode:  getEnv("PG_SSLMODE", "disable"),
		Stage:    getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

func InitLogger(w io.Writer, stage string) *slog.Logger {
	levelNames := map[slog.Leveler]string{
		LevelTrace: "FATAL",
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}

	var handler slog.Handler = slog.NewTextHandler(w, opts)
	if stage == "production" {
		handler = slog.NewJSONHandler(w, opts)
	}

	l := slog.New(handler)
	return l
}
