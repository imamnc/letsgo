package auth

import (
	"errors"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"
	sharedjwt "letsgo/shared/jwt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FetchUserResponse struct {
	User UserResponse `json:"user"`
}

type AccessTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// AccessToken godoc
// @Summary Authenticate a user and receive access/refresh tokens
// @Description Authenticate using email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AccessTokenRequest true "Credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/access-token [post]
func (h *Handler) AccessToken(c *fiber.Ctx) error {
	var request AccessTokenRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid request body"})
	}

	if strings.TrimSpace(request.Email) == "" || strings.TrimSpace(request.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "email and password are required"})
	}

	user, accessToken, refreshToken, err := h.service.Authenticate(c.Context(), request.Email, request.Password)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	})
}

// RefreshToken godoc
// @Summary Refresh a user's access token
// @Description Exchange a refresh token for a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh-token [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var request RefreshTokenRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "invalid request body"})
	}

	tokenString := strings.TrimSpace(request.RefreshToken)
	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "refresh_token is required"})
	}

	user, accessToken, refreshToken, err := h.service.Refresh(c.Context(), tokenString)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	})
}

// FetchUser godoc
// @Summary Get current authenticated user
// @Description Fetch logged-in user details from access token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} FetchUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/user [get]
func (h *Handler) FetchUser(c *fiber.Ctx) error {
	tokenString := bearerToken(c.Get("Authorization"))
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "UNAUTHORIZED", "message": "authorization token is required"})
	}

	user, err := h.service.FetchUser(c.Context(), tokenString)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(FetchUserResponse{User: toUserResponse(user)})
}

func (h *Handler) respondError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "UNAUTHORIZED", "message": err.Error()})
	case errors.Is(err, sharedjwt.ErrInvalidToken):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "UNAUTHORIZED", "message": err.Error()})
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "NOT_FOUND", "message": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "ERROR", "message": "something went wrong"})
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

func toUserResponse(user dbsql.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      defaultUserRole,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
