package config

import (
	"fmt"
	"os"
)

type config struct {
	Port      string
	Username  string
	Password  string
	Host      string
	DBPort    string
	DBName    string
	SSLMode   string
	JWTSecret string
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
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
