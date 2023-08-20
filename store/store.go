package store

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var (
	Store    *session.Store
	Storage  *sqlite3.Storage
	AUTH_KEY string = "authenticated"
	USER_ID  string = "user_id"
)

func SetStore() {
	Storage = sqlite3.New(sqlite3.Config{
		Database: "./db.sqlite3",
		Table:    "storage",
		Reset:    false,
	})
	Store = session.New(session.Config{
		CookieHTTPOnly: true,
		// CookieSecure: true,
		Expiration: time.Hour * 5,
		Storage:    Storage,
	})
}
