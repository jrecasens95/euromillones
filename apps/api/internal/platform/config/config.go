package config

import (
	"os"
	"strings"
)

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
		CORSAllowedOrigins: allowedOrigins(),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}

func allowedOrigins() string {
	return normalizeOrigins(getEnvFrom(
		[]string{"CORS_ALLOWED_ORIGINS", "FRONTEND_URL", "PUBLIC_FRONTEND_URL", "VERCEL_URL"},
		"http://localhost:5173,http://127.0.0.1:5173",
	))
}

func normalizeOrigins(value string) string {
	origins := strings.Split(value, ",")
	normalized := make([]string, 0, len(origins))

	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		origin = strings.TrimRight(origin, "/")
		if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			origin = "https://" + origin
		}
		normalized = append(normalized, origin)
	}

	return strings.Join(normalized, ",")
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvFrom(keys []string, fallback string) string {
	for _, key := range keys {
		value := os.Getenv(key)
		if value != "" {
			return value
		}
	}

	return fallback
}
