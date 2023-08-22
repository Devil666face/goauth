package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	IP            = env("IP", "127.0.0.1")
	HTTP_PORT     = env("HTTP_PORT", "8000")
	HTTPS_PORT    = env("HTTPS_PORT", "4443")
	DB            = env("DB", "db.sqlite3")
	ALLOW_HOST    = env("ALLOW_HOST", "localhost")
	TLS           = boolenv(env("TLS", "False"))
	TLS_KEY       = env("TLS_KEY", "server.key")
	TLS_CRT       = env("TLC_CRT", "server.crt")
	CONNECT_HTTP  = fmt.Sprintf("%v:%v", IP, HTTP_PORT)
	CONNECT_HTTPS = fmt.Sprintf("%v:%v", IP, HTTPS_PORT)
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

func boolenv(env string) bool {
	if env == "True" || env == "true" {
		return true
	}
	return false
}

func GetSuperuser() (string, string) {
	var (
		USER = env("SUUSER", "superuser")
		PASS = env("SUPASS", "Qwerty123!")
	)
	return USER, PASS
}
