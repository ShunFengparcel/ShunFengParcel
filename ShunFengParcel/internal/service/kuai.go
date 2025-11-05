package service

import (
	pb "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/basic/config"
	"ShunFengParcel/internal/biz"
	"ShunFengParcel/internal/model"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"ShunFengParcel/internal/queue"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *KuaiService) CourierAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				courierIDStr := tr.RequestHeader().Get("X-Courier-ID")
				if courierID, err := strconv.ParseInt(courierIDStr, 10, 64); err == nil {
					ctx = context.WithValue(ctx, "courier_id", courierID)
				}
			}
			return handler(ctx, req)
		}
	}
}

type KuaiService struct {
	pb.UnimplementedKuaiServer
	uc  *biz.GreeterUsecase
	RDB redis.UniversalClient
}

func NewKuaiService(uc *biz.GreeterUsecase, rdb redis.UniversalClient) *KuaiService {
	return &KuaiService{
		uc:  uc,
		RDB: rdb, // 把 Redis 客户端带进来
	}
}
func (s *KuaiService) DeleteKuai(ctx context.Context, req *pb.DeleteKuaiRequest) (*pb.DeleteKuaiReply, error) {
	return &pb.DeleteKuaiReply{}, nil
}

var (
	taskTypeMap = map[string]pb.TaskType{
		"pickup":   pb.TaskType_Pickup,
		"delivery": pb.TaskType_Delivery,
	}
	taskStatusMap = map[string]pb.TaskStatus{
		"pending":   pb.TaskStatus_Pending,
		"accepted":  pb.TaskStatus_Accepted,
		"completed": pb.TaskStatus_Completed,
		"cancelled": pb.TaskStatus_Cancelled,
	}
)

func (s *KuaiService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskReply, error) {
	var t model.SfCourierTasks
	config.DB.Where("id=?", req.Id).Find(&t)
	var courier model.SfCouriers
	config.DB.Where("id = ?", t.CourierId).First(&courier)
	return &pb.GetTaskReply{
		Id:          t.Id,
		CourierId:   t.CourierId,
		CourierName: courier.RealName,
		TaskType:    taskTypeMap[t.TaskType],
		Priority:    int32(t.Priority),
		TaskStatus:  taskStatusMap[t.TaskStatus],
	}, nil
}
func (s *KuaiService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	var c model.SfCouriers
	result := config.DB.Where("phone = ?", req.Phone).First(&c)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或手机号错误")
		}
		return nil, result.Error
	}
	if c.PasswordHash != req.PasswordHash {
		return nil, errors.New("密码错误")
	}
	fmt.Println("111")
	return &pb.LoginReply{
		Id: c.Id,
	}, nil
}
func (s *KuaiService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	couriers := model.SfCouriers{
		Phone:        req.Phone,
		PasswordHash: req.PasswordHash,
	}
	config.DB.Create(&couriers)
	return &pb.RegisterReply{
		Id: couriers.Id,
	}, nil
}
func (s *KuaiService) StuUpd(ctx context.Context, req *pb.StuUpdRequest) (*pb.StuUpdReply, error) {
	couriers := model.SfCouriers{Status: req.Status}
	config.DB.Where("phone=?", req.Phone).Updates(&couriers)

	return &pb.StuUpdReply{
		Status: req.Status,
	}, nil
}

func (s *KuaiService) TaskList(ctx context.Context, req *pb.TaskListRequest) (*pb.TaskListReply, error) {
	// 使用联查获取任务和对应的订单信息
	var results []struct {
		model.SfCourierTasks
		SenderAddress   string `gorm:"column:sender_address"`
		ReceiverAddress string `gorm:"column:receiver_address"`
	}

	if err := config.DB.Table("sf_courier_tasks").
		Select("sf_courier_tasks.*, sf_orders.sender_address, sf_orders.receiver_address").
		Joins("LEFT JOIN sf_orders ON sf_courier_tasks.order_id = sf_orders.id").
		Find(&results).Error; err != nil {
		return nil, err
	}

	pbTasks := make([]*pb.Task, len(results))
	for i, result := range results {
		var courier model.SfCouriers
		config.DB.Where("id = ?", result.CourierId).First(&courier)

		pbTasks[i] = &pb.Task{
			Id:              result.Id,
			CourierId:       result.CourierId,
			CourierName:     courier.RealName,
			TaskType:        taskTypeMap[result.TaskType],
			Priority:        int32(result.Priority),
			TaskStatus:      taskStatusMap[result.TaskStatus],
			CreatedAt:       result.CreatedAt.Format("2006-01-02 15:04:05"),
			SenderAddress:   result.SenderAddress,
			ReceiverAddress: result.ReceiverAddress,
		}
	}
	return &pb.TaskListReply{
		Tasks: pbTasks,
		Total: int32(len(pbTasks)),
	}, nil
}

