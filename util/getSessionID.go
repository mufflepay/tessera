package util

import (
	"github.com/gofiber/fiber/v2"
)

func GetSessionIDFromCookie(c *fiber.Ctx) string {
	cookie := c.Cookies("session_id")

	return cookie
}
