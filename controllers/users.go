package controllers

import (
	"backend/models"
	"backend/repos"
	"backend/util"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController interface {
	GetMe(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	// PatchUser(c *fiber.Ctx) error
	// DeleteUser(c *fiber.Ctx) error
}

type userController struct {
	userRepo repos.IUserRepository
}

func NewUserController(userRepo repos.IUserRepository) UserController {
	return &userController{userRepo}
}

func (u *userController) GetMe(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store := util.GetSessionStore()
	session, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte("Failed to get session in Auth middleware"))
	}

	userID := c.Locals("user_id")
	if userID == nil {
		session.Destroy()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	fmt.Println("localsId: ", userID)

	userIDToUUID := uuid.MustParse(userID.(string))

	// TODO: exclude password
	user, err := u.userRepo.GetMe(ctx, userIDToUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if user == (models.User{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		//"request_id":
		"message": "User found",
		"data":    user,
	})
}

func (u *userController) GetUser(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Please specify a valid user id",
		})
	}

	idToUuid := uuid.MustParse(id)

	user, err := u.userRepo.GetUserByID(ctx, idToUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		//"request_id":
		"message": "User found",
		"data":    user,
	})
}

func (u *userController) GetUsers(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := u.userRepo.GetUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Users not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users found",
		"data":    users,
	})
}
