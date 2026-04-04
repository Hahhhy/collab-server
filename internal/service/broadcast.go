package service

import (
	"TO/internal/models"
	"encoding/json"
	"sync"
	"time"
)

//这个对象是用来告诉客户端其他人？？
// 其他像document就是对于我们服务端后面来说要储存的对象，有数据库的

type DocumentEventType string

const (
	EventDocumentEdited DocumentEventType = "document_edited"
	EventUserJoined     DocumentEventType = "user_joined"
	EventUserLeft       DocumentEventType = "user_left"
	EventCursorMoved    DocumentEventType = "cursor_moved"
	EventTypingStarted  DocumentEventType = "typing_started"
	EventTypingStopped  DocumentEventType = "typing_stopped"
)

type DocumentEvent struct {
	Type       DocumentEventType `json:"type"`
	DocID      string            `json:"doc_id"`
	Version    int               `json:"version,omitempty"`
	Op         *models.Operation `json:"op,omitempty"`
	OccurredAt time.Time         `json:"occurred_at"`
}

//Client表示一个在线连接
//这里先不接websocket.Conn，只保留send队列
//就是说先不完成双向交流，显示单向发送，send

type Client struct {
	//UserID和DocID是为了后续广播消息时知道发给谁和哪个文档的用户
	UserID string
	DocID  string
	Send   chan []byte
}

// DocRoom表示一个文档房间
type DocRoom struct {
	DocID   string
	mu      sync.RWMutex
	clients map[string]*Client // userID -> Client,来识别？？？？？？这是什么意思呢🤔
}

// Broadcast
func (r *DocRoom) Broadcast(evt DocumentEvent, excludeUserID string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	jsonEvt, err := json.Marshal(evt)
	if err != nil {
		// 处理错误
		return
	}
	//遍历房间里的所有客户端，给除了excludeUserID以外的客户端发送消息
	for userID, client := range r.clients {
		if userID != excludeUserID {
			client.Send <- jsonEvt //这个有可能失败吗？？？？这里怎么写比较好呢，如果缓冲区满了呢？？？
			//对于这么多客户端，如果有一个满了的话比如说网络不好在这里阻塞，
			// 会影响后面的，所以可能要改进一下？？
			//可以直接丢掉？？select+default的方式？？？？
			//或者goroutine
		}
	}
}
