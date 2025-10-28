package biz

import (
	"context"

	v1 "ShunFengParcel/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
	RedisGet(context.Context, string) (string, error)
	RedisSet(context.Context, string, string) error
	GetUserByPhone(context.Context, string) (*Greeter, error)
	CreateUser(context.Context, string) error
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

// SetSendSms 保存验证码到 Redis
func (uc *GreeterUsecase) SetSendSms(ctx context.Context, phone, code string) error {
	key := "verify_code:" + phone
	return uc.repo.RedisSet(ctx, key, code)
}

// GetSendSmd 从 Redis 获取验证码
func (uc *GreeterUsecase) GetSendSmd(ctx context.Context, phone string) string {
	key := "verify_code:" + phone
	val, _ := uc.repo.RedisGet(ctx, key)
	return val
}

// GetUserByPhone 根据手机号获取用户
func (uc *GreeterUsecase) GetUserByPhone(ctx context.Context, phone string) (*Greeter, error) {
	return uc.repo.GetUserByPhone(ctx, phone)
}

// CreateUser 创建新用户
func (uc *GreeterUsecase) CreateUser(ctx context.Context, phone string) error {
	return uc.repo.CreateUser(ctx, phone)
}

// SetToken 保存用户 token 到 Redis
func (uc *GreeterUsecase) SetToken(ctx context.Context, token string, userID int64) error {
	key := "user_token:" + token
	return uc.repo.RedisSet(ctx, key, string(rune(userID)))
}

// GetUserByToken 通过 token 获取用户ID
func (uc *GreeterUsecase) GetUserByToken(ctx context.Context, token string) (string, error) {
	key := "user_token:" + token
	return uc.repo.RedisGet(ctx, key)
}
