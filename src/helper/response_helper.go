package helper

import "github.com/gofiber/fiber/v2"

type ResponseHelper struct {
}

func (r *ResponseHelper) Success(c *fiber.Ctx, result []interface{}, msg string, code int) error {

	return c.Status(code).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": msg,
		"code":    code,
	})

}

func (r *ResponseHelper) Failed(c *fiber.Ctx, msg string, code int) error {

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    code,
	})

}
