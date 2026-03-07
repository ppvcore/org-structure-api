package cfg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
}

type PostgresConfig struct {
	DatabaseURL string
}

type ServerConfig struct {
	Addr string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	cfg := &Config{
		Postgres: PostgresConfig{
			DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/orgdb?sslmode=disable"),
		},
		Server: ServerConfig{
			Addr: getEnv("ADDR", ":8080"),
		},
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
