package service

import (
	"context"
	"domain-manager/ent"
	"domain-manager/ent/user"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	client *ent.Client
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{client: client}
}

// EnsureUser 在用户不存在时创建默认管理员用户。
func (s *UserService) EnsureUser(ctx context.Context, username, password string) error {
	u, err := s.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("检查默认用户失败: %w", err)
	}

	if ent.IsNotFound(err) {
		// 创建 admin 用户
		if _, err := s.createUserWithRole(ctx, username, password, "admin"); err != nil {
			return fmt.Errorf("创建默认用户失败: %w", err)
		}
		return nil
	}

	// 如果已存在但不是 admin，升级为 admin
	if u.Role != "admin" {
		if err := s.client.User.UpdateOne(u).SetRole("admin").Exec(ctx); err != nil {
			return fmt.Errorf("升级默认用户为管理员失败: %w", err)
		}
	}
	return nil
}

// CreateUser 创建普通用户（注册接口使用）
func (s *UserService) CreateUser(ctx context.Context, username, password string) (*ent.User, error) {
	return s.createUserWithRole(ctx, username, password, "user")
}

// createUserWithRole 内部创建用户
func (s *UserService) createUserWithRole(ctx context.Context, username, password, role string) (*ent.User, error) {
	exists, err := s.client.User.Query().Where(user.UsernameEQ(username)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已被注册")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return s.client.User.Create().
		SetUsername(username).
		SetPassword(string(hashedPassword)).
		SetRole(role).
		SetStatus("active").
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Save(ctx)
}

// VerifyUser 验证用户登录
func (s *UserService) VerifyUser(ctx context.Context, username, password string) (*ent.User, error) {
	u, err := s.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if u.Status == "disabled" {
		return nil, errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return u, nil
}

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	u, err := s.client.User.Get(ctx, userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.client.User.UpdateOne(u).
		SetPassword(string(hashedPassword)).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// ListUsers 管理员获取所有用户列表
func (s *UserService) ListUsers(ctx context.Context) ([]*ent.User, error) {
	return s.client.User.Query().Order(ent.Asc(user.FieldID)).All(ctx)
}

// AdminCreateUser 管理员创建用户
func (s *UserService) AdminCreateUser(ctx context.Context, username, password, role string) (*ent.User, error) {
	if role != "admin" && role != "user" {
		role = "user"
	}
	return s.createUserWithRole(ctx, username, password, role)
}

// UpdateUserStatus 更新用户状态（active/disabled）
func (s *UserService) UpdateUserStatus(ctx context.Context, userID int, status string) error {
	return s.client.User.UpdateOneID(userID).
		SetStatus(status).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// UpdateUserRole 更新用户角色（admin/user）
func (s *UserService) UpdateUserRole(ctx context.Context, userID int, role string) error {
	return s.client.User.UpdateOneID(userID).
		SetRole(role).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// AdminResetPassword 管理员重置用户密码
func (s *UserService) AdminResetPassword(ctx context.Context, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.client.User.UpdateOneID(userID).
		SetPassword(string(hashedPassword)).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}
