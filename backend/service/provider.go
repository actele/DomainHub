package service

import (
	"context"
	"domain-manager/ent"
	"domain-manager/ent/provider"
	"fmt"
	"time"
)

type ProviderService struct {
	client *ent.Client
}

func NewProviderService(client *ent.Client) *ProviderService {
	return &ProviderService{client: client}
}

// EnsureDefaultProviders 初始化内置三家服务商（幂等）
func (s *ProviderService) EnsureDefaultProviders(ctx context.Context) error {
	defaults := []struct {
		Name        string
		Type        string
		Description string
	}{
		{"阿里云", "aliyun", "阿里云 DNS 解析（Alidns）"},
		{"腾讯云", "tencent", "腾讯云 DNSPod"},
		{"Cloudflare", "cloudflare", "Cloudflare DNS"},
	}

	for _, d := range defaults {
		exists, err := s.client.Provider.Query().Where(provider.TypeEQ(d.Type)).Exist(ctx)
		if err != nil {
			return fmt.Errorf("检查服务商失败: %w", err)
		}
		if exists {
			continue
		}
		now := time.Now()
		if _, err := s.client.Provider.Create().
			SetName(d.Name).
			SetType(d.Type).
			SetDescription(d.Description).
			SetEnabled(true).
			SetCreatedAt(now).
			SetUpdatedAt(now).
			Save(ctx); err != nil {
			return fmt.Errorf("创建服务商失败: %w", err)
		}
	}
	return nil
}

// ListProviders 获取所有服务商（管理员用）
func (s *ProviderService) ListProviders(ctx context.Context) ([]*ent.Provider, error) {
	return s.client.Provider.Query().Order(ent.Asc(provider.FieldID)).All(ctx)
}

// ListEnabledProviders 获取已启用的服务商（普通用户用）
func (s *ProviderService) ListEnabledProviders(ctx context.Context) ([]*ent.Provider, error) {
	return s.client.Provider.Query().
		Where(provider.EnabledEQ(true)).
		Order(ent.Asc(provider.FieldID)).
		All(ctx)
}

// UpdateProviderStatus 启用/禁用服务商
func (s *ProviderService) UpdateProviderStatus(ctx context.Context, id int, enabled bool) error {
	return s.client.Provider.UpdateOneID(id).
		SetEnabled(enabled).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// CreateProvider 新增自定义服务商
func (s *ProviderService) CreateProvider(ctx context.Context, name, ptype, description string) (*ent.Provider, error) {
	now := time.Now()
	return s.client.Provider.Create().
		SetName(name).
		SetType(ptype).
		SetDescription(description).
		SetEnabled(true).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Save(ctx)
}

// DeleteProvider 删除服务商
func (s *ProviderService) DeleteProvider(ctx context.Context, id int) error {
	return s.client.Provider.DeleteOneID(id).Exec(ctx)
}
