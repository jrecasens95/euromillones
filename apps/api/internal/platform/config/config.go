package config

import "os"

type AppConfig struct {
	Port         string
	DrawsDataDir string
	JWTSecret    string
}

var Current AppConfig

func Load() {
	Current = AppConfig{
		Port:         getEnv("PORT", "4000"),
		DrawsDataDir: getEnv("DRAWS_DATA_DIR", "data"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
