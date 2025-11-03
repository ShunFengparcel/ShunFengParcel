package service

import (
	"ShunFengParcel/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// NotificationType 通知类型
type NotificationType string

const (
	// 通知类型常量
	NotifyTypePayment   NotificationType = "payment"   // 支付通知
	NotifyTypeOrder     NotificationType = "order"     // 订单通知
	NotifyTypeSystem    NotificationType = "system"    // 系统通知
	NotifyTypeDelivery  NotificationType = "delivery"  // 配送通知
	NotifyTypeRefund    NotificationType = "refund"    // 退款通知
	NotifyTypeHeartbeat NotificationType = "heartbeat" // 心跳
)

// NotificationMessage 通知消息结构
type NotificationMessage struct {
	Type      NotificationType       `json:"type"`           // 通知类型
	Title     string                 `json:"title"`          // 通知标题
	Content   string                 `json:"content"`        // 通知内容
	Data      map[string]interface{} `json:"data,omitempty"` // 附加数据
	Timestamp int64                  `json:"timestamp"`      // 时间戳
	MessageID string                 `json:"message_id"`     // 消息ID
	Read      bool                   `json:"read"`           // 是否已读
}

// NotificationClient WebSocket客户端
type NotificationClient struct {
	ID       string                    // 客户端唯一ID（用户ID）
	Conn     *websocket.Conn           // WebSocket连接
	Send     chan *NotificationMessage // 发送消息通道
	Manager  *NotificationManager      // 管理器引用
	LastPing time.Time                 // 最后心跳时间
	mu       sync.Mutex                // 互斥锁
}

// NotificationManager 通知管理器
type NotificationManager struct {
	clients    map[string]*NotificationClient // 所有连接的客户端
	broadcast  chan *BroadcastMsg             // 广播消息通道
	register   chan *NotificationClient       // 注册客户端通道
	unregister chan *NotificationClient       // 注销客户端通道
	mu         sync.RWMutex                   // 读写锁
}

// BroadcastMsg 广播消息
type BroadcastMsg struct {
	UserIDs []string             // 目标用户ID列表（空表示广播给所有人）
	Message *NotificationMessage // 消息内容
}

// 全局通知管理器
var notificationManager *NotificationManager

// WebSocket升级器
var notificationUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查Origin
	},
	HandshakeTimeout: 10 * time.Second,
}

// InitNotificationSystem 初始化通知系统
func InitNotificationSystem() {
	notificationManager = NewNotificationManager()
	go notificationManager.Start()
	utils.Info("✅ WebSocket通知系统已初始化")
}

// NewNotificationManager 创建通知管理器
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		clients:    make(map[string]*NotificationClient),
		broadcast:  make(chan *BroadcastMsg, 256),
		register:   make(chan *NotificationClient, 64),
		unregister: make(chan *NotificationClient, 64),
	}
}

// Start 启动通知管理器
func (m *NotificationManager) Start() {
	utils.Info("📡 通知管理器已启动")

	// 启动定时清理任务
	go m.cleanInactiveClients()

	for {
		select {
		case client := <-m.register:
			m.registerClient(client)

		case client := <-m.unregister:
			m.unregisterClient(client)

		case broadcast := <-m.broadcast:
			m.broadcastMessage(broadcast)
		}
	}
}

// registerClient 注册客户端
func (m *NotificationManager) registerClient(client *NotificationClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已存在同ID的连接，先关闭旧连接
	if oldClient, exists := m.clients[client.ID]; exists {
		close(oldClient.Send)
		oldClient.Conn.Close()
		utils.Warnf("用户 %s 重新连接，关闭旧连接", client.ID)
	}

	m.clients[client.ID] = client
	utils.Infof("✅ 客户端已连接: UserID=%s, 当前在线=%d", client.ID, len(m.clients))

	// 发送欢迎消息
	welcomeMsg := &NotificationMessage{
		Type:      NotifyTypeSystem,
		Title:     "连接成功",
		Content:   "您已成功连接到顺丰速递通知系统",
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}
	client.Send <- welcomeMsg
}

