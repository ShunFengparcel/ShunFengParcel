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
	"github.com/robfig/cron/v3"
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

	// 加载支付宝公钥用于验证回调签名
	// ⚠️ 注意：这里必须使用支付宝公钥，不是应用公钥！
	// 支付宝公钥从支付宝开放平台的"开发信息"中获取
	alipayPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuTeRNKVmgIl3iMzw3X7y76e5M0PWKT8LmfRXbOSC5/JMIV27MIy9Zns/ym98jW4dOI4W+0PO47k4JKdeZzzJBjOBVcXM1lru4m8aMclRqVGLhKWDVSn2CiL8NbOf2fDMflDYlW8eCrBBUTnoTGkhYijTVPL6av3GNCB5WhtGUjyRinMXBo43Xdujgz8SqU2V5Y+tRZoZeMu8Hz68OBCO+yLwT59JDnASsUaHZ6Axk71ZOkINruWMGoxMKL87/Q7+Zggh1tAcVm+UPnaIGf9DTkddBg45+mJL2KOC/J1b7FQqE8VfftR1Ybd88/fxLVF/1hxkvWAMyrpcSXkt/dP4qwIDAQAB"
	if err := client.LoadAliPayPublicKey(alipayPublicKey); err != nil {
		log.NewHelper(logger).Errorf("加载支付宝公钥失败: %v", err)
		return &PaymentService{log: log.NewHelper(logger)}
	}

	return &PaymentService{
		log:          log.NewHelper(logger),
		alipayClient: client,
	}
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentReply, error) {
	defer Wg.Done()
	Wg.Add(1)
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

	// 支付宝回调可能通过 POST Form 或 GET Query 传递参数
	// 需要合并两种方式的参数
	params := rawReq.PostForm

	// 先处理 POST 参数
	for k, v := range params {
		params[k] = v
	}

	// 再处理 GET 参数（如果有的话）
	for k, v := range rawReq.URL.Query() {
		if _, exists := params[k]; !exists {
			params[k] = v
		}
	}

	// 创建 map 用于日志和后续使用
	paramMap := make(map[string]string)
	for k, v := range params {
		if len(v) > 0 {
			paramMap[k] = v[0]
		}
	}

	// 记录接收到的签名和部分参数用于调试
	s.log.Infof("收到支付宝回调 - 订单号: %s, 签名: %s", paramMap["out_trade_no"], paramMap["sign"][:50]+"...")

	// 使用 url.Values 进行验签
	if err := s.alipayClient.VerifySign(params); err != nil {
		s.log.Errorf("❌ 签名验证失败: %v", err)
		s.log.Errorf("请检查以下配置:")
		s.log.Errorf("1. 确认使用的是【支付宝公钥】不是【应用公钥】")
		s.log.Errorf("2. 支付宝公钥需要从支付宝开放平台获取")
		s.log.Errorf("3. 沙箱环境: https://openhome.alipay.com/platform/appDaily.htm")
		s.log.Errorf("4. 正式环境: https://open.alipay.com/")
		s.log.Errorf("当前 app_id: %s, 通知类型: %s", paramMap["app_id"], paramMap["notify_type"])
		return &pb.UpdatePaymentReply{
			Result: "fail",
		}, nil
	}

	s.log.Info("✅ 签名验证成功")

	outTradeNo := paramMap["out_trade_no"] // 你的系统订单号
	fmt.Println(outTradeNo)
	tradeStatus := paramMap["trade_status"] // 支付状态（SUCCESS 表示成功）
	totalAmount := paramMap["total_amount"] // 支付金额
	alipayTradeNo := paramMap["trade_no"]   // 支付宝交易号

	fmt.Printf("收到支付宝回调: 订单号=%s, 状态=%s, 金额=%s\n", outTradeNo, tradeStatus, totalAmount)

	go func() {
		err := CreatePayment(outTradeNo, alipayTradeNo, tradeStatus)
		if err != nil {
			return
		}
	}()

	var order config.SfOrders

	//TRADE_FINISHED	交易完成	true（触发通知）
	//TRADE_SUCCESS	支付成功	true（触发通知）
	//WAIT_BUYER_PAY	交易创建	false（不触发通知）
	//TRADE_CLOSED	交易关闭	true（触发通知）mm

	err := order.FIndByOrderSn(inits.DB, outTradeNo)
	if err != nil {
		return nil, err
	}

	if tradeStatus == "TRADE_SUCCESS" {
		order.PaymentStatus = "paid"

		err := order.UpdateOrderStatus(inits.DB, outTradeNo)
		if err != nil {
			utils.LogPaymentError(outTradeNo, "订单状态修改失败", err)
			return nil, errors.New(400, "ORDER_UPDATE_FAILED", "订单状态修改失败")
		}
		inits.RDB.Del(context.Background(), "payment-order:"+outTradeNo)
	}
	Wg.Wait()
	// 5. 返回结果（必须返回 "success"，否则支付宝会重复回调）
	return &pb.UpdatePaymentReply{Result: "success"}, nil
}

