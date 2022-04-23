package module

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type AccountHandler struct {
}

func (u *AccountHandler) GetUserInfo(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) CreateUser(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) UserLogin(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) EditUser(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) EditUserAsAdmin(c *fiber.Ctx) error {
	return c.SendString("Test")
}

func (u *AccountHandler) DeleteUser(c *fiber.Ctx) error {
	return c.SendString("Test")
}
