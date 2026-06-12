package response

import "github.com/gofiber/fiber/v2"

func Success[T any](
	c *fiber.Ctx,
	message string,
	data T,
) error {
	return c.JSON(Response[T]{
		Success: true,
		Status:  "SUCCESS",
		Message: message,
		Data:    data,
	})
}

func Created[T any](
	c *fiber.Ctx,
	message string,
	data T,
) error {
	return c.Status(fiber.StatusCreated).JSON(Response[T]{
		Success: true,
		Status:  "SUCCESS",
		Message: message,
		Data:    data,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
