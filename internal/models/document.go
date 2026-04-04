package models

import (
	"fmt"
	"sync"
	"time"
)

// OperationType represents the type of an operation performed on a document, such as creating, updating, or deleting content.
//是从文档呈现内容角度的model

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

type DocumentSnapshot struct {
	DocID   string
	Content string
	Version int
}

func (d *Document) Apply(op Operation) (Operation, error) {
	//文档上锁
	d.mu.Lock()
	//解锁
	defer d.mu.Unlock()

	//一致性检验

	//版本不对直接崩溃退出来处理吗？？？？？？
	//修改的时候基于的版本和当前版本不一样————修改时基于的版本在业务逻辑哪一步进行更新？？？？？？
	if op.BaseVersion != d.Version {
		return Operation{}, fmt.Errorf("version conflict: expected %d, got %d", d.Version, op.BaseVersion)
	}

	//处理操作
	//返回err，err为nil就是成功，
	var err error
	switch op.Type {
	case OpInsert:
		//插入操作
		//根据op.Position在文档内容中插入op.Text
		//更新文档版本号和更新时间
		//应该这些操作要在一个函数里面封装
		err = d.applyInsert(op.Position, op.Text)
	case OpDelete:
		//删除操作
		//根据op.Position和op.Length删除文档内容
		//更新文档版本号和更新时间
		err = d.applyDelete(op.Position, op.Length)
	default:
		return Operation{}, fmt.Errorf("unknown operation type: %s", op.Type)
	}

	if err != nil {
		return Operation{}, err
	}

	//如果版本正确，则要改变版本，更新
	//更新服务端的版本号
	op.Version = d.Version + 1
	//更新文档版本号
	d.Version = op.Version
	d.UpdatedAt = time.Now()
	//time.Now还是UpdatedAt
	op.CreatedAt = d.UpdatedAt

	return op, nil
} //返回操作：操作里面包含版本号，文档ID，用户ID，操作类型，位置，文本内容，长度，基于版本，创建时间等信息

func (d *Document) applyDelete(pos int, length int) error {
	//判断删除的位置和长度是否合法
	//这个无效操作是怎么造成的？？感觉不应该是用户可以做出来的效果？？？，length小于零是传输有误？？？咋来的这种bug呢？？？？？？？？？？？？？？？？？？
	// ？？？？？？？？？？？？？？？？？？？？
	// ？？？？？？？？？？？？？？？？？？
	if pos < 0 || pos > len(d.Content) || length < 0 || pos+length > len(d.Content) {
		return fmt.Errorf("invalid delete operation: position %d, length %d", pos, length)
	}

	//删除操作
	//根据在哪个位置开始删除，删除多长，然后把这个前面的和后面的进行拼接，我只记录多长，显示在界面上？
	d.Content = d.Content[:pos] + d.Content[pos+length:]
	return nil

}

func (d *Document) applyInsert(pos int, text string) error {
	//判断插入的位置是否合法
	if pos < 0 || pos > len(d.Content) {
		return fmt.Errorf("invalid insert operation: position %d", pos)
	}
	//插入操作
	//根据在哪个位置插入，插入什么内容，然后把这个前面的和后面的进行拼接，我只记录插入的内容，显示在界面上？
	d.Content = d.Content[:pos] + text + d.Content[pos:]
	return nil
}

func (d *Document) Snapshot() DocumentSnapshot {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return DocumentSnapshot{
		DocID:   d.ID,
		Content: d.Content,
		Version: d.Version,
	}
}
