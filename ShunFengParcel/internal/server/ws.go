package server

import (
    "context"
    "net/http"
    "sync"
    "time"

    "ShunFengParcel/internal/basic/config"

    kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/gorilla/websocket"
)

// WSHub 管理订单订阅的 WebSocket 连接，并复用 Redis PubSub 通道
type WSHub struct {
    mu     sync.RWMutex
    rooms  map[string]map[*websocket.Conn]struct{}
    subs   map[string]context.CancelFunc // 每个订单一个 Redis 订阅 goroutine
    up     websocket.Upgrader
}

func NewWSHub() *WSHub {
    return &WSHub{
        rooms: make(map[string]map[*websocket.Conn]struct{}),
        subs:  make(map[string]context.CancelFunc),
        up: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool { return true }, // CORS 放行，生产建议白名单
        },
    }
}

// RegisterRoutes 注册 /ws/order/{orderId}
func (h *WSHub) RegisterRoutes(s *kratoshttp.Server) {
    r := s.Route("/")
    r.GET("/ws/order/{orderId}", func(ctx kratoshttp.Context) error {
        // 升级 WebSocket
        w := ctx.Response()
        req := ctx.Request()
        conn, err := h.up.Upgrade(w, req, nil)
        if err != nil {
            return err
        }

        vars := ctx.Vars()["orderId"]
        if len(vars) == 0 {
            return nil // 没有 orderId 参数
        }
        orderID := vars[0]
        h.addConn(orderID, conn)
        h.ensureSub(ctx, orderID)

        // 心跳与读循环（读丢弃，仅用于检测断开）
        conn.SetReadLimit(512)
        _ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        conn.SetPongHandler(func(string) error {
            _ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
            return nil
        })

        go func() {
            ticker := time.NewTicker(30 * time.Second)
            defer ticker.Stop()
            for {
                if err := conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(5*time.Second)); err != nil {
                    break
                }
                <-ticker.C
            }
        }()

        // 读阻塞直到断开
        for {
            if _, _, err := conn.ReadMessage(); err != nil {
                break
            }
        }

        h.removeConn(orderID, conn)
        _ = conn.Close()
        return nil
    })
}

func (h *WSHub) addConn(orderID string, c *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if h.rooms[orderID] == nil {
        h.rooms[orderID] = make(map[*websocket.Conn]struct{})
    }
    h.rooms[orderID][c] = struct{}{}
}

func (h *WSHub) removeConn(orderID string, c *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if set := h.rooms[orderID]; set != nil {
        delete(set, c)
        if len(set) == 0 {
            delete(h.rooms, orderID)
            if cancel, ok := h.subs[orderID]; ok {
                cancel()
                delete(h.subs, orderID)
            }
        }
    }
}

// ensureSub 确保为指定订单启动一个 Redis 订阅
func (h *WSHub) ensureSub(ctx context.Context, orderID string) {
    h.mu.Lock()
    if _, ok := h.subs[orderID]; ok {
        h.mu.Unlock()
        return
    }
    // 启动订阅 goroutine
    cctx, cancel := context.WithCancel(context.Background())
    h.subs[orderID] = cancel
    h.mu.Unlock()

    go func() {
        // 若未配置 Redis，直接返回
        if config.RDB == nil {
            return
        }
        ch := "order:location:" + orderID
        sub := config.RDB.Subscribe(cctx, ch)
        defer sub.Close()
        for {
            msg, err := sub.ReceiveMessage(cctx)
            if err != nil {
                return
            }
            // 广播给房间内所有连接
            h.mu.RLock()
            conns := h.rooms[orderID]
            for c := range conns {
                _ = c.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
            }
            h.mu.RUnlock()
        }
    }()
}