package service

import (
	pb "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/basic/config"
	"ShunFengParcel/internal/biz"
	"ShunFengParcel/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"ShunFengParcel/internal/queue"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

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
		OrderId:    orders.Id,  // 关联订单ID
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

func (s *KuaiService) Performance(ctx context.Context, req *pb.PerformanceRequest) (*pb.PerformanceReply, error) {
	db := config.DB

	var courierRanks []struct {
		CourierID   int64   `json:"courier_id"`
		CourierName string  `json:"courier_name"`
		TaskCount   int64   `json:"task_count"`
		Score       float64 `json:"score"`
	}

	// 查询快递员绩效排行榜
	err := db.Table("sf_courier_tasks t").
		Select("t.courier_id, c.name as courier_name, COUNT(t.id) as task_count, COUNT(t.id) * 10 as score").
		Joins("LEFT JOIN sf_couriers c ON t.courier_id = c.id").
		Where("t.task_status = ?", "completed").
		Group("t.courier_id, c.name").
		Order("score DESC").
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
			PerformanceScore: rank.Score,
		})
	}

	return &pb.PerformanceReply{
		RankList: rankList,
	}, nil
}

// HandleException 处理异常情况
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
			Success:              true,
			Message:              message,
			NewStatus:            newOrderStatus,
			RequiresNegotiation:  requiresNegotiation,
			SuggestedFee:         suggestedFee,
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
