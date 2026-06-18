package response

import "github.com/gofiber/fiber/v2"

func Error(
	c *fiber.Ctx,
	code int,
	status string,
	message string,
) error {
	return c.Status(code).JSON(
		Response[any]{
			Status:  status,
			Message: message,
		},
	)
}

func BadRequest(c *fiber.Ctx, message ...string) error {
	msg := "Bad request"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusBadRequest, "BAD_REQUEST", msg)
}

func Forbidden(c *fiber.Ctx, message ...string) error {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusForbidden, "FORBIDDEN", msg)
}

func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", msg)
}

func ValidationError(c *fiber.Ctx, message ...string) error {
	msg := "Validation error"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", msg)
}

func InternalServerError(c *fiber.Ctx, message ...string) error {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", msg)
}

func NotFound(c *fiber.Ctx, message ...string) error {
	msg := "Not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", msg)
}

func Conflict(c *fiber.Ctx, message ...string) error {
	msg := "Conflict"
	if len(message) > 0 {
		msg = message[0]
	}
	return Error(c, fiber.StatusConflict, "CONFLICT", msg)
}