// 收入明细列表：展示已完成任务，联查订单实际费用、完成时间与寄件人姓名
func (s *KuaiService) IncomeList(ctx context.Context, req *pb.IncomeListRequest) (*pb.IncomeListReply, error) {
	// 获取快递员ID：优先取请求参数，其次取请求头中间件注入
	courierID := req.CourierId
	if courierID == 0 {
		if v, ok := ctx.Value("courier_id").(int64); ok {
			courierID = v
		}
	}
	if courierID == 0 {
		return nil, status.Error(codes.InvalidArgument, "缺少快递员ID")
	}

	// 无分页：一次性返回所有已完成任务的收入明细

	var results []struct {
		Id          int64     `gorm:"column:id"`
		OrderId     int64     `gorm:"column:order_id"`
		Fee         float64   `gorm:"column:fee"`
		CompletedAt time.Time `gorm:"column:completed_at"`
		SenderName  string    `gorm:"column:sender_name"`
	}

	db := config.DB.Table("sf_courier_tasks AS t").
		Select("t.id AS id, t.order_id AS order_id, o.actual_fee AS fee, COALESCE(o.actual_delivery_time, t.updated_at) AS completed_at, o.sender_name AS sender_name").
		Joins("LEFT JOIN sf_orders AS o ON t.order_id = o.id").
		Where("t.task_status = ? AND t.courier_id = ?", "completed", courierID).
		Order("completed_at DESC")

	if err := db.Find(&results).Error; err != nil {
		return nil, status.Error(codes.Internal, "查询收入明细失败")
	}

	list := make([]*pb.IncomeItem, 0, len(results))
	for _, r := range results {
		completed := ""
		if !r.CompletedAt.IsZero() {
			completed = r.CompletedAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, &pb.IncomeItem{
			TaskId:        r.Id,
			OrderId:       r.OrderId,
			Fee:           r.Fee,
			CompletedTime: completed,
			SenderName:    r.SenderName,
		})
	}

	return &pb.IncomeListReply{
		List:  list,
		Total: int32(len(list)),
	}, nil
}

// 绩效统计
// GetCourierPerformanceRequest
func (s *KuaiService) GetCourierPerformance(ctx context.Context, req *pb.GetCourierPerformanceRequest) (*pb.GetCourierPerformanceReply, error) {
	courierID, ok := ctx.Value("courier_id").(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid courier id")
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	dailyRevenue, err := s.calculateRevenue(ctx, courierID, startOfDay)
	if err != nil {
		return nil, err
	}

	monthlyRevenue, err := s.calculateRevenue(ctx, courierID, startOfMonth)
	if err != nil {
		return nil, err
	}

	// 统计接单数：根据快递员ID统计 accepted 或 completed 的任务数量
	var acceptedCount int64
	if err := config.DB.WithContext(ctx).
		Model(&model.SfCourierTasks{}).
		Where("courier_id = ? AND task_status IN (?)", courierID, []string{"accepted", "completed"}).
		Count(&acceptedCount).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to count accepted tasks")
	}

	return &pb.GetCourierPerformanceReply{
		DailyRevenue:   dailyRevenue,
		MonthlyRevenue: monthlyRevenue,
		AcceptedCount:  acceptedCount,
	}, nil
}

func (s *KuaiService) calculateRevenue(ctx context.Context, courierID int64, startTime time.Time) (float64, error) {
	var tasks []model.SfCourierTasks
	if err := config.DB.WithContext(ctx).Where("courier_id = ? AND task_status = ? AND updated_at >= ?", courierID, "completed", startTime).Find(&tasks).Error; err != nil {
		return 0, status.Error(codes.Internal, "failed to query tasks")
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	var orderIDs []int64
	for _, task := range tasks {
		orderIDs = append(orderIDs, task.OrderId)
	}

	var orders []model.SfOrders
	if err := config.DB.WithContext(ctx).Where("id IN ?", orderIDs).Find(&orders).Error; err != nil {
		return 0, status.Error(codes.Internal, "failed to query orders")
	}

	var totalRevenue float64
	for _, order := range orders {
		totalRevenue += order.ActualFee
	}

	return totalRevenue, nil
}

// 订单列表
func (s *KuaiService) OrderList(ctx context.Context, req *pb.OrderListRequest) (*pb.OrderListReply, error) {
	var orders []model.SfOrders

	// 默认只查询 pending 状态的订单
	if err := config.DB.Where("order_status = ?", "pending").Order("id desc").Find(&orders).Error; err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}

	items := make([]*pb.OrderItem, 0, len(orders))
	for _, o := range orders {
		items = append(items, &pb.OrderItem{
			Id:                   o.Id,
			OrderNo:              o.OrderNo,
			SenderName:           o.SenderName,
			ReceiverName:         o.ReceiverName,
			ServiceType:          o.ServiceType,
			OrderStatus:          o.OrderStatus,
			CreatedAt:            o.CreatedAt.Format("2006-01-02 15:04:05"),
			UserId:               o.UserId,
			CourierId:            o.CourierId,
			SenderPhone:          o.SenderPhone,
			SenderAddress:        o.SenderAddress,
			ReceiverPhone:        o.ReceiverPhone,
			ReceiverAddress:      o.ReceiverAddress,
			ProductType:          o.ProductType,
			IsInsured:            int32(o.IsInsured),
			InsuredValue:         o.InsuredValue,
			EstimatedWeight:      o.EstimatedWeight,
			ActualWeight:         o.ActualWeight,
			VolumeWeight:         o.VolumeWeight,
			EstimatedFee:         o.EstimatedFee,
			ActualFee:            o.ActualFee,
			PaymentMethod:        o.PaymentMethod,
			PaymentStatus:        o.PaymentStatus,
			ExpectedPickupTime:   o.ExpectedPickupTime.Unix(),
			ActualPickupTime:     o.ActualPickupTime.Unix(),
			ExpectedDeliveryTime: o.ExpectedDeliveryTime.Unix(),
			ActualDeliveryTime:   o.ActualDeliveryTime.Unix(),
			UpdatedAt:            o.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: func() string {
				if o.DeleteAt.IsZero() {
					return ""
				} else {
					return o.DeleteAt.Format("2006-01-02 15:04:05")
				}
			}(),
		})
	}
	return &pb.OrderListReply{List: items}, nil
}

