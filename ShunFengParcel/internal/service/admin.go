package service

import (
	pb "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/config"
	"ShunFengParcel/inits"
	"ShunFengParcel/pkg"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type AdminService struct {
	pb.UnimplementedAdminServer
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) AdminRegister(ctx context.Context, req *pb.AdminRegisterReq) (*pb.AdminRegisterResp, error) {
	a := config.SysAdmin{
		Username: req.Username,
		Password: pkg.Md5(req.Password),
	}
	inits.DB.Create(&a)
	return &pb.AdminRegisterResp{
		Id: a.Id,
	}, nil
}
func (s *AdminService) FindAdminByUsername(ctx context.Context, req *pb.FindAdminByUsernameReq) (*pb.FindAdminByUsernameResp, error) {
	// 在方法开始时立即创建 Span
	ctx, span := otel.Tracer("admin-service").Start(ctx, "FindAdminByUsername")
	defer span.End() // 确保 Span 结束时被上报

	// 添加 Span 属性，记录关键信息
	span.SetAttributes(
		attribute.String("username", req.Username),
		attribute.String("operation", "find_admin_by_username"),
	)

	// 打印日志，确认埋点逻辑执行
	log.Printf("触发 Trace 埋点，用户名: %s", req.Username)

	var a config.SysAdmin
	inits.DB.Where("username = ?", req.Username).Find(&a)
	if a.Password != pkg.Md5(req.Password) {
		span.AddEvent("密码验证失败")
		return nil, nil
	}

	span.AddEvent("用户查询成功")
	return &pb.FindAdminByUsernameResp{
		Id: a.Id,
	}, nil
}
func (s *AdminService) UpdateAdminPersonal(ctx context.Context, req *pb.UpdateAdminPersonalReq) (*pb.UpdateAdminPersonalResp, error) {
	a := config.SysAdmin{
		Id:       req.Id,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
	}
	inits.DB.Updates(&a)
	return &pb.UpdateAdminPersonalResp{}, nil

}
func (s *AdminService) AdminStatus(ctx context.Context, req *pb.AdminStatusReq) (*pb.AdminStatusResp, error) {
	a := config.SysAdmin{
		Username: req.Username,
	}
	inits.DB.Where("username = ?", req.Username).Updates(&a)
	if a.Status == 0 {
		return nil, nil
	}
	return &pb.AdminStatusResp{
		Status: int64(a.Status),
	}, nil

}
func (s *AdminService) UpdateAdminStatus(ctx context.Context, req *pb.UpdateAdminStatusReq) (*pb.UpdateAdminStatusResp, error) {
	a := config.SysAdmin{
		Id:     req.Id,
		Status: int8(req.Status),
	}
	inits.DB.Updates(&a)
	return &pb.UpdateAdminStatusResp{}, nil

}
