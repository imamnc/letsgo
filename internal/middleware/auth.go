package middleware

import (
	"strings"

	"letsgo/internal/app"

	"github.com/gofiber/fiber/v2"
)

func Auth(application *app.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := strings.TrimSpace(c.Get("Authorization"))
		if authorization == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "UNAUTHORIZED",
				"message": "authorization token is required",
			})
		}

		token := bearerToken(authorization)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "UNAUTHORIZED",
				"message": "authorization token is required",
			})
		}

		claims, err := application.Jwt.Decode(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "UNAUTHORIZED",
				"message": "invalid authorization token",
			})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func bearerToken(authorization string) string {
	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return ""
	}

	if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return strings.TrimSpace(authorization[7:])
	}

	return authorization
}
