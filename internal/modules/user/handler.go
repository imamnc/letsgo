package user

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetUserResponse struct {
	User UserResponse `json:"user"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListUsers godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Maximum results"
// @Param offset query int false "Result offset"
// @Success 200 {object} ListUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/ [get]
func (h *Handler) List(c *fiber.Ctx) error {
	limit := int32(20)
	offset := int32(0)

	if q := c.Query("limit"); q != "" {
		parsed, err := strconv.Atoi(q)
		if err != nil || parsed < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "limit must be a positive integer"})
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = int32(parsed)
	}

	if q := c.Query("offset"); q != "" {
		parsed, err := strconv.Atoi(q)
		if err != nil || parsed < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "offset must be a non-negative integer"})
		}
		offset = int32(parsed)
	}

	users, err := h.service.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return h.respondError(c, err)
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = toUserResponse(user)
	}

	return c.JSON(ListUsersResponse{
		Users:  response,
		Limit:  limit,
		Offset: offset,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var request CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid request body"})
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Name == "" || request.Email == "" || request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "name, email, and password are required"})
	}

	user, err := h.service.CreateUser(c.Context(), request)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toUserResponse(user))
}

// GetUser godoc
// @Summary Get user details
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid user id"})
	}

	user, err := h.service.FindUserByID(c.Context(), int32(id))
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(toUserResponse(user))
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user fields by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid user id"})
	}

	var request UpdateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid request body"})
	}

	if request.Name == nil && request.Email == nil && request.Password == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "at least one field is required to update"})
	}

	user, err := h.service.UpdateUser(c.Context(), int32(id), request)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(toUserResponse(user))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {string} string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid user id"})
	}

	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		return h.respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) respondError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "NOT_FOUND", "message": err.Error()})
	case errors.Is(err, ErrEmailAlreadyUsed):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "CONFLICT", "message": err.Error()})
	case strings.Contains(err.Error(), "cannot be empty"):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "ERROR", "message": "something went wrong"})
	}
}

func toUserResponse(user dbsql.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