// 根据新的优先级规则计算任务优先级
// P0（最高）：超时订单
// P1（高）：预约时间即将到期的订单
// P2（中）：普通已接单订单
// P3（低）：新订单
func calcPriority(req *pb.CreateOrderRequest) int8 {
	// 获取当前时间作为基准
	now := time.Now()

	// 将请求中的期望上门时间(Unix 秒时间戳)转换为 time.Time 类型
	expectedPickupTime := time.Unix(req.ExpectedPickupTime, 0)

	// P0: 如果期望上门时间早于当前时间，则认为订单已超时，优先级最高(0)
	if !expectedPickupTime.IsZero() && expectedPickupTime.Before(now) {
		return 0
	}

	// P1: 如果期望上门时间与当前时间差小于等于 1 小时，则视为即将到期，次高优先级(1)
	if !expectedPickupTime.IsZero() && expectedPickupTime.Sub(now) <= time.Minute*30 {
		return 1
	}

	// P2: 如果订单状态为已接单，则按普通已接单订单处理，优先级(2)
	if req.OrderStatus == pb.OrderStatus_ORDER_STATUS_ACCEPTED {
		return 2
	}

	// P3: 其他情况(如新订单)默认最低优先级(3)
	return 3
}

// 创建订单放入消息队列优先级
func (s *KuaiService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	var serviceTypeMap = map[pb.ServiceType]string{
		pb.ServiceType_DELIVERY_OPTION_SAME_DAY: "same_day",
		pb.ServiceType_DELIVERY_OPTION_NEXT_DAY: "next_day",
		pb.ServiceType_DELIVERY_OPTION_STANDARD: "standard",
		pb.ServiceType_DELIVERY_OPTION_ECONOMY:  "economy",
	}
	var PaymentMethodToDB = map[pb.PaymentMethod]string{
		pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:  "sender_pay", // 默认寄付
		pb.PaymentMethod_PAYMENT_METHOD_SENDER_PAY:   "sender_pay",
		pb.PaymentMethod_PAYMENT_METHOD_RECEIVER_PAY: "receiver_pay",
		pb.PaymentMethod_PAYMENT_METHOD_MONTHLY:      "monthly",
	}
	var PaymentStatusToDB = map[pb.PaymentStatus]string{
		pb.PaymentStatus_PAYMENT_STATUS_PENDING:  "pending",
		pb.PaymentStatus_PAYMENT_STATUS_PAID:     "paid",
		pb.PaymentStatus_PAYMENT_STATUS_FAILED:   "failed",
		pb.PaymentStatus_PAYMENT_STATUS_REFUNDED: "refunded",
	}
	var OrderStatusToDB = map[pb.OrderStatus]string{
		pb.OrderStatus_ORDER_STATUS_PENDING:          "pending",
		pb.OrderStatus_ORDER_STATUS_ACCEPTED:         "accepted",
		pb.OrderStatus_ORDER_STATUS_PICKED_UP:        "picked_up",
		pb.OrderStatus_ORDER_STATUS_IN_TRANSIT:       "in_transit",
		pb.OrderStatus_ORDER_STATUS_OUT_FOR_DELIVERY: "out_for_delivery",
		pb.OrderStatus_ORDER_STATUS_DELIVERED:        "delivered",
		pb.OrderStatus_ORDER_STATUS_CANCELLED:        "cancelled",
		pb.OrderStatus_ORDER_STATUS_EXCEPTION:        "exception",
	}
	orders := model.SfOrders{
		OrderNo:              uuid.NewString(),
		UserId:               req.UserId,
		SenderName:           req.SenderName,
		SenderPhone:          req.SenderPhone,
		SenderAddress:        req.SenderAddress,
		ReceiverName:         req.ReceiverName,
		ReceiverPhone:        req.ReceiverPhone,
		ReceiverAddress:      req.ReceiverAddress,
		ServiceType:          serviceTypeMap[req.ServiceType],
		ProductType:          req.ProductType,
		IsInsured:            int8(req.IsInsured),
		InsuredValue:         req.InsuredValue,
		EstimatedWeight:      req.EstimatedWeight,
		ActualWeight:         req.ActualWeight,
		VolumeWeight:         req.VolumeWeight,
		EstimatedFee:         req.EstimatedFee,
		ActualFee:            req.ActualFee,
		PaymentMethod:        PaymentMethodToDB[req.PaymentMethod],
		PaymentStatus:        PaymentStatusToDB[req.PaymentStatus],
		ExpectedPickupTime:   time.Unix(req.ExpectedPickupTime, 0),
		ActualPickupTime:     time.Unix(req.ActualPickupTime, 0),
		ExpectedDeliveryTime: time.Unix(req.ExpectedDeliveryTime, 0),
		ActualDeliveryTime:   time.Unix(req.ActualDeliveryTime, 0),
		OrderStatus:          OrderStatusToDB[req.OrderStatus],
	}
	if err := config.DB.Create(&orders).Error; err != nil {
		return nil, errors.New("订单创建失败")
	}
	priority := calcPriority(req)

	// 创建待接任务并入库，等待快递员接单
	task := model.SfCourierTasks{
		OrderId:    orders.Id, // 关联订单ID
		CourierId:  0,
		TaskType:   "pickup",
		Priority:   priority,
		TaskStatus: "pending",
	}
	if err := config.DB.Create(&task).Error; err != nil {
		return nil, errors.New("任务创建失败")
	}

	// 翻译成 Redis 队列名并入队任务ID
	key := queue.KeyByPriority(priority)
	// 用有序集合 zadd，Member 使用任务ID
	if err := config.RDB.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: task.Id,
	}).Err(); err != nil {
		return nil, errors.New("加入队列失败")
	}

	return &pb.CreateOrderReply{OrderId: int32(orders.Id)}, nil
}

