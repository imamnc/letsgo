package auth

import (
	"errors"
	"strings"
	"time"

	dbsql "letsgo/db/sqlc"
	sharedjwt "letsgo/shared/jwt"
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
// @Tags Auth
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
		return response.BadRequest(c, "invalid request body")
	}

	if strings.TrimSpace(request.Email) == "" || strings.TrimSpace(request.Password) == "" {
		return response.BadRequest(c, "email and password are required")
	}

	user, accessToken, refreshToken, err := h.service.Authenticate(c.Context(), request.Email, request.Password)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "Authentication successful", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         User(user),
	})
}

// RefreshToken godoc
// @Summary Refresh a user's access token
// @Description Exchange a refresh token for a new access token
// @Tags Auth
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
		return response.BadRequest(c, "invalid request body")
	}

	tokenString := strings.TrimSpace(request.RefreshToken)
	if tokenString == "" {
		return response.BadRequest(c, "refresh_token is required")
	}

	user, accessToken, refreshToken, err := h.service.Refresh(c.Context(), tokenString)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "Token refreshed successfully", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         User(user),
	})
}

// FetchUser godoc
// @Summary Get current authenticated user
// @Description Fetch logged-in user details from access token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} FetchUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /auth/user [get]
func (h *Handler) FetchUser(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*sharedjwt.Claims)
	if !ok || claims == nil {
		return response.Unauthorized(c, "authorization token is required")
	}

	user, err := h.service.FindUserByID(c.Context(), claims.UserID)
	if err != nil {
		return h.Error(c, err)
	}

	return response.Success(c, "User fetched successfully", User(user))
}

func (h *Handler) Error(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return response.Unauthorized(c, err.Error())
	case errors.Is(err, sharedjwt.ErrInvalidToken):
		return response.Unauthorized(c, err.Error())
	case errors.Is(err, ErrUserNotFound):
		return response.NotFound(c, err.Error())
	default:
		return response.InternalServerError(c, "something went wrong")
	}
}

func User(user dbsql.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      defaultUserRole,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
