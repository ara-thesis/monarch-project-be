package helper

import "github.com/gofiber/fiber/v2"

type ResponseHelper struct {
}

func (r *ResponseHelper) Data(c *fiber.Ctx, result []interface{}, msg string, code int) error {

	var success bool

	if code >= 400 {
		success = false
	} else {
		success = true
	}

	return c.Status(code).JSON(fiber.Map{
		"success": success,
		"data":    result,
		"message": msg,
		"code":    code,
	})

}

func (r *ResponseHelper) Success(c *fiber.Ctx, result []interface{}, msg string) error {

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": msg,
		"code":    404,
	})

}

func (r *ResponseHelper) Created(c *fiber.Ctx, msg string) error {

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": msg,
		"code":    404,
	})

}

func (r *ResponseHelper) NotFound(c *fiber.Ctx, msg string) error {

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    404,
	})

}

func (r *ResponseHelper) ServerError(c *fiber.Ctx, msg string) error {

	return c.Status(500).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    500,
	})

}

func (r *ResponseHelper) Unauthorized(c *fiber.Ctx, msg string) error {

	return c.Status(403).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    401,
	})

}

func (r *ResponseHelper) Forbidden(c *fiber.Ctx, msg string) error {

	return c.Status(403).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    403,
	})

}
