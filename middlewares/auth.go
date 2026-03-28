package middlewares

import (
	"backend/controllers"
	"backend/util"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

var allowedURLs = map[string]bool{
	"/api/v1/auth/login":    true,
	"/api/v1/auth/register": true,
}

func AuthorizeAll(c *fiber.Ctx) error {
	if _, ok := allowedURLs[c.Path()]; ok {
		return c.Next()
	}
	store := util.GetSessionStore()

	// Retrieve the session data
	session, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to get session in Auth middleware"))
	}

	// Check if the user is authenticated
	userInfoJSON := session.Get("user_info")

	if userInfoJSON == nil {
		err := session.Destroy()
		if err != nil {
			return err
		}
		// Redirect to login page
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Unauthorized"))
	}

	userByte := userInfoJSON.([]byte)

	var userInfo controllers.UserInfo
	err = json.Unmarshal(userByte, &userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to unmarshal json"))
	}

	if userInfo.ExpiresOn.Before(time.Now()) {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Session expired"))
	}

	if userInfo.IsNotAdminOrUser() {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusForbidden).Send([]byte("Unauthorized access"))
	}

	// Store the user ID in the context
	c.Locals("user_id", userInfo.UserID)

	// Renew the session
	if err := session.Regenerate(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to renew session"))
	}

	// Set the session cookie
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to save session cookie"))
	}

	// Continue processing the request
	return c.Next()
}

func AuthorizeAdmin(c *fiber.Ctx) error {
	if _, ok := allowedURLs[c.Path()]; ok {
		return c.Next()
	}

	store := util.GetSessionStore()

	// Retrieve the session data
	session, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to get session in Auth middleware"))
	}

	// Check if the user is authenticated
	userInfoJSON := session.Get("user_info")

	if userInfoJSON == nil {
		err := session.Destroy()
		if err != nil {
			return err
		}
		// Redirect to login page
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Unauthorized"))
	}

	userByte := userInfoJSON.([]byte)

	var userInfo controllers.UserInfo
	err = json.Unmarshal(userByte, &userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to unmarshal json"))
	}

	if userInfo.ExpiresOn.Before(time.Now()) {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Session expired"))
	}

	if userInfo.IsNotAdmin() {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusForbidden).Send([]byte("Unauthorized access"))
	}

	// Store the user ID in the context
	c.Locals("user_id", userInfo.UserID)

	// Renew the session
	if err := session.Regenerate(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to renew session"))
	}

	// Set the session cookie
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to save session cookie"))
	}

	// Continue processing the request
	return c.Next()
}

func AuthorizeAdminOrCorporate(c *fiber.Ctx) error {
	if _, ok := allowedURLs[c.Path()]; ok {
		return c.Next()
	}

	store := util.GetSessionStore()

	// Retrieve the session data
	session, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to get session in Auth middleware"))
	}

	// Check if the user is authenticated
	userInfoJSON := session.Get("user_info")

	if userInfoJSON == nil {
		err := session.Destroy()
		if err != nil {
			return err
		}
		// Redirect to login page
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Unauthorized"))
	}

	userByte := userInfoJSON.([]byte)

	var userInfo controllers.UserInfo
	err = json.Unmarshal(userByte, &userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to unmarshal json"))
	}

	if userInfo.ExpiresOn.Before(time.Now()) {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusUnauthorized).Send([]byte("Session expired"))
	}

	if userInfo.IsNotAdminOrCorporate() {
		err := session.Destroy()
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusForbidden).Send([]byte("Unauthorized access"))
	}

	// Store the user ID in the context
	c.Locals("user_id", userInfo.UserID)

	// Renew the session
	if err := session.Regenerate(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to renew session"))
	}

	// Set the session cookie
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to save session cookie"))
	}

	// Continue processing the request
	return c.Next()
}

// func authenticate(c *fiber.Ctx) (controllers.UserInfo, *session.Session, error) {
// 	store := util.GetSessionStore()

// 	// Retrieve the session data
// 	session, err := store.Get(c)
// 	if err != nil {
// 		return controllers.UserInfo{}, nil, c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to get session in Auth middleware"))
// 	}

// 	// Check if the user is authenticated
// 	userInfoJSON := session.Get("user_info")

// 	if userInfoJSON == nil {
// 		session.Destroy()
// 		return controllers.UserInfo{}, nil, c.Status(fiber.StatusUnauthorized).Send([]byte("Unauthorized"))
// 	}

// 	userByte := userInfoJSON.([]byte)

// 	var userInfo controllers.UserInfo
// 	err = json.Unmarshal(userByte, &userInfo)
// 	if err != nil {
// 		return controllers.UserInfo{}, nil, c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to unmarshal json"))
// 	}

// 	return userInfo, session, nil
// }
