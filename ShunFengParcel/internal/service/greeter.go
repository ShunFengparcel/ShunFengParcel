package service

import (
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/biz"
	"ShunFengParcel/internal/data"
	"ShunFengParcel/pkg"
	"context"
	"math/rand"
	"strconv"

	"gorm.io/gorm"
)

var DB *gorm.DB

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, d *data.Data) *GreeterService {
	DB = d.DB // 从 Data 结构体中获取 DB 实例
	return &GreeterService{uc: uc}
}

// SendSms implements helloworld.GreeterServer.
func (s *GreeterService) SendSms(ctx context.Context, in *v1.SendSmsRequest) (*v1.SendSmsReply, error) {
	// 生成验证码
	i := rand.Intn(9000) + 1000

	// 保存验证码到 Redis
	err := s.uc.SetSendSms(ctx, in.Phone, strconv.Itoa(i))
	if err != nil {
		return &v1.SendSmsReply{
			Code:    500,
			Message: "验证码保存失败",
			Data:    "",
		}, nil
	}

	// 发送短信
	sms := pkg.SendSms(in.Phone, strconv.Itoa(i))

	return &v1.SendSmsReply{
		Code:    200,
		Message: "短信发送成功",
		Data:    sms,
	}, nil
}

// UserLoginSendSms implements helloworld.GreeterServer.
// 验证码登录：如果用户不存在则自动注册
func (s *GreeterService) UserLoginSendSms(ctx context.Context, in *v1.UserLoginSendSmsRequest) (*v1.UserLoginSendSmsReply, error) {
	// 1. 验证验证码
	savedCode := s.uc.GetSendSmd(ctx, in.Phone)
	if savedCode == "" {
		return &v1.UserLoginSendSmsReply{
			Code:    400,
			Message: "验证码已过期，请重新获取",
			Data:    "",
		}, nil
	}

	if savedCode != in.SendSms {
		return &v1.UserLoginSendSmsReply{
			Code:    400,
			Message: "验证码错误",
			Data:    "",
		}, nil
	}

	// 2. 检查用户是否存在
	user, err := s.uc.GetUserByPhone(ctx, in.Phone)

	if err != nil {
		// 用户不存在，创建新用户（注册）
		createErr := s.uc.CreateUser(ctx, in.Phone)
		if createErr != nil {
			return &v1.UserLoginSendSmsReply{
				Code:    500,
				Message: "注册失败: " + createErr.Error(),
				Data:    "",
			}, nil
		}

		// 重新获取用户信息
		user, err = s.uc.GetUserByPhone(ctx, in.Phone)
		if err != nil {
			return &v1.UserLoginSendSmsReply{
				Code:    500,
				Message: "获取用户信息失败",
				Data:    "",
			}, nil
		}

		// 新用户注册成功
		return &v1.UserLoginSendSmsReply{
			Code:    200,
			Message: "注册成功",
			Data:    "欢迎加入, " + user.Hello,
		}, nil
	}

	// 3. 老用户登录成功
	return &v1.UserLoginSendSmsReply{
		Code:    200,
		Message: "登录成功",
		Data:    "欢迎回来, " + user.Hello,
	}, nil
}

// WxLogin implements helloworld.GreeterServer.
// 微信小程序登录
func (s *GreeterService) WxLogin(ctx context.Context, in *v1.WxLoginRequest) (*v1.WxLoginReply, error) {
	// 调用 wxlogin.go 中的实现
	resp := s.doWxLogin(ctx, in.GetCode(), in.GetNickname(), in.GetAvatarUrl())

	return &v1.WxLoginReply{
		Code:     int64(resp.Code),
		Message:  resp.Message,
		Token:    resp.Token,
		Openid:   resp.OpenId,
		Nickname: resp.NickName,
		UserId:   resp.UserId,
	}, nil
}
