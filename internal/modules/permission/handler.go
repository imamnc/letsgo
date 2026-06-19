package permission

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"
	"letsgo/shared/format"
	response "letsgo/shared/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PermissionResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	ParentID    *int32    `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListPermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Limit       int32                `json:"limit"`
	Offset      int32                `json:"offset"`
}

type CreatePermissionRequest struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
	ParentID    *int32  `json:"parent_id,omitempty"`
}

type UpdatePermissionRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *int32  `json:"parent_id,omitempty"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListPermissions godoc
// @Summary List permissions
// @Description Get a paginated list of permissions
// @Tags Permissions
// @Accept json
// @Produce json
// @Param limit query int false "Maximum results"
// @Param offset query int false "Result offset"
// @Success 200 {object} ListPermissionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /permissions/ [get]
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

	permissions, err := h.service.ListPermissions(c.Context(), limit, offset)
	if err != nil {
		return h.Error(c, err)
	}

	responseData := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responseData[i] = Permission(permission)
	}

	return response.Success(c, "Permissions fetched successfully", ListPermissionsResponse{
		Permissions: responseData,
		Limit:       limit,
		Offset:      offset,
	})
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a permission with name, code, optional description, and optional parent permission
// @Tags Permissions
// @Accept json
// @Produce json
// @Param request body CreatePermissionRequest true "Create permission request"
// @Success 201 {object} PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /permissions/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var request CreatePermissionRequest
	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.TrimSpace(request.Code)

	if request.Name == "" || request.Code == "" {
		return response.BadRequest(c, "name and code are required")
	}

	permission, err := h.service.CreatePermission(c.Context(), request)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Created(c, "Permission created successfully", Permission(permission))
}

// GetPermission godoc
// @Summary Get permission details
// @Description Get a permission by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /permissions/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid permission id")
	}

	permission, err := h.service.FindPermissionByID(c.Context(), int32(id))
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "Permission fetched successfully", Permission(permission))
}

// UpdatePermission godoc
// @Summary Update a permission
// @Description Update permission fields by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param request body UpdatePermissionRequest true "Update permission request"
// @Success 200 {object} PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /permissions/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid permission id")
	}

	var request UpdatePermissionRequest
	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if request.Name == nil && request.Code == nil && request.Description == nil && request.ParentID == nil {
		return response.BadRequest(c, "at least one field is required to update")
	}

	permission, err := h.service.UpdatePermission(c.Context(), int32(id), request)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "Permission updated successfully", Permission(permission))
}

// DeletePermission godoc
// @Summary Delete a permission
// @Description Delete a permission by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 204 {string} string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /permissions/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid permission id")
	}

	if err := h.service.DeletePermission(c.Context(), int32(id)); err != nil {
		return h.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) Error(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrPermissionNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, ErrPermissionAlreadyUsed):
		return response.Conflict(c, err.Error())
	case strings.Contains(err.Error(), "cannot be empty"):
		return response.BadRequest(c, err.Error())
	default:
		return response.InternalServerError(c, "something went wrong")
	}
}

func Permission(permission dbsql.Permission) PermissionResponse {
	return PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Code:        permission.Code,
		Description: format.FromNullString(permission.Description),
		ParentID:    format.FromNullInt32(permission.ParentID),
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}
