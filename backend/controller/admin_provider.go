package controller

import (
	"domain-manager/response"
	"domain-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminProviderController struct {
	providerService *service.ProviderService
}

func NewAdminProviderController(providerService *service.ProviderService) *AdminProviderController {
	return &AdminProviderController{providerService: providerService}
}

// ListProviders 管理员获取所有服务商
func (c *AdminProviderController) ListProviders(ctx *gin.Context) {
	providers, err := c.providerService.ListProviders(ctx)
	if err != nil {
		response.ErrorMessage(ctx, response.InternalError, "获取服务商列表失败")
		return
	}
	response.Success(ctx, providers)
}

// ListEnabledProviders 普通用户获取已启用的服务商
func (c *AdminProviderController) ListEnabledProviders(ctx *gin.Context) {
	providers, err := c.providerService.ListEnabledProviders(ctx)
	if err != nil {
		response.ErrorMessage(ctx, response.InternalError, "获取服务商列表失败")
		return
	}
	response.Success(ctx, providers)
}

// CreateProvider 管理员新增服务商
func (c *AdminProviderController) CreateProvider(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	p, err := c.providerService.CreateProvider(ctx, req.Name, req.Type, req.Description)
	if err != nil {
		response.ErrorMessage(ctx, response.InternalError, err.Error())
		return
	}
	response.Success(ctx, p)
}

// UpdateProviderStatus 管理员启用/禁用服务商
func (c *AdminProviderController) UpdateProviderStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	if err := c.providerService.UpdateProviderStatus(ctx, id, req.Enabled); err != nil {
		response.ErrorMessage(ctx, response.InternalError, "更新服务商状态失败")
		return
	}
	response.Success(ctx, nil)
}

// DeleteProvider 管理员删除服务商
func (c *AdminProviderController) DeleteProvider(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, response.InvalidParams)
		return
	}
	if err := c.providerService.DeleteProvider(ctx, id); err != nil {
		response.ErrorMessage(ctx, response.InternalError, "删除服务商失败")
		return
	}
	response.Success(ctx, nil)
}
