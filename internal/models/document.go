package models

import (
	"sync"
	"time"
)

// OperationType represents the type of an operation performed on a document, such as creating, updating, or deleting content.
//是从文档呈现内容角度的model

type OperationType string

const (
	OpInsert OperationType = "insert"
	OpDelete OperationType = "delete"
)

// 文档的属性
type Document struct {
	mu        sync.RWMutex
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Version   int       `json:"version"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt time.Time `json:"deleted_at"`
	// IsDeleted bool      `json:"is_deleted"`
}

// 操作
type Operation struct {
	//id是操作的唯一标识，doc_id是操作所属的文档id，
	// user_id是执行操作的用户id，type是操作类型（插入、删除等），
	// position是操作的位置，text是插入的文本内容，
	// length是删除的文本长度，base_version是客户端认为自己编辑时基于的版本，
	// version是服务端分配的新版本，created_at是操作的创建时间。
	ID       string `json:"id"`
	DocID    string `json:"doc_id"`
	UserID   string `json:"user_id"`
	Position int    `json:"position"`
	Text     string `json:"text"`
	Length   int    `json:"length"`
	//客户端认为自己编辑时基于的版本，有什么用？？
	BaseVersion int64 `json:"base_version"`
	//服务端分配的新版本
	Version int64 `json:"version"`

	//操作及定位操作的数据属性？？？？

	CreatedAt time.Time     `json:"created_at"`
	Type      OperationType `json:"type"`
}

type DocumentSnapshot struct {
	DocID   string
	Content string
	Version int64
}
