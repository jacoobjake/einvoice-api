package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/internal/services"
	pkgError "github.com/jacoobjake/einvoice-api/pkg/error"
	"github.com/jacoobjake/einvoice-api/pkg/response"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response.JSONApiResponse{
			Success:          false,
			Code:             http.StatusUnprocessableEntity,
			Message:          "invalid request data asd",
			ValidationErrors: pkgError.FormatValidationError(err),
		})
		return
	}

	token, refreshToken, err := h.AuthService.Token(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		log.Println("Error during login:", err.Error())
		c.JSON(http.StatusUnauthorized, response.JSONApiResponse{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, response.JSONApiResponse{
		Success: true,
		Message: "login successful",
		Data: gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetString("auth_token")
	err := h.AuthService.RevokeToken(c, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.JSONApiResponse{
			Success: false,
			Message: "an error occurred while logging out",
		})
	}
	c.JSON(http.StatusOK, response.JSONApiResponse{
		Success: true,
		Message: "logged out successfully",
	})
}

func NewAuthHandler(AuthService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: AuthService,
	}
}
