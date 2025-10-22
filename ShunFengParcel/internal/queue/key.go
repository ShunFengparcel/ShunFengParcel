package queue

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Redis 队列键
const (
	KeyPending = "task:pending" // P0: 超时订单（最高优先级）
	KeyNormal  = "task:normal"  // P1: 预约时间即将到期的订单（高优先级）
	KeyLow     = "task:low"     // P2: 普通已接单订单（中优先级）
	KeyLowest  = "task:lowest"  // P3: 新订单（低优先级）
)

// KeyByPriority 根据优先级返回对应 Redis 键
// P0（最高）：超时订单
// P1（高）：预约时间即将到期的订单
// P2（中）：普通已接单订单
// P3（低）：新订单
func KeyByPriority(p int8) string {
	switch p {
	case 0:
		return KeyPending  // P0: 超时订单
	case 1:
		return KeyNormal   // P1: 预约时间即将到期
	case 2:
		return KeyLow      // P2: 普通已接单订单
	default:
		return KeyLowest   // P3: 新订单
	}
}

// queue/redis_queue.go
func PopTask(rdb redis.UniversalClient) (int64, error) {
	//从高到低依次弹出：P0(超时) -> P1(即将到期) -> P2(已接单) -> P3(新订单)
	keys := []string{"task:pending", "task:normal", "task:low", "task:lowest"}
	for _, k := range keys {
		// 原子弹出最早 1 个元素
		vals, err := rdb.ZPopMin(context.Background(), k, 1).Result()
		if err == nil && len(vals) > 0 {
			// 安全的类型转换
			switch member := vals[0].Member.(type) {
			case int64:
				return member, nil
			case string:
				// 如果是字符串，尝试转换为 int64
				if taskID, err := strconv.ParseInt(member, 10, 64); err == nil {
					return taskID, nil
				}
			case float64:
				// 如果是 float64，转换为 int64
				return int64(member), nil
			}
		}
	}
	return 0, redis.Nil
}