// ------- 幂等与审计辅助 -------
func (s *KuaiService) getIdempotencyKey(ctx context.Context, key string) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if h := tr.RequestHeader().Get("X-Idempotency-Key"); h != "" {
			return h
		}
	}
	if key != "" {
		return key
	}
	return uuid.NewString()
}

// 返回 true 表示首次执行（已登记），false 表示重复
func (s *KuaiService) checkAndSetIdempotency(ctx context.Context, scope, key string, ttl time.Duration) (bool, error) {
	if s.RDB == nil {
		// 无 Redis 时降级允许执行
		return true, nil
	}
	k := fmt.Sprintf("idem:%s:%s", scope, key)
	ok, err := s.RDB.SetNX(ctx, k, 1, ttl).Result()
	return ok, err
}

func operatorFromCtx(ctx context.Context) (id int64, role string) {
	if v, ok := ctx.Value("courier_id").(int64); ok && v > 0 {
		return v, "courier"
	}
	return 0, "system"
}

// CancelOrder 处理“取消订单”请求，支持用户/快递员主动取消，全程幂等、事务、审计。
func (s *KuaiService) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderReply, error) {

	// 从请求上下文里取出客户端生成的幂等键（Idempotency-Key）
	idem := s.getIdempotencyKey(ctx, req.IdempotencyKey)

	// 去 Redis 登记：如果已经登记过 表示重复请求
	first, err := s.checkAndSetIdempotency(ctx,
		fmt.Sprintf("cancel:order:%d", req.OrderId), // Redis key 前缀
		idem,
		time.Minute*10) // 10 分钟内重复请求直接短路
	if err != nil {
		return nil, status.Errorf(codes.Internal, "幂等登记失败: %v", err)
	}

	// 订单日志表，订单改派记录表
	_ = config.DB.AutoMigrate(&model.SfOrderAuditLogs{}, &model.SfOrderReassignments{})

	// 取出当前操作人信息（JWT 解析）
	opId, opRole := operatorFromCtx(ctx)

	//  事务开始
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // panic 保护
		}
	}()

	// 查询订单并加行锁（FOR UPDATE）
	//查到订单并立刻加写锁，防止别人同时改这条订单，避免并发冲突
	var order model.SfOrders
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", req.OrderId).
		First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.CancelOrderReply{
				Success:   false,
				Message:   "订单不存在",
				NewStatus: pb.OrderStatus_ORDER_STATUS_UNSPECIFIED,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "查询订单失败: %v", err)
	}

	// 已签收 → 不能取消
	if order.OrderStatus == "delivered" {
		tx.Rollback()
		return &pb.CancelOrderReply{
			Success:   false,
			Message:   "订单已签收，无法取消",
			NewStatus: pb.OrderStatus_ORDER_STATUS_DELIVERED,
		}, nil
	}
	// 已取消 → 幂等返回成功
	if order.OrderStatus == "cancelled" {
		return &pb.CancelOrderReply{
			Success:   true,
			Message:   "订单已取消（幂等）",
			NewStatus: pb.OrderStatus_ORDER_STATUS_CANCELLED,
		}, nil
	}
	// 快递员只能取消自己的单
	if req.Reason == pb.CancelReason_COURIER_REQUEST &&
		order.CourierId > 0 &&
		opRole == "courier" &&
		opId != order.CourierId {
		tx.Rollback()
		return &pb.CancelOrderReply{
			Success:   false,
			Message:   "无权取消其他快递员订单",
			NewStatus: pb.OrderStatus_ORDER_STATUS_UNSPECIFIED,
		}, nil
	}

	// 更新订单状态
	before := order.OrderStatus
	order.OrderStatus = "cancelled"
	if err := tx.Model(&order).
		Update("order_status", order.OrderStatus).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新订单状态失败: %v", err)
	}

	// 同步取消对应快递员任务
	var task model.SfCourierTasks
	if err := tx.Where("order_id = ?", req.OrderId).First(&task).Error; err == nil {
		_ = tx.Model(&task).
			Updates(map[string]interface{}{"task_status": "cancelled"}).Error
	}

	// 写审计日志
	_ = tx.Create(&model.SfOrderAuditLogs{
		OrderId:        req.OrderId,
		ActionType:     "cancel",
		OperatorId:     opId,
		OperatorRole:   opRole,
		BeforeStatus:   before,
		AfterStatus:    order.OrderStatus,
		ReasonType:     req.Reason.String(),
		ReasonText:     req.ReasonText,
		IdempotencyKey: idem, // 把本次幂等键也落库，方便以后对账
	}).Error

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}

	// 如果是重复请求，告诉调用方“已去重”；第一次请求返回正常成功
	if !first {
		return &pb.CancelOrderReply{
			Success:   true,
			Message:   "重复请求已去重",
			NewStatus: pb.OrderStatus_ORDER_STATUS_CANCELLED,
		}, nil
	}
	return &pb.CancelOrderReply{
		Success:   true,
		Message:   "订单已取消",
		NewStatus: pb.OrderStatus_ORDER_STATUS_CANCELLED,
	}, nil
}

