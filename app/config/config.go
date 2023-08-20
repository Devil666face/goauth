package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	IP      = env("IP", "127.0.0.1")
	PORT    = env("PORT", "8000")
	DB      = env("DB", "db.sqlite3")
	ALLOWED = env("ALLOWED", "localhost")
)

func env(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	if fallback != "" {
		return fallback
	}
	panic(fmt.Sprintf(`Environment variable not found %v`, name))
}

func GetSuperuser() (string, string) {
	var (
		USER = env("SUUSER", "superuser")
		PASS = env("SUPASS", "Qwerty123!")
	)
	return USER, PASS
}
