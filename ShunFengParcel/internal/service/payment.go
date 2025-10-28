package service

import (
	pb "ShunFengParcel/api/helloworld/payment"
	"ShunFengParcel/config"
	"ShunFengParcel/inits"
	"ShunFengParcel/utils"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/smartwalle/alipay/v3"
	logs "log"
	"time"
)

type PaymentService struct {
	pb.UnimplementedPaymentServer
	log *log.Helper
	// 支付宝客户端实例（暂时注释，避免未使用警告）
	alipayClient *alipay.Client
}

func NewPaymentService(logger log.Logger) *PaymentService {
	// 初始化支付宝客户端
	privateKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5+awnsflJyb+3KNUjgCT0rMZ7Q8vutr9A4Ex10qwv0hOaoiZP8PpAhe/mkf+SyonaymtMZiybiZ3oaXAZxvlIvNq6+O7uX/MLE+87/8UzI2aQqdblkZyym/ehbszPVMg/YjHfDvjq/RTyBYTrLva5gWTdlB7d6FPuBLLJP8ezEe+DSXzYbX0RJAiQHxs4Uf6u4cYQitv63zz1WaM0Jwhk8OMz+xf6+6Yrbc3bZ/y7MzcLyuTWL6t2oFxiQdhisFVv90jTfn2pwsLWCMvjEQrwoH+z1NH3bfdu7YO2PPSGTycuelgqpCIDDTMB8ykll0WHSDZVUunC501TetDpUP5PAgMBAAECggEBAKze52byqY4safggjZb1S+TgIZN9HrijyL3OgWRYU1QUS8LGjSRsFMMbfTdvaQkTGfd1hy26ICboUg8yy1A6w4gYfJ2mPbr5F5upiQzVoIH7myBHr4kvxF/GWPYHX3wuDAHXDhKRAVK1s92QwdA3lII1RAQv5k1R44pxdZtanQ/WuJHNTj4pktme+M7ULgNJmWXsflxNIjkc8tuzF9p6sgrkP6w8lSxR/UfJMHK1X1s2U7I5GgnPRGU7e765usPiA3CvtepDSDThQVhuRnbz3rUOxthRwos2SS6wr2pnqq1nj0K6G9aCYQyazoNjfPtnLpNWbI0pBpDm8rFJbOC05nECgYEA9lEz2GA2vPAFlQApt7UrjznfpfuIef5vQbUUrHaa7svUx0KhxXLkYSMmTF/DQFgVd/kXEsg3vAPdx/ZSsnV/HdL7kitixB1eAITM1xb1N0/UvbZ/NIgKxrXkXc/8JEkzDdR9rWa1WnYbA+wgp0hP8KR2C0KxBFwl+efbE01EFecCgYEAwUk5EmZa2BGuMYuGccXgO6lYHrnAA3vTRPPWkBahKffAuCYPvyySjvbhO8Abr5sDKKdQPvL9Ck++KLc7tCAps8ExoGJ0Tq3TPaXVE3CEpFPKefqCCKx7Mkgt5Y/S+BjksW59C4Q8z0dW95wek2ap7iZxxJa7yI6dtXevIPjxd1kCgYA/GbxTYQqEymRTsG//fO0EywmtRsvGnNS5m38JU3ULWbJPvZUdtPomnE+SXzHwyN/vFSPBDwOgKclmEYdL6me/Jy6FWpH5taBAN8UWEO6O1eelFhxuQ6+nCi/PjJmGXi4zC82KX3Z0Dy+KiLIwyIiaGDeZWONqP4UHCUuJHVEk3wKBgE/IldJOhbkiszCoUzqrXz/BSyqDqgrGFhMkQ7D+ZlAYgGiC7YUQNP7mUVqElekKp2ckiS8yxdh3yhqsZ+yWSiB04rM9cJz5i3Sq+yUnENlz7OQkz4AdEk1TFf7oO0FFpUDIRr12PFOjMvbKbqSRgBtZqyRmw+SpWdgKKzQFDkchAoGBALZYSAe/COTJ2d5WZW6+HGMTNkuzG7FDn0dkwfTcBcdaxfD/i2d3lRpe5S3HKUMTOTfX2bbaJTZYMG6Vx3YXWa92zhnnI0fXRIboYjye1bMQ3FW+AZWhAEFxYuJOC8G1IDcHSf0nu+1E1AXKvtLGkqhwWJGjM3n8hIww5HrOY5uD"
	client, err := alipay.New("2021000148652076", privateKey, false)
	if err != nil {
		log.NewHelper(logger).Errorf("初始化支付宝客户端失败: %v", err)
		return &PaymentService{log: log.NewHelper(logger)}
	}

	return &PaymentService{
		log:          log.NewHelper(logger),
		alipayClient: client,
	}
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentReply, error) {
	rawReq, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return &pb.UpdatePaymentReply{
			Result: "fail",
		}, nil
	}

	if err := rawReq.ParseForm(); err != nil {
		s.log.Error("解析表单参数失败:", err)
		return &pb.UpdatePaymentReply{
			Result: "fail",
		}, nil
	}
	params := rawReq.PostForm

	paramMap := make(map[string]string)
	for k, v := range params {
		if len(v) > 0 {
			paramMap[k] = v[0]
		}
	}
	if err := s.alipayClient.VerifySign(params); err != nil {
		s.log.Error("签名验证失败:", err, "参数:", paramMap)
		return &pb.UpdatePaymentReply{
			Result: "fail",
		}, nil
	}

	outTradeNo := paramMap["out_trade_no"] // 你的系统订单号
	fmt.Println(outTradeNo)
	tradeStatus := paramMap["trade_status"] // 支付状态（SUCCESS 表示成功）
	totalAmount := paramMap["total_amount"] // 支付金额
	alipayTradeNo := paramMap["trade_no"]   // 支付宝交易号

	fmt.Printf("收到支付宝回调: 订单号=%s, 状态=%s, 金额=%s\n", outTradeNo, tradeStatus, totalAmount)

	var order config.SfOrders

	//TRADE_FINISHED	交易完成	true（触发通知）
	//TRADE_SUCCESS	支付成功	true（触发通知）
	//WAIT_BUYER_PAY	交易创建	false（不触发通知）
	//TRADE_CLOSED	交易关闭	true（触发通知）

	if tradeStatus == "TRADE_SUCCESS" {
		order.OrderStatus = "paid"
		err := order.UpdateOrderStatus(inits.DB, alipayTradeNo)
		if err != nil {
			fmt.Println("订单状态修改失败")
			return nil, errors.New(400, "ORDER_UPDATE_FAILED", "订单状态修改失败")
		}

	}

	// 5. 返回结果（必须返回 "success"，否则支付宝会重复回调）
	return &pb.UpdatePaymentReply{Result: "success"}, nil
}