// ReassignOrder 订单改派接口（含协商逻辑）
// 协商场景：用户/运营/系统提出改派 → 校验合法性 → 决定是否立即生效 or 需要额外确认
func (s *KuaiService) ReassignOrder(ctx context.Context, req *pb.ReassignOrderRequest) (*pb.ReassignOrderReply, error) {

	// ① 幂等协商：同一请求 10 分钟内多次点击，直接返回上次结果，避免重复转单
	idem := s.getIdempotencyKey(ctx, req.IdempotencyKey)
	scope := fmt.Sprintf("reassign:order:%d:to:%d", req.OrderId, req.TargetCourierId)
	first, err := s.checkAndSetIdempotency(ctx, scope, idem, time.Minute*10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "幂等登记失败: %v", err)
	}

	// ② 表结构协商：开发阶段自动建表，生产环境可关闭
	_ = config.DB.AutoMigrate(&model.SfOrderAuditLogs{}, &model.SfOrderReassignments{})

	// ③ 操作人协商：从 JWT 解析出“谁”在发起改派（用户/运营/系统）
	opId, opRole := operatorFromCtx(ctx)

	// ④ 事务协商：开启事务，保证“订单+任务+链路+审计”要么一起成功，一起失败
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 异常时回滚，避免半吊子数据
		}
	}()

	// ⑤ 订单存在性协商：锁行读，防止并发改派时订单被删除
	var order model.SfOrders
	if err := tx.Where("id = ?", req.OrderId).First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ReassignOrderReply{Success: false, Message: "订单不存在"}, nil
		}
		return nil, status.Errorf(codes.Internal, "查询订单失败: %v", err)
	}

	// ⑥ 状态合法性协商：已取消/已签收 → 不允许再改派，直接拒绝
	if order.OrderStatus == "cancelled" || order.OrderStatus == "delivered" {
		tx.Rollback()
		return &pb.ReassignOrderReply{Success: false, Message: "订单状态不可改派"}, nil
	}

	// ⑦ 目标快递员协商：已是当前快递员 → 幂等成功，无需再转
	from := order.CourierId
	if from == req.TargetCourierId {
		return &pb.ReassignOrderReply{
			Success:       true,
			Message:       "已改派（幂等）",
			FromCourierId: from,
			ToCourierId:   req.TargetCourierId,
		}, nil
	}

	// ⑧ 改派执行协商：更新订单 courier_id（单行更新，锁住行）
	if err := tx.Model(&order).Update("courier_id", req.TargetCourierId).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新订单快递员失败: %v", err)
	}

	// ⑨ 任务同步协商：把关联的快递员任务也转给新小哥
	var task model.SfCourierTasks
	if err := tx.Where("order_id = ?", req.OrderId).First(&task).Error; err == nil {
		_ = tx.Model(&task).Updates(map[string]interface{}{"courier_id": req.TargetCourierId}).Error
	}

	// ⑩ 改派链路协商：记录“从谁→到谁+原因”，方便后台追溯
	_ = tx.Create(&model.SfOrderReassignments{
		OrderId:       req.OrderId,
		FromCourierId: from,
		ToCourierId:   req.TargetCourierId,
		OperatorId:    opId,
		ReasonText:    req.ReasonText,
	}).Error

	// ⑪ 审计协商：写审计日志，记录“谁”在“什么时间”做了“改派”动作
	_ = tx.Create(&model.SfOrderAuditLogs{
		OrderId:        req.OrderId,
		ActionType:     "reassign",
		OperatorId:     opId,
		OperatorRole:   opRole,
		BeforeStatus:   order.OrderStatus,
		AfterStatus:    order.OrderStatus,
		ReasonType:     "reassign",
		ReasonText:     req.ReasonText,
		IdempotencyKey: idem,
	}).Error

	// ⑫ 提交协商：所有步骤无错误，事务提交，改派正式生效
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}

	// ⑬ 重复请求协商：非首次调用，告诉调用方“已去重”，但仍返回 from/to 快递员 ID
	if !first {
		return &pb.ReassignOrderReply{
			Success:       true,
			Message:       "重复请求已去重",
			FromCourierId: from,
			ToCourierId:   req.TargetCourierId,
		}, nil
	}

	// ⑭ 首次成功协商：返回改派成功 + 原快递员与新快递员 ID
	return &pb.ReassignOrderReply{
		Success:       true,
		Message:       "改派成功",
		FromCourierId: from,
		ToCourierId:   req.TargetCourierId,
	}, nil
}

