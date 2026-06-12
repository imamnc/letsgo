package health

import (
	response "letsgo/shared/response"

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
		return response.InternalServerError(c, "service unavailable")
	}

	count, err := h.service.GetUserCount(c)
	if err != nil {
		return response.InternalServerError(c, "service unavailable")
	}

	return response.Success(c, "Service is up", map[string]any{
		"project":      "LetsGO",
		"project_desc": "API Boilerplate built with Go and Fiber. Just clone and LetsGO!",
		"version":      "1.0.0",
		"status":       "UP",
		"user_count":   count,
	})

}
