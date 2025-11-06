package handlers

import (
	"net/http"

	"auth-service/services"
	"auth-service/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{svc: services.NewAuthService()}
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body registerReq
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := h.svc.Register(body.Email, body.Password, body.FullName)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	// do not return password
	u.Password = ""
	utils.JSONSuccess(c, http.StatusCreated, u)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body loginReq
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := h.svc.Authenticate(body.Email, body.Password)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := h.svc.GenerateJWT(u)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{"token": token})
}
