package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/redis"
)

func GetSessionStore() *session.Store {
	storage := redis.New(redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "redis",
		Database: 0,
	})

	store := session.New(session.Config{
		Storage:        storage,
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
		Expiration:     12 * time.Hour,
	})

	return store
}

func DelSession(c *fiber.Ctx, sessionID string) {
	storage := redis.New(redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "redis",
		Database: 0,
	})

	storage.Delete(sessionID)
}
