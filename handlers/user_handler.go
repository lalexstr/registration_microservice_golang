package handlers

import (
	"net/http"
	"strconv"

	"auth-service/middleware"
	"auth-service/models"
	"auth-service/repositories"
	"auth-service/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repositories.UserRepo
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		repo: repositories.NewUserRepo(),
	}
}

// Create user by admin (create different roles)
type createUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
	Role     string `json:"role"` // admin or user
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	// only admin allowed — middleware should ensure it
	var body createUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if body.Role != "admin" {
		body.Role = "user"
	}
	// check exists
	if _, err := h.repo.FindByEmail(body.Email); err == nil {
		utils.JSONError(c, http.StatusBadRequest, "email already used")
		return
	}
	u := &models.User{
		Email:    body.Email,
		FullName: body.FullName,
		Role:     body.Role,
	}
	if err := u.SetPassword(body.Password); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}
	if err := h.repo.Create(u); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.Password = ""
	utils.JSONSuccess(c, http.StatusCreated, u)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	uCtx := c.MustGet(middleware.CtxUserKey).(*middleware.ContextUser)
	user, err := h.repo.FindByID(uCtx.ID)
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	user.Password = ""
	utils.JSONSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	requester := c.MustGet(middleware.CtxUserKey).(*middleware.ContextUser)
	// if requester is not admin and not requesting own -> forbidden
	if requester.Role != "admin" && requester.ID != uint(id64) {
		utils.JSONError(c, http.StatusForbidden, "forbidden")
		return
	}
	u, err := h.repo.FindByID(uint(id64))
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	u.Password = ""
	utils.JSONSuccess(c, http.StatusOK, u)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// only admin
	// query params: page, size
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if pv, err := strconv.Atoi(p); err == nil && pv > 0 {
			page = pv
		}
	}
	if s := c.Query("size"); s != "" {
		if sv, err := strconv.Atoi(s); err == nil && sv > 0 && sv <= 100 {
			size = sv
		}
	}
	offset := (page - 1) * size
	users, total, err := h.repo.List(offset, size)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	// sanitize
	for i := range users {
		users[i].Password = ""
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{
		"items": users,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

type updateUserReq struct {
	FullName *string `json:"full_name"`
	Password *string `json:"password"`
	Role     *string `json:"role"` // admin only
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	requester := c.MustGet(middleware.CtxUserKey).(*middleware.ContextUser)
	if requester.Role != "admin" && requester.ID != uint(id64) {
		utils.JSONError(c, http.StatusForbidden, "forbidden")
		return
	}
	u, err := h.repo.FindByID(uint(id64))
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	var body updateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if body.FullName != nil {
		u.FullName = *body.FullName
	}
	if body.Password != nil {
		if err := u.SetPassword(*body.Password); err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
			return
		}
	}
	if body.Role != nil {
		if requester.Role != "admin" {
			utils.JSONError(c, http.StatusForbidden, "only admin can change role")
			return
		}
		if *body.Role != "admin" && *body.Role != "user" {
			utils.JSONError(c, http.StatusBadRequest, "invalid role")
			return
		}
		u.Role = *body.Role
	}
	if err := h.repo.Update(u); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.Password = ""
	utils.JSONSuccess(c, http.StatusOK, u)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	requester := c.MustGet(middleware.CtxUserKey).(*middleware.ContextUser)
	if requester.Role != "admin" && requester.ID != uint(id64) {
		utils.JSONError(c, http.StatusForbidden, "forbidden")
		return
	}
	u, err := h.repo.FindByID(uint(id64))
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	if err := h.repo.Delete(u); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{"deleted": true})
}
