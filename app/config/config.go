package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var (
	IP   = getEnv("IP", "127.0.0.1")
	PORT = getEnv("PORT", "8000")
	DB   = getEnv("DB", "db.sqlite3")
)

func getEnv(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Environment variable not found %v`, name))
}
