package service

import (
	pb "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/config"
	"ShunFengParcel/inits"
	"ShunFengParcel/pkg"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) AddOrder(ctx context.Context, req *pb.AddOrderReq) (*pb.AddOrderResp, error) {
	o := config.SfOrders{
		OrderNo:         uuid.NewString(),
		UserId:          req.Userid,
		CourierId:       req.Courierid,
		SenderName:      req.Sendername,
		SenderPhone:     req.Senderphone,
		SenderAddress:   req.Senderaddress,
		ReceiverName:    req.Receivername,
		ReceiverPhone:   req.Receiverphone,
		ReceiverAddress: req.Senderaddress,
	}
	inits.DB.Create(&o)
	rabbitmq := pkg.NewRabbitMQSimple("" +
		"orders")
	rabbitmq.PublishSimple(req.Sendername + req.Receivername)
	fmt.Println("发送成功！")
	return &pb.AddOrderResp{
		Id: o.Id,
	}, nil
}
