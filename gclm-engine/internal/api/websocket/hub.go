package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gclm/gclm-flow/gclm-engine/internal/logger"
)

// Hub WebSocket 连接中心
type Hub struct {
	// 订阅管理: taskID -> []connection
	subscriptions map[string]map[*Connection]bool
	mu            sync.RWMutex

	// 连接管理
	connections map[*Connection]string // connection -> taskID

	// 事件通道
	broadcast chan *Event

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc
}

// Event WebSocket 事件
type Event struct {
	Type    string      `json:"type"`
	TaskID  string      `json:"task_id,omitempty"`
	Data    interface{} `json:"data"`
}

// Connection WebSocket 连接
type Connection struct {
	id     string
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Event
	taskID string
	mu     sync.Mutex
}

// NewHub 创建 WebSocket Hub
func NewHub(ctx context.Context) *Hub {
	ctx, cancel := context.WithCancel(ctx)

	return &Hub{
		subscriptions: make(map[string]map[*Connection]bool),
		connections:   make(map[*Connection]string),
		broadcast:     make(chan *Event, 256),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run 启动 Hub 事件循环
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("WebSocket Hub stopped")
			return
		case event := <-h.broadcast:
			h.broadcastEvent(event)
		}
	}
}

// Subscribe 订阅任务更新
func (h *Hub) Subscribe(taskID string, conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 添加到订阅映射
	if h.subscriptions[taskID] == nil {
		h.subscriptions[taskID] = make(map[*Connection]bool)
	}
	h.subscriptions[taskID][conn] = true

	// 记录连接的任务ID
	h.connections[conn] = taskID
	conn.taskID = taskID

	logger.Info().
		Str("task_id", taskID).
		Str("conn_id", conn.id).
		Msg("WebSocket subscribed")
}

// Unsubscribe 取消订阅
func (h *Hub) Unsubscribe(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	taskID, exists := h.connections[conn]
	if !exists {
		return
	}

	// 从订阅映射中删除
	delete(h.subscriptions[taskID], conn)
	delete(h.connections, conn)

	// 关闭发送通道
	close(conn.send)

	logger.Info().
		Str("task_id", taskID).
		Str("conn_id", conn.id).
		Msg("WebSocket unsubscribed")
}

// Broadcast 广播事件到订阅者
func (h *Hub) Broadcast(event *Event) {
	select {
	case h.broadcast <- event:
	case <-h.ctx.Done():
	}
}

// broadcastEvent 内部广播实现
func (h *Hub) broadcastEvent(event *Event) {
	h.mu.RLock()
	targetTaskID := event.TaskID

	// 如果没有指定 taskID，广播到所有连接
	if targetTaskID == "" {
		for _, conns := range h.subscriptions {
			for conn := range conns {
				h.sendToConnection(conn, event)
			}
		}
		h.mu.RUnlock()
		return
	}

	// 广播到特定任务的订阅者
	conns, exists := h.subscriptions[targetTaskID]
	if !exists {
		h.mu.RUnlock()
		return
	}

	for conn := range conns {
		h.sendToConnection(conn, event)
	}
	h.mu.RUnlock()
}

// sendToConnection 发送事件到单个连接 (带超时)
func (h *Hub) sendToConnection(conn *Connection, event *Event) {
	select {
	case conn.send <- event:
	case <-time.After(5 * time.Second):
		// 发送超时，关闭连接
		h.Unsubscribe(conn)
		logger.Warn().
			Str("conn_id", conn.id).
			Msg("WebSocket send timeout, closing connection")
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Hub) HandleWebSocket(c *gin.Context) {
	taskID := c.Param("id")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境需限制）
	},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error().Err(err).Msg("WebSocket upgrade failed")
		return
	}

	// 创建连接
	connection := &Connection{
		id:   generateConnID(),
		hub:  h,
		conn: conn,
		send: make(chan *Event, 256),
	}

	// 订阅任务
	h.Subscribe(taskID, connection)

	// 启动读写 goroutine
	go connection.writePump()
	connection.readPump()

	logger.Info().
		Str("task_id", taskID).
		Str("conn_id", connection.id).
		Msg("WebSocket connection established")
}

// readPump 读取消息循环
func (c *Connection) readPump() {
	defer c.hub.Unsubscribe(c)

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Debug().Err(err).Msg("WebSocket read error")
			}
			break
		}

		// 处理接收到的消息（如果需要）
		_ = message
	}
}

// writePump 写入消息循环
func (c *Connection) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			// 发送 ping 消息保持连接
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}

		case event, ok := <-c.send:
			if !ok {
				return
			}

			// 发送事件消息
			data, err := marshalEvent(event)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal event")
				continue
			}

			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to send WebSocket message")
				return
			}
		}
	}
}

// marshalEvent 序列化事件
func marshalEvent(event *Event) ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type":"%s","task_id":"%s","data":%s}`,
		event.Type,
		event.TaskID,
		marshalData(event.Data),
	)), nil
}

// marshalData 简单的数据序列化
func marshalData(data interface{}) string {
	if data == nil {
		return "{}"
	}
	// 简化处理，实际应使用 json.Marshal
	return fmt.Sprintf("%v", data)
}

// generateConnID 生成连接ID
func generateConnID() string {
	return fmt.Sprintf("conn-%d", time.Now().UnixNano())
}
