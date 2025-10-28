package data

import (
	"context"
	"time"

	"ShunFengParcel/config"
	"ShunFengParcel/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}

// RedisSet 保存数据到 Redis，5分钟过期
func (r *greeterRepo) RedisSet(ctx context.Context, key, value string) error {
	return r.data.Rdb.Set(ctx, key, value, time.Minute*5).Err()
}

// RedisGet 从 Redis 获取数据
func (r *greeterRepo) RedisGet(ctx context.Context, key string) (string, error) {
	return r.data.Rdb.Get(ctx, key).Result()
}

// GetUserByPhone 根据手机号查询用户
func (r *greeterRepo) GetUserByPhone(ctx context.Context, phone string) (*biz.Greeter, error) {
	var user config.SfUsers
	err := r.data.DB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &biz.Greeter{Hello: user.Nickname}, nil
}

// CreateUser 创建用户
func (r *greeterRepo) CreateUser(ctx context.Context, phone string) error {
	user := &config.SfUsers{
		Phone:      phone,
		Nickname:   "用户" + phone[7:],
		Status:     "active",
		IsRealname: 0,
	}
	result := r.data.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		r.log.Errorf("创建用户失败: %v", result.Error)
		return result.Error
	}
	r.log.Infof("成功创建用户: %s, ID: %d", phone, user.Id)
	return nil
}
