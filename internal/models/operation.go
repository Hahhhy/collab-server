package models

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OpInsert OperationType = "insert"
	OpDelete OperationType = "delete"
)

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
	BaseVersion int `json:"base_version"`
	//服务端分配的新版本
	Version int `json:"version"`

	//操作及定位操作的数据属性？？？？

	CreatedAt time.Time     `json:"created_at"`
	Type      OperationType `json:"type"`
}

func NewOperation(docID, userID string, opType OperationType, pos int, text string, length int, baseVersion int) Operation {
	return Operation{
		//这里使用uuid生成一个唯一的操作ID，确保每个操作都有一个独特的标识符，方便后续的追踪和管理。
		ID:          uuid.New().String(),
		DocID:       docID,
		UserID:      userID,
		Type:        opType,
		Position:    pos,
		Text:        text,
		Length:      length,
		BaseVersion: baseVersion,
		CreatedAt:   time.Now(),
	}
}
