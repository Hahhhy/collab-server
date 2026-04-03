package service

import "fmt"

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

	return op, nil
}
