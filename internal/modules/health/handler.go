package health

import (
	"github.com/gofiber/fiber/v2"
)

// Handler struct to hold application context
type Handler struct {
	service *Service
}

// NewHandler initializes a new Handler with the application context
func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

// Check is the handler function for the health check endpoint
func (h *Handler) Check(
	c *fiber.Ctx,
) error {
	if !h.service.Check(c) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "DOWN",
		})
	}

	count, err := h.service.GetUserCount(c)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "DOWN",
		})
	}

	return c.JSON(fiber.Map{
		"project":      "LetsGO",
		"project_desc": "API Boilerplate built with Go, and Fiber. Just clone and LetsGO!",
		"version":      "1.0.0",
		"status":       "UP",
		"user_count":   count,
	})

}
