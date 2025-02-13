package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port string
	AllowedOrigins []string
	DatabaseURL string
	JWTSecret string
	JWTExpiration time.Duration
	Environment string
}

func LoadConfig() *Config {
	config := &Config{
		Port: "8080",
		Environment: "development",
		JWTExpiration: 24 * time.Hour,
		AllowedOrigins: []string{"http://localhost:5173"},
	}

	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if environment := os.Getenv("ENVIRONMENT"); environment != "" {
		config.Environment = environment
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL 
	} else {
		config.DatabaseURL = "host=localhost user=hafizh password=Sudarmi12 dbname=test_db port=5432 sslmode=disable"
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	} else {
		config.JWTSecret = "CedAfWIWdd2FimRPl8A6cwTwq9VS4QjzgyjWbpXt3Ib6a1qj0GHbJSxfcowETF3SqAJarBLcyYOseTIUhDA0sDJUfgHbDii2bUrJoFlZhvXEIAypfAi2HeTi9gfHBbJIXD7V8XHgk6CZ5b9tjV4KtvFwjirDMqAteyzWLSiuBpAqYm7JmdIufBysguvUxXW038x4qz4Jqf5894g2JlTve0CFwFCMLfFFxiYgiXX7HzvGgxRU7hC8GfuHkSSwYXhT"
	}

	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		config.AllowedOrigins = strings.Split(origins, ",")
	}

	return config


}