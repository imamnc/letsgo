package user

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"
	response "letsgo/shared/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
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
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param limit query int false "Maximum results"
// @Param offset query int false "Result offset"
// @Success 200 {object} ListUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/ [get]
func (h *Handler) List(c *fiber.Ctx) error {
	limit := int32(20)
	offset := int32(0)

	if q := c.Query("limit"); q != "" {
		parsed, err := strconv.Atoi(q)
		if err != nil || parsed < 1 {
			return response.BadRequest(c, "limit must be a positive integer")
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = int32(parsed)
	}

	if q := c.Query("offset"); q != "" {
		parsed, err := strconv.Atoi(q)
		if err != nil || parsed < 0 {
			return response.BadRequest(c, "offset must be a non-negative integer")
		}
		offset = int32(parsed)
	}

	users, err := h.service.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return h.Error(c, err)
	}

	responseData := make([]UserResponse, len(users))
	for i, user := range users {
		responseData[i] = User(user)
	}

	return response.Success(c, "Users fetched successfully", ListUsersResponse{
		Users:  responseData,
		Limit:  limit,
		Offset: offset,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateUserRequest true "Create user request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var request CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Name == "" || request.Email == "" || request.Password == "" {
		return response.BadRequest(c, "name, email, and password are required")
	}

	user, err := h.service.CreateUser(c.Context(), request)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Created(c, "User created successfully", User(user))
}

// GetUser godoc
// @Summary Get user details
// @Description Get a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid user id")
	}

	user, err := h.service.FindUserByID(c.Context(), int32(id))
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "User fetched successfully", User(user))
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user fields by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid user id")
	}

	var request UpdateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if request.Name == nil && request.Email == nil && request.Password == nil {
		return response.BadRequest(c, "at least one field is required to update")
	}

	user, err := h.service.UpdateUser(c.Context(), int32(id), request)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "User updated successfully", User(user))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 204 {string} string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid user id")
	}

	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		return h.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) Error(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, ErrEmailAlreadyUsed):
		return response.Conflict(c, err.Error())
	case strings.Contains(err.Error(), "cannot be empty"):
		return response.BadRequest(c, err.Error())
	default:
		return response.InternalServerError(c, "something went wrong")
	}
}

func User(user dbsql.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
