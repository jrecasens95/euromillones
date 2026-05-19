package config

import "os"

type AppConfig struct {
	Port               string
	DrawsDataDir       string
	CORSAllowedOrigins string
	JWTSecret          string
}

var Current AppConfig

func Load() {
	Current = AppConfig{
		Port:               getEnv("PORT", "4000"),
		DrawsDataDir:       getEnv("DRAWS_DATA_DIR", "data"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
