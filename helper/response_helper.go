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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": msg,
		"code":    fiber.StatusOK,
	})

}

func (r *ResponseHelper) Created(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusCreated,
	})

}

func (r *ResponseHelper) BadRequest(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusBadRequest,
	})

}

func (r *ResponseHelper) NotFound(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusNotFound,
	})

}

func (r *ResponseHelper) ServerError(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusInternalServerError,
	})

}

func (r *ResponseHelper) Unauthorized(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusUnauthorized,
	})

}

func (r *ResponseHelper) Forbidden(c *fiber.Ctx, msg string) error {

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"message": msg,
		"code":    fiber.StatusForbidden,
	})

}