// unregisterClient 注销客户端
func (m *NotificationManager) unregisterClient(client *NotificationClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[client.ID]; exists {
		delete(m.clients, client.ID)
		close(client.Send)
		utils.Infof("❌ 客户端已断开: UserID=%s, 当前在线=%d", client.ID, len(m.clients))
	}
}

// broadcastMessage 广播消息
func (m *NotificationManager) broadcastMessage(broadcast *BroadcastMsg) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果目标用户列表为空，广播给所有人
	if len(broadcast.UserIDs) == 0 {
		sentCount := 0
		for _, client := range m.clients {
			select {
			case client.Send <- broadcast.Message:
				sentCount++
			default:
				utils.Warnf("客户端 %s 消息通道已满，跳过", client.ID)
			}
		}
		utils.Infof("📢 广播消息给 %d 个用户: %s", sentCount, broadcast.Message.Title)
		return
	}

	// 发送给指定用户
	sentCount := 0
	for _, userID := range broadcast.UserIDs {
		if client, exists := m.clients[userID]; exists {
			select {
			case client.Send <- broadcast.Message:
				sentCount++
			default:
				utils.Warnf("客户端 %s 消息通道已满，跳过", userID)
			}
		}
	}
	utils.Infof("📤 发送消息给 %d/%d 个用户: %s", sentCount, len(broadcast.UserIDs), broadcast.Message.Title)
}

// cleanInactiveClients 清理不活跃的客户端
func (m *NotificationManager) cleanInactiveClients() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, client := range m.clients {
			// 超过2分钟没有心跳，视为离线
			if now.Sub(client.LastPing) > 2*time.Minute {
				utils.Warnf("客户端 %s 超时未心跳，强制断开", id)
				client.Conn.Close()
				delete(m.clients, id)
				close(client.Send)
			}
		}
		m.mu.Unlock()
	}
}

// HandleNotificationWebSocket WebSocket连接处理器
func HandleNotificationWebSocket(ctx kratoshttp.Context) error {
	// 从查询参数获取用户ID
	userID := ctx.Request().URL.Query().Get("user_id")
	if userID == "" {
		utils.Error("WebSocket连接缺少user_id参数")
		return fmt.Errorf("缺少user_id参数")
	}

	// 升级为WebSocket连接
	conn, err := notificationUpgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		utils.LogAPIError("WebSocket", "Upgrade", 500, err)
		return err
	}

	// 创建客户端
	client := &NotificationClient{
		ID:       userID,
		Conn:     conn,
		Send:     make(chan *NotificationMessage, 256),
		Manager:  notificationManager,
		LastPing: time.Now(),
	}

	// 注册客户端
	notificationManager.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()

	return nil
}

// readPump 读取客户端消息
func (c *NotificationClient) readPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Conn.Close()
	}()

	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		c.mu.Lock()
		c.LastPing = time.Now()
		c.mu.Unlock()
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Errorf("WebSocket读取错误: %v", err)
			}
			break
		}

		// 处理客户端消息（如心跳、已读确认等）
		c.handleClientMessage(message)
	}
}

// writePump 向客户端发送消息
func (c *NotificationClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送JSON消息
			data, err := json.Marshal(message)
			if err != nil {
				utils.Error("序列化消息失败", zap.Error(err))
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				utils.Error("发送消息失败", zap.String("user_id", c.ID), zap.Error(err))
				return
			}

			utils.Debugf("💬 发送消息给 %s: %s", c.ID, message.Title)

		case <-ticker.C:
			// 发送心跳
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleClientMessage 处理客户端消息
func (c *NotificationClient) handleClientMessage(data []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "pong":
		// 心跳响应
		c.mu.Lock()
		c.LastPing = time.Now()
		c.mu.Unlock()

	case "read":
		// 标记消息已读
		if messageID, ok := msg["message_id"].(string); ok {
			utils.Debugf("用户 %s 已读消息: %s", c.ID, messageID)
		}
	}
}

