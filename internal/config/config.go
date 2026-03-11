package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigins string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8000"),
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/zguba_gov?parseTime=true&charset=utf8mb4"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:4200,http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
