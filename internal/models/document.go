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
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	IsDeleted bool      `json:"is_deleted"`
}

// 操作
type Operation struct {
	//客户端认为自己编辑时基于的版本
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
