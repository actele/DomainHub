package controller

import (
	"domain-manager/response"
	"domain-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	userService *service.UserService
}

func NewAdminUserController(userService *service.UserService) *AdminUserController {
	return &AdminUserController{userService: userService}
}

// ListUsers 获取所有用户列表
func (c *AdminUserController) ListUsers(ctx *gin.Context) {
	users, err := c.userService.ListUsers(ctx)
	if err != nil {
		response.ErrorMessage(ctx, response.InternalError, "获取用户列表失败")
		return
	}
	// 隐藏密码字段
	type safeUser struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
	}
	result := make([]safeUser, 0, len(users))
	for _, u := range users {
		result = append(result, safeUser{
			ID:        u.ID,
			Username:  u.Username,
			Role:      u.Role,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.Success(ctx, result)
}

// CreateUser 管理员创建用户
func (c *AdminUserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}

	user, err := c.userService.AdminCreateUser(ctx, req.Username, req.Password, req.Role)
	if err != nil {
		response.ErrorMessage(ctx, response.InternalError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"id": user.ID, "username": user.Username, "role": user.Role})
}

// UpdateUserStatus 更新用户状态
func (c *AdminUserController) UpdateUserStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	if req.Status != "active" && req.Status != "disabled" {
		response.ErrorMessage(ctx, response.InvalidParams, "状态值无效")
		return
	}
	if err := c.userService.UpdateUserStatus(ctx, id, req.Status); err != nil {
		response.ErrorMessage(ctx, response.InternalError, "更新用户状态失败")
		return
	}
	response.Success(ctx, nil)
}

// UpdateUserRole 更新用户角色
func (c *AdminUserController) UpdateUserRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	if req.Role != "admin" && req.Role != "user" {
		response.ErrorMessage(ctx, response.InvalidParams, "角色值无效")
		return
	}
	if err := c.userService.UpdateUserRole(ctx, id, req.Role); err != nil {
		response.ErrorMessage(ctx, response.InternalError, "更新用户角色失败")
		return
	}
	response.Success(ctx, nil)
}

// ResetPassword 管理员重置用户密码
func (c *AdminUserController) ResetPassword(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	if err := c.userService.AdminResetPassword(ctx, id, req.Password); err != nil {
		response.ErrorMessage(ctx, response.InternalError, "重置密码失败")
		return
	}
	response.Success(ctx, nil)
}