func (s *PaymentService) CreatedReconciliation(ctx context.Context, req *pb.CreatedReconciliationRequest) (*pb.CreatedReconciliationReply, error) {
	var orders []config.SfOrders
	todaystr3 := time.Now().Format("2006-01-02")
	fmt.Println("todaystr3:", todaystr3)

	// 使用请求中的日期，如果没有则使用今天
	queryDate := req.CreatedTime
	if queryDate == "" {
		queryDate = todaystr3
	}

	err := inits.DB.Model(orders).Where("created_at like ?", "%"+queryDate+"%").Find(&orders).Error
	if err != nil {
		fmt.Println("查询今日订单失败", err.Error())
		return nil, err
	}
	fmt.Println(orders)

	sum := 0.00
	for _, order := range orders {
		sum += order.ActualFee
	}

	// 查询前一天的对账数据
	var previousReconciliation config.Reconciliation
	nowTime := time.Now()
	getTime := nowTime.AddDate(0, 0, -1) //年，月，日   获取一天前的时间
	resTime := getTime.Format("2006-01-02")

	err = inits.DB.Model(previousReconciliation).Where("created_at like ?", "%"+resTime+"%").Limit(1).Find(&previousReconciliation).Error
	if err != nil {
		fmt.Println("查询前一天数据失败", err.Error())
		return nil, err
	}

	// 计算涨幅比例，避免除零错误
	var proportion float64
	if previousReconciliation.ActualAmount > 0 {
		proportion = sum / previousReconciliation.ActualAmount
	} else {
		proportion = 0.0 // 如果没有前一天数据，设为0.0
	}

	reconciliation := config.Reconciliation{
		TaskName:           fmt.Sprintf("%s对账", todaystr3),
		ReconciliationNo:   uuid.NewString(),
		ReconciliationDate: time.Now(),
		Amount:             sum,
		ActualAmount:       sum,
		Proportion:         proportion,
		HandlerStatus:      "待处理", // 设置默认状态
	}

	err = reconciliation.Created(inits.DB)
	if err != nil {
		fmt.Println("添加日账单失败", err.Error())
		return nil, err
	}

	return &pb.CreatedReconciliationReply{
		ReconciliationId: int64(reconciliation.ID),
	}, nil
}
func (s *PaymentService) ListReconciliation(ctx context.Context, req *pb.ListReconciliationRequest) (*pb.ListReconciliationReply, error) {
	var ra config.Reconciliation

	list, err := ra.ReconciliationItemList(inits.DB, req.PlayTime, req.EndTime)
	if err != nil {
		fmt.Println("查询失败", err.Error())
		return nil, err
	}
	return &pb.ListReconciliationReply{
		List: list,
	}, nil
}

func (s *PaymentService) PaymentOrder(ctx context.Context, req *pb.PaymentOrderRequest) (*pb.PaymentOrderReply, error) {
	var o config.SfOrders
	err := o.FIndByOrderSn(inits.DB, req.OrderSn)
	if err != nil {
		logs.Println("查询订单号错误", err.Error())
		return nil, err
	}

	var pay utils.AliPay

	price := fmt.Sprintf("%.2f", o.ActualFee)
	url := pay.Pay(o.OrderNo, price)
	return &pb.PaymentOrderReply{
		Url: url,
	}, nil
}