// generateMsgID 生成消息ID
func generateMsgID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// --- 业务通知函数 ---

// SendPaymentNotification 发送支付通知
func SendPaymentNotification(userID, orderNo string, amount float64, status string) {
	if notificationManager == nil {
		utils.Warn("通知系统未初始化")
		return
	}

	msg := &NotificationMessage{
		Type:    NotifyTypePayment,
		Title:   "支付通知",
		Content: fmt.Sprintf("您的订单 %s 支付%s，金额: ¥%.2f", orderNo, status, amount),
		Data: map[string]interface{}{
			"order_no": orderNo,
			"amount":   amount,
			"status":   status,
		},
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}

	notificationManager.broadcast <- &BroadcastMsg{
		UserIDs: []string{userID},
		Message: msg,
	}

	utils.LogBusinessInfo("通知系统", "发送支付通知", map[string]interface{}{
		"user_id":  userID,
		"order_no": orderNo,
		"amount":   amount,
	})
}

// SendOrderNotification 发送订单通知
func SendOrderNotification(userID, orderNo, status, message string) {
	if notificationManager == nil {
		utils.Warn("通知系统未初始化")
		return
	}

	msg := &NotificationMessage{
		Type:    NotifyTypeOrder,
		Title:   "订单通知",
		Content: fmt.Sprintf("订单 %s %s", orderNo, message),
		Data: map[string]interface{}{
			"order_no": orderNo,
			"status":   status,
		},
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}

	notificationManager.broadcast <- &BroadcastMsg{
		UserIDs: []string{userID},
		Message: msg,
	}
}

// SendDeliveryNotification 发送配送通知
func SendDeliveryNotification(userID, orderNo, courierName, phone, status string) {
	if notificationManager == nil {
		utils.Warn("通知系统未初始化")
		return
	}

	msg := &NotificationMessage{
		Type:    NotifyTypeDelivery,
		Title:   "配送通知",
		Content: fmt.Sprintf("您的订单 %s 正在配送中，快递员：%s", orderNo, courierName),
		Data: map[string]interface{}{
			"order_no":      orderNo,
			"courier_name":  courierName,
			"courier_phone": phone,
			"status":        status,
		},
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}

	notificationManager.broadcast <- &BroadcastMsg{
		UserIDs: []string{userID},
		Message: msg,
	}
}

// SendSystemNotification 发送系统通知
func SendSystemNotification(userIDs []string, title, content string, data map[string]interface{}) {
	if notificationManager == nil {
		utils.Warn("通知系统未初始化")
		return
	}

	msg := &NotificationMessage{
		Type:      NotifyTypeSystem,
		Title:     title,
		Content:   content,
		Data:      data,
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}

	notificationManager.broadcast <- &BroadcastMsg{
		UserIDs: userIDs,
		Message: msg,
	}
}

// BroadcastToAll 广播给所有在线用户
func BroadcastToAll(title, content string) {
	if notificationManager == nil {
		utils.Warn("通知系统未初始化")
		return
	}

	msg := &NotificationMessage{
		Type:      NotifyTypeSystem,
		Title:     title,
		Content:   content,
		Timestamp: time.Now().Unix(),
		MessageID: generateMsgID(),
		Read:      false,
	}

	notificationManager.broadcast <- &BroadcastMsg{
		UserIDs: nil, // nil表示广播给所有人
		Message: msg,
	}

	utils.Infof("📢 系统广播: %s - %s", title, content)
}

// GetOnlineUsersCount 获取在线用户数
func GetOnlineUsersCount() int {
	if notificationManager == nil {
		return 0
	}

	notificationManager.mu.RLock()
	defer notificationManager.mu.RUnlock()
	return len(notificationManager.clients)
}

// IsUserOnline 检查用户是否在线
func IsUserOnline(userID string) bool {
	if notificationManager == nil {
		return false
	}

	notificationManager.mu.RLock()
	defer notificationManager.mu.RUnlock()
	_, exists := notificationManager.clients[userID]
	return exists
}