func CreatePayment(ordersn string, alipayNo string, status string) error {
	var order config.SfOrders
	err := order.FIndByOrderSn(inits.DB, ordersn)
	if err != nil {
		logs.Println("查询订单失败", err.Error())
		return err
	}
	fmt.Println(order)
	var pay config.SfPayments
	pay = config.SfPayments{
		OrderId:   order.Id,
		PayNo:     ordersn,
		Channel:   "alipay",
		Amount:    order.ActualFee,
		Status:    status,
		PaidAt:    time.Now(),
		ThirdTxId: alipayNo,
	}

	err = pay.Created(inits.DB)
	if err != nil {
		logs.Println("支付记录失败", err.Error())
		return err
	}

	return nil
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
	count := 0
	for _, order := range orders {
		count++
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
		Count:              int64(count),
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

	// 转换为 protobuf 类型
	var pbList []*pb.ReconciliationItem
	for _, rec := range list {
		pbList = append(pbList, &pb.ReconciliationItem{
			TaskName:           rec.TaskName,
			ReconciliationNo:   rec.ReconciliationNo,
			ReconciliationDate: rec.ReconciliationDate.Format("2006-01-02 15:04:05"),
			Amount:             float32(rec.Amount),
			ActualAmount:       float32(rec.ActualAmount),
			Proportion:         fmt.Sprintf("%.2f", rec.Proportion),
			HandlerStatus:      rec.HandlerStatus,
		})
	}

	return &pb.ListReconciliationReply{
		List: pbList,
	}, nil
}

func (s *PaymentService) PaymentOrder(ctx context.Context, req *pb.PaymentOrderRequest) (*pb.PaymentOrderReply, error) {
	var o config.SfOrders
	err := o.FIndByOrderSn(inits.DB, req.OrderSn)
	if err != nil {
		utils.LogOrderError(req.OrderSn, "查询订单号错误", err)
		return nil, err
	}

	if o.Id == 0 {
		utils.LogOrderError(req.OrderSn, "该订单不存在", err)
		return nil, errors.New(400, err.Error(), "该订单不存在")
	}

	if o.PaymentStatus == "paid" {
		return nil, errors.New(500, "paid", "该订单已被支付")
	}

	var url string
	if req.PaymentType == "1" {
		var pay utils.AliPay

		price := fmt.Sprintf("%.2f", o.ActualFee)
		url = pay.Pay(o.OrderNo, price)
	}

	inits.RDB.Set(context.Background(), "payment-order:"+o.OrderNo, o, time.Minute*15)

	return &pb.PaymentOrderReply{
		Url: url,
	}, nil
}

func (s *PaymentService) MonitorCreate(ctx context.Context, req *pb.MonitorCreateRequest) (*pb.MonitorCreateReply, error) {
	var err error
	var o config.SfOrders
	err = o.FIndByOrderSn(inits.DB, req.OrderNo)
	if err != nil {
		logs.Println("查询订单号错误", err.Error())
		return nil, err
	}
	Level := ""
	score := 0
	if o.ActualFee >= 100 {
		Level = "low"
	} else if o.ActualFee >= 1000 {
		Level = "medium"
	} else if o.ActualFee >= 10000 {
		Level = "high"
	}

	var m config.TransactionMonitor
	m = config.TransactionMonitor{
		MonitorNo:         uuid.NewString(),
		OrderNo:           req.OrderNo,
		UserId:            o.UserId,
		TransactionAmount: o.ActualFee,
		MonitorType:       req.MonitorType,
		RiskLevel:         Level,
		RiskScore:         int64(score),
	}

	if err = m.Created(inits.DB); err != nil {
		logs.Println("监控订单失败", err.Error())
		return nil, err
	}
	return &pb.MonitorCreateReply{
		Id: int64(m.ID),
	}, nil
}
func (s *PaymentService) MonitorUpdate(ctx context.Context, req *pb.MonitorUpdateRequest) (*pb.MonitorUpdateReply, error) {
	var err error
	var m config.TransactionMonitor
	if err = m.Updated(inits.DB, req.MonitorNo); err != nil {
		logs.Println("监控情况查询失败", err.Error())
		return nil, err
	}

	m.HandleResult = req.HandleResult
	if err = m.Updated(inits.DB, req.MonitorNo); err != nil {
		logs.Println("监控情况修改失败", err.Error())
		return nil, err
	}
	return &pb.MonitorUpdateReply{
		Id: int64(m.ID),
	}, nil
}
func (s *PaymentService) MonitorDelete(ctx context.Context, req *pb.MonitorDeleteRequest) (*pb.MonitorDeleteReply, error) {
	var err error
	var m config.TransactionMonitor

	if err = m.Deleted(inits.DB, int(req.Id)); err != nil {
		logs.Println("监控删除失败", err.Error())
		return nil, err
	}

	if err = m.FindByID(inits.DB, int(req.Id)); err != nil {
		logs.Println("监控情况查询失败", err.Error())
		return nil, err
	}
	message := ""
	if m.ID == 0 {
		message = "删除成功"
	} else {
		logs.Println("未删除监控", m.ID)
	}

	return &pb.MonitorDeleteReply{
		Message: message,
	}, nil
}

func (s *PaymentService) OrderList(ctx context.Context, req *pb.OrderListRequest) (*pb.OrderListReply, error) {
	var o config.SfOrders
	list, err := o.FIndByList(inits.DB, req.Status)
	if err != nil {
		logs.Println(err.Error())
		return nil, err
	}

	// 转换为 protobuf 类型
	var pbList []*pb.OrderItem
	for _, order := range list {
		pbList = append(pbList, &pb.OrderItem{
			OrderNo:         order.OrderNo,
			SenderName:      order.SenderName,
			SenderPhone:     order.SenderPhone,
			SenderAddress:   order.SenderAddress,
			ReceiverName:    order.ReceiverName,
			ReceiverPhone:   order.ReceiverPhone,
			ReceiverAddress: order.ReceiverAddress,
			ProductType:     order.ProductType,
			ActualFee:       fmt.Sprintf("%.2f", order.ActualFee),
			CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.OrderListReply{
		List: pbList,
	}, nil
}

func (s *PaymentService) PaymentList(ctx context.Context, req *pb.PaymentListRequest) (*pb.PaymentListReply, error) {
	var p config.SfPayments
	list, err := p.FIndByList(inits.DB, req.Status)
	if err != nil {
		logs.Println("查询失败", err.Error())
		return nil, err
	}
	fmt.Println(list)
	if list == nil {
		logs.Println("未找到支付信息")
	}

	return &pb.PaymentListReply{
		List: list,
	}, nil
}

func AddFunc() {
	// 创建调度器（默认支持分钟级，如需秒级需添加 WithSeconds() 选项）
	c := cron.New(cron.WithSeconds()) // 支持秒级调度

	// 注册任务：每5秒执行一次（Cron表达式："0 2 * * *"）
	_, err := c.AddFunc("0 2  * * *", func() {

		fmt.Printf("任务执行时间：%v\n", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Printf("注册任务失败：%v\n", err)
		return
	}

	// 启动调度器
	c.Start()
	defer c.Stop() // 程序退出时停止调度器

	// 阻塞主线程（避免程序退出）
	select {}
}
