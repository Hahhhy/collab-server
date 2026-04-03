package models

import "time"

// Presence represents the presence information of a user in a collaborative editing session.
//对于每个用户，记录他们的光标位置和选择范围等信息，以便其他用户能够看到他们的编辑活动。
//是从用户操作角度的model

//实时状态对象的简单实现

type UserRole string

//三种用户状态：编辑者、评论者和查看者
const (
	RoleEditor    UserRole = "editor"
	RoleCommenter UserRole = "commenter"
	RoleViewer    UserRole = "viewer"
)

//连接状态
type ConnectionStatus string

const (
	StatusConnected    ConnectionStatus = "online"
	StatusDisconnected ConnectionStatus = "offline"
)

//光标位置
//暂时分这两种状态？？？？
type Cursor struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

//文本选区
//用开头和结尾的光标位置来表示文本选区
type Selection struct {
	Start Cursor `json:"start"`
	End   Cursor `json:"end"`
}

//判断是否处于输入状态
type TypingStatus struct {
	IsTyping  bool      `json:"is_typing"`
	StartedAt time.Time `json:"started_at"`
}

//用户的实时状态信息
type Presence struct {
	UserID string           `json:"user_id"`
	Role   UserRole         `json:"role"`
	DocID  string           `json:"doc_id"`
	State  ConnectionStatus `json:"state"`
	//这三个为什么要指针？因为有可能没有光标位置、选区或者输入状态，所以用指针来表示可选的字段
	Cursor    *Cursor       `json:"cursor"`
	Selection *Selection    `json:"selection"`
	Typing    *TypingStatus `json:"typing"`
	//这两个感觉可能没必要？？？？？
	LastSeenAt time.Time `json:"last_seen_at"`
	JoinedAt   time.Time `json:"joined_at"`
}