// ------- 订单详情（缓存优先，读库回源） -------
func (s *KuaiService) OrderDetail(ctx context.Context, req *pb.OrderDetailRequest) (*pb.OrderDetailReply, error) {
	// 权限隔离：如果请求来自快递员，只能查看由其接单的订单
	courierID, role := operatorFromCtx(ctx)

	// 构建缓存键
	keyPart := req.OrderNo
	if keyPart == "" {
		keyPart = fmt.Sprintf("%d", req.OrderId)
	}
	if keyPart == "" {
		return nil, status.Error(codes.InvalidArgument, "缺少订单号或订单ID")
	}
	cacheKey := fmt.Sprintf("order:detail:%s", keyPart)

	// 优先走缓存，减少数据库压力，提高接口性能。

	if s.RDB != nil {
		if val, err := s.RDB.Get(ctx, cacheKey).Result(); err == nil && val != "" {
			// 直接反序列化
			var res pb.OrderDetailReply
			if err := json.Unmarshal([]byte(val), &res); err == nil {
				// 权限二次检查（缓存命中也要防越权）
				if role == "courier" && res.CourierId != courierID {
					return nil, status.Error(codes.PermissionDenied, "无权访问该订单详情")
				}
				return &res, nil
			}
		}
	}

	// 读库回源：精确检索订单
	var order model.SfOrders
	db := config.DB.WithContext(ctx)
	if req.OrderNo != "" {
		if err := db.Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Error(codes.NotFound, "订单不存在")
			}
			return nil, status.Errorf(codes.Internal, "查询订单失败: %v", err)
		}
	} else {
		if err := db.Where("id = ?", req.OrderId).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Error(codes.NotFound, "订单不存在")
			}
			return nil, status.Errorf(codes.Internal, "查询订单失败: %v", err)
		}
	}

	// 权限隔离：快递员仅能查看其订单
	if role == "courier" && order.CourierId != 0 && order.CourierId != courierID {
		return nil, status.Error(codes.PermissionDenied, "无权访问该订单详情")
	}

	// 订单对应任务（用于时间点推导）
	var task model.SfCourierTasks
	_ = db.Where("order_id = ?", order.Id).First(&task).Error

	// 审计日志（状态变更记录）
	var logs []model.SfOrderAuditLogs
	_ = db.Where("order_id = ?", order.Id).Order("created_at ASC").Find(&logs).Error

	// 改派链路
	var rs []model.SfOrderReassignments
	_ = db.Where("order_id = ?", order.Id).Order("created_at ASC").Find(&rs).Error

	// 构造响应
	resp := &pb.OrderDetailReply{
		Id:              order.Id,
		OrderNo:         order.OrderNo,
		UserId:          order.UserId,
		CourierId:       order.CourierId,
		SenderName:      order.SenderName,
		SenderPhone:     order.SenderPhone,
		SenderAddress:   order.SenderAddress,
		ReceiverName:    order.ReceiverName,
		ReceiverPhone:   order.ReceiverPhone,
		ReceiverAddress: order.ReceiverAddress,
		ServiceType:     order.ServiceType,
		ProductType:     order.ProductType,
		OrderStatus:     order.OrderStatus,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
		Fee: &pb.FeeDetail{
			EstimatedFee: order.EstimatedFee,
			ActualFee:    order.ActualFee,
			// 计费重量：对齐计费口径，采用体积重量与实际称重取较大值
			ChargeWeight: func() float64 {
				vw := order.VolumeWeight
				aw := order.ActualWeight
				if vw >= aw {
					return vw
				}
				return aw
			}(),
			PaymentMethod: order.PaymentMethod,
			PaymentStatus: order.PaymentStatus,
		},
	}

	// 司机接单与到达时间点推导
	// 接单时间：若任务存在且状态为 accepted，则使用任务 UpdatedAt；否则若订单状态 accepted，用订单 UpdatedAt
	if task.TaskStatus == "accepted" {
		resp.AcceptedAt = task.UpdatedAt.Format("2006-01-02 15:04:05")
	} else if order.OrderStatus == "accepted" {
		resp.AcceptedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	if !order.ActualPickupTime.IsZero() {
		resp.PickupActualAt = order.ActualPickupTime.Format("2006-01-02 15:04:05")
	}
	if !order.ActualDeliveryTime.IsZero() {
		resp.DeliveryActualAt = order.ActualDeliveryTime.Format("2006-01-02 15:04:05")
	}

	// 审计记录映射
	resp.StatusChanges = make([]*pb.AuditLogItem, 0, len(logs))
	for _, l := range logs {
		resp.StatusChanges = append(resp.StatusChanges, &pb.AuditLogItem{
			ActionType:   l.ActionType,
			OperatorId:   l.OperatorId,
			OperatorRole: l.OperatorRole,
			BeforeStatus: l.BeforeStatus,
			AfterStatus:  l.AfterStatus,
			ReasonType:   l.ReasonType,
			ReasonText:   l.ReasonText,
			CreatedAt:    l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 改派链路映射
	resp.Reassignments = make([]*pb.ReassignmentItem, 0, len(rs))
	for _, x := range rs {
		resp.Reassignments = append(resp.Reassignments, &pb.ReassignmentItem{
			FromCourierId: x.FromCourierId,
			ToCourierId:   x.ToCourierId,
			ReasonText:    x.ReasonText,
			CreatedAt:     x.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 风控标记与异常码：简化策略
	// 若订单状态为 exception 或存在异常类型审计记录，则标记风险；异常码取最近一次异常的 ReasonType
	riskFlag := order.OrderStatus == "exception" // 如果订单当前状态就是异常，先立个 flag
	exceptionCode := ""                          // 准备存最近一次异常原因
	// 逆序扫审计日志，找到最后一次出现异常的记录
	for i := len(logs) - 1; i >= 0; i-- {
		if logs[i].ActionType == "exception" || logs[i].AfterStatus == "exception" {
			riskFlag = true                    // 只要有异常动作，就强制标风险
			exceptionCode = logs[i].ReasonType // 取异常原因代码
			break                              // 找到最新的就停
		}
	}
	resp.Risk = &pb.RiskInfo{RiskFlag: riskFlag, ExceptionCode: exceptionCode, Notes: ""}

	// 写缓存下次同样订单号再进来，直接走缓存，不再查库，60 秒后自动过期，防止脏数据长期残留。
	if s.RDB != nil {
		if b, err := json.Marshal(resp); err == nil {
			_ = s.RDB.Set(ctx, cacheKey, string(b), time.Minute).Err() // TTL 60s
		}
	}

	// 可选：慢查询 explain 审核（仅日志输出，不影响功能）
	// 仅对 order_no 精确检索进行 explain
	if req.OrderNo != "" {
		type explainRow struct {
			id            int
			select_type   string
			table         string
			type_         string `gorm:"column:type"`
			possible_keys string
			key           string
			rows          int
			extra         string
		}
		// 忽略错误，避免影响主流程
		var e []explainRow
		_ = config.DB.Raw("EXPLAIN SELECT * FROM sf_orders WHERE order_no = ? LIMIT 1", req.OrderNo).Scan(&e).Error
	}

	return resp, nil
}

// convertTask 数据库模型 -> protobuf 消息
func convertTask(src *model.SfCourierTasks) *pb.Task {
	var taskStatusMap = map[string]pb.TaskStatus{
		"pending":   pb.TaskStatus_Pending,
		"accepted":  pb.TaskStatus_Accepted,
		"completed": pb.TaskStatus_Completed,
		"cancelled": pb.TaskStatus_Cancelled,
	}
	var taskTypeMap = map[string]pb.TaskType{
		"pickup":   pb.TaskType_Pickup,
		"delivery": pb.TaskType_Delivery,
	}
	return &pb.Task{
		Id:         src.Id,
		CourierId:  src.CourierId,
		TaskType:   taskTypeMap[src.TaskType], // 如果库里是 string 先用 map 转
		Priority:   int32(src.Priority),
		TaskStatus: taskStatusMap[src.TaskStatus],
		// 其他字段按需继续补
	}
}

// 快递员接任务列表
func (s *KuaiService) TakeTask(ctx context.Context, req *pb.TakeTaskRequest) (*pb.TakeTaskReply, error) {
	//从队里弹出一个任务id，原子操作
	taskID, err := queue.PopTask(s.RDB)
	if err == redis.Nil {
		return nil, status.Error(codes.NotFound, "暂无可接任务")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	//乐观锁只改“待抢”状态，防止重复抢
	res := config.DB.Model(&model.SfCourierTasks{}).
		Where("id = ? AND task_status = ?", taskID, "pending").
		Updates(map[string]interface{}{
			"courier_id":  req.CourierId,
			"task_status": "accepted",
		})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.AlreadyExists, "手慢啦，任务已被接走")
	}

	var task model.SfCourierTasks
	config.DB.First(&task, taskID)
	return &pb.TakeTaskReply{Task: convertTask(&task)}, nil
}

// 快递员绩效排行榜
func (s *KuaiService) Performance(ctx context.Context, req *pb.PerformanceRequest) (*pb.PerformanceReply, error) {
	db := config.DB
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var courierRanks []struct {
		CourierID   int64   `json:"courier_id"`
		CourierName string  `json:"courier_name"`
		TotalFee    float64 `json:"total_fee"`
	}
	// 查询快递员绩效排行榜
	err := db.Table("sf_courier_tasks t").
		Select("t.courier_id, c.real_name as courier_name, SUM(o.actual_fee) as total_fee").
		Joins("LEFT JOIN sf_couriers c ON t.courier_id = c.id").
		Joins("LEFT JOIN sf_orders o ON t.order_id = o.id").
		Where("t.task_status = ? AND t.updated_at >= ?", "completed", startOfMonth).
		Group("t.courier_id, c.real_name").
		Order("total_fee DESC").
		Limit(int(req.TopN)).
		Scan(&courierRanks).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询绩效排行榜失败: %v", err)
	}

	var rankList []*pb.CourierRank
	for i, rank := range courierRanks {
		rankList = append(rankList, &pb.CourierRank{
			Rank:             int32(i + 1),
			CourierId:        rank.CourierID,
			CourierName:      rank.CourierName,
			PerformanceScore: rank.TotalFee,
		})
	}

	return &pb.PerformanceReply{
		RankList: rankList,
	}, nil
}

// HandleException 处理异常情况处理多种异常 ，包括用户取消、快递员取消、派送失败、地址错误等
func (s *KuaiService) HandleException(ctx context.Context, req *pb.HandleExceptionRequest) (*pb.HandleExceptionReply, error) {

	db := config.DB

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 根据订单ID查询订单信息
	var order model.SfOrders
	err := tx.Where("id = ?", req.OrderId).First(&order).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.HandleExceptionReply{
				Success: false,
				Message: "订单不存在",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "查询订单失败: %v", err)
	}

	// 2. 根据订单ID查询相关的快递员任务
	var task model.SfCourierTasks
	err = tx.Where("order_id = ?", req.OrderId).First(&task).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "查询任务失败: %v", err)
	}

	var newOrderStatus pb.OrderStatus
	var message string

	// 3. 根据异常类型处理不同的逻辑
	switch req.ExceptionType {
	case pb.ExceptionType_USER_CANCEL:
		// 用户取消订单处理逻辑
		if order.OrderStatus == "cancelled" {
			return &pb.HandleExceptionReply{
				Success: false,
				Message: "订单已经是取消状态",
			}, nil
		}

		// 判断任务状态来决定处理方式
		var requiresNegotiation bool
		var suggestedFee float64

		if task.Id != 0 {
			// 检查任务状态
			switch task.TaskStatus {
			case "pending":
				// 未出发：直接同意取消
				message = "订单取消成功"
				requiresNegotiation = false
				suggestedFee = 0
			case "accepted":
				// 已出发：需要协商空跑费
				message = "快递员已出发，需要协商空跑费"
				requiresNegotiation = true
				suggestedFee = 5.0 // 建议空跑费5元
			case "completed":
				return &pb.HandleExceptionReply{
					Success: false,
					Message: "任务已完成，无法取消",
				}, nil
			case "cancelled":
				return &pb.HandleExceptionReply{
					Success: false,
					Message: "任务已取消",
				}, nil
			}
		} else {
			// 没有相关任务，直接取消
			message = "订单取消成功"
			requiresNegotiation = false
			suggestedFee = 0
		}

		// 更新订单状态为已取消
		err = tx.Model(&order).Updates(map[string]interface{}{
			"order_status": "cancelled",
			"updated_at":   time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "更新订单状态失败: %v", err)
		}

		// 如果存在相关任务，也更新任务状态
		if task.Id != 0 {
			err = tx.Model(&task).Updates(map[string]interface{}{
				"task_status": "cancelled",
				"updated_at":  time.Now(),
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, status.Errorf(codes.Internal, "更新任务状态失败: %v", err)
			}
		}

		newOrderStatus = pb.OrderStatus_ORDER_STATUS_CANCELLED

		// 提交事务
		if err = tx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
		}

		return &pb.HandleExceptionReply{
			Success:             true,
			Message:             message,
			NewStatus:           newOrderStatus,
			RequiresNegotiation: requiresNegotiation,
			SuggestedFee:        suggestedFee,
		}, nil

	case pb.ExceptionType_COURIER_CANCEL:
		// 快递员取消订单
		message = "快递员取消订单"
		newOrderStatus = pb.OrderStatus_ORDER_STATUS_PENDING

		// 将订单状态重置为待接单
		err = tx.Model(&order).Updates(map[string]interface{}{
			"order_status": "pending",
			"courier_id":   nil,
			"updated_at":   time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "更新订单状态失败: %v", err)
		}

		// 取消相关任务
		if task.Id != 0 {
			err = tx.Model(&task).Updates(map[string]interface{}{
				"task_status": "cancelled",
				"updated_at":  time.Now(),
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, status.Errorf(codes.Internal, "更新任务状态失败: %v", err)
			}
		}

		// 提交事务
		if err = tx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
		}

		return &pb.HandleExceptionReply{
			Success:             true,
			Message:             message,
			NewStatus:           newOrderStatus,
			RequiresNegotiation: false,
			SuggestedFee:        0,
		}, nil

	case pb.ExceptionType_DELIVERY_FAILED, pb.ExceptionType_PICKUP_FAILED,
		pb.ExceptionType_ADDRESS_ERROR, pb.ExceptionType_CONTACT_FAILED:
		// 其他异常情况，标记为异常状态
		message = fmt.Sprintf("订单异常: %s", req.Reason)
		newOrderStatus = pb.OrderStatus_ORDER_STATUS_EXCEPTION

		err = tx.Model(&order).Updates(map[string]interface{}{
			"order_status": "exception",
			"updated_at":   time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "更新订单状态失败: %v", err)
		}

	default:
		tx.Rollback()
		return &pb.HandleExceptionReply{
			Success: false,
			Message: "不支持的异常类型",
		}, nil
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}

	return &pb.HandleExceptionReply{
		Success:             true,
		Message:             message,
		NewStatus:           newOrderStatus,
		RequiresNegotiation: false, // 其他异常类型不需要协商
		SuggestedFee:        0,
	}, nil
}
