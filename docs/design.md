## 问题一：数据对象的实现
- 指路：[数据对象实现](..\internal\models)

## 问题二：数据库建立
- 指路：[数据库建立](..\internal\migrations)

## 问题三：Operation对象实现
- 指路：[Operation对象实现](..\internal\models\operation.go)
- 每次对文档的修改感觉更像是事件驱动，就是说我记录每次修改的操作，而不是每次记录这个文档修改之后的快照，这样更不冗余
- 服务器处理：参数校验——>检验是否在线是否有权限（没有怎么办，怎么处理？返回一个信号给前端吗）——>获取文档（拿取最新的文档快照和这个之后的操作日志，然后进行更新到最新的版本）——>构造Operation——>应用apply（将op传进去，根据type进行操作）——>持久化操作（保存这个operation）——>生成广播事件`DocumentEvent`结构体——>广播给其他用户`BroadcastToRoom`（根据文档ID找到房间ID，然后广播给其他用户）——>返回结果给用户`EditResponse`

## 问题四：如何对这个广播事件进行处理，对哪些用户进行广播
- 思路：首先根据文档Id找到room——根据excludeUserID之外的所有用户进行广播(每个client有一个send的channel，然后这个序列化之后的jsonEvt通过channel发送给对应的client)
TODO：好像还没有写把这个传回给前端的代码，maybe是websocket吗？
TODO：传输jsonEvt的步骤里面阻塞问题`Broadcast`

## 问题五：考虑消息发出过后不重复、不断档、不冲突
- 1.不重复：每个operation的id是唯一的（ai：uuid全球唯一），在`apply`函数里面检查`(doc_id,version)`这个唯一约束，防止重复插入
- <span style="color: yellow;">优化</span>：在`HandleEdit`中，可以先检查`storage.OperationExists(op.ID)`，如果存在则返回错误，防止重复插入

- 不冲突：`Apply`中检查`op.BaseVersion != d.Version`时返回错误，客户端刷新后重试
- <span style="color: yellow;">优化</span>：实现OT(ai提的建议)，不过感觉思路很简单：就是取这两个之间的所有操作，依次调整，最后转换（还没完成）

- 不断档：在`broadcast`中的最后给客户端穿这个jsonEvt的时候，某个客户端阻塞的两个解决办法（未完成）；客户端定期发送最后收到的Version，服务端可重发丢失消息（？？？？这个ACK哪里实现了）；断线重连：客户端重连时携带lastVersion，服务端返回缺失操作（？？？？？感觉这个我也没实现是吗？？？）
- <span style="color: yellow;">优化</span>：增加异步重试队列，避免因临时存储故障而丢失操作

## 问题六：服务器重启或节点故障时如何恢复文档状态
- 就是`getDocument`函数，根据文档ID获取文档快照和操作日志，然后进行更新到最新的版本
- <span style="color: yellow;">优化</span>：定期保存快照，现在似乎是每次都保存了

## challenge：分布式
- 分布式需要关心“数据在哪”，“如何同步”，“节点故障怎么办”
- 1.如何将多个文档


