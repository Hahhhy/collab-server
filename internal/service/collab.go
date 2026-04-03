package service

import (
	"fmt"
	"time"
)

// TODO: 这个函数是干嘛用的，检验版本对不对为什么要放这个函数里面？？？？？？
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

// 删除操作
func (d *Document) applyDelete(pos int, length int) error {
	//判断删除的位置和长度是否合法
	if pos < 0 || pos > len(d.Content) || length < 0 || pos+length > len(d.Content) {
		return fmt.Errorf("invalid delete operation: position %d, length %d", pos, length)
	}

	//删除操作
	//根据在哪个位置开始删除，删除多长，然后把这个前面的和后面的进行拼接，我只记录多长，显示在界面上？
	d.Content = d.Content[:pos] + d.Content[pos+length:]
	return nil

}

// 插入操作
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

// 前面是方法函数包装，后面就要写这个过程，利用到前面这些函数
// 这个s是什么东西，collabService哪里定义的
func (s *CollabService) HandleEdit(req EditRequest) (*EditResponse, error) {
	//1.基本参数校验
	if req.DocID == "" || req.UserID == "" {
		return nil, fmt.Errorf("invalid request: missing docID or userID")
	}

	//2.校验用户权限
	if !s.checkUserPermission(req.UserID, req.DocID, RoleEditor) {
		return &EditResponse{Accepted: false, Reason: "no permission"}, nil
	}

	// 3) 获取文档
	doc, err := s.getDocument(req.DocID)
	if err != nil {
		return nil, err
	}

	// 4) 构造操作
	op := Operation{
		ID:          uuid.New().String(),
		DocID:       req.DocID,
		UserID:      req.UserID,
		Type:        req.Type,
		Position:    req.Position,
		Text:        req.Text,
		Length:      req.Length,
		BaseVersion: req.BaseVersion,
		CreatedAt:   time.Now(),
	}

	// 5) 应用操作
	appliedOp, err := doc.Apply(op)
	if err != nil {
		return &EditResponse{Accepted: false, Reason: err.Error()}, nil
	}

	// 6) 持久化操作
	//前面前是生成了operation，现在就是真正的把这个operation保存到数据库里面去，保存成功了就可以广播了，保存失败了就记录错误日志但是不影响用户的编辑体验，因为内存状态已经更新了
	//哪一步骤是更新内存状态的？？？？？？前面哪一步啊？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？

	if err := s.storage.SaveOperation(appliedOp); err != nil {
		// 记录错误但继续，因为内存状态已更新
		s.logger.Error("failed to save operation", err)
	}

	// 7) 生成广播事件
	evt := &DocumentEvent{
		Type:       EventDocumentEdited,
		DocID:      appliedOp.DocID,
		Version:    appliedOp.Version,
		Op:         appliedOp,
		OccurredAt: time.Now(),
	}

	// 8) 广播给其他用户
	s.broadcaster.BroadcastToRoom(req.DocID, evt, req.UserID)

	// 9) 返回结果
	snapshot := doc.Snapshot()
	return &EditResponse{
		Accepted:     true,
		Reason:       "ok",
		AppliedOp:    &appliedOp,
		NewSnapshot:  &snapshot,
		BroadcastEvt: evt,
	}, nil
}
