# 自洽性修复方案：

## 持久化bug
- 在 storage/snapshot.go 中实现 `SaveOperation` 方法（执行 SQL INSERT）。<span style="color: green;">已完成</span>
- 在 collab.go 的 [HandleEdit](..\internal\service\collab.go) 中，取消注释并启用 `s.storage.SaveOperation(appliedOp)`。<span style="color: green;">已完成</span>
- 决策： 每次编辑都保存快照太慢。建议仅在操作累积到一定数量或文档关闭时调用 [SaveSnapshot](。。\internal\storage\snapshot.go)。只要操作日志（`Operations Log`）被持久化，重启后通过 [getDocument](..\internal\service\collab.go) 重放即可恢复状态。<span style="color: green;">已完成</span>


## 乐观锁bug
-  客户端需要处理这个错误（例如：获取最新内容，合并更改，重新提交）。你的 [EditResponse](..\internal\service\collab.go) 返回 `Accepted: false`，这是合理的。<span style="color: red;">未完成</span>



## main.go路由逻辑bug
<span style="color: green;">已完成</span>


## 状态同步bug
-  实现 `WebSocket` 连接建立时更新 `user_sessions` 状态的逻辑，否则 [checkUserOnline](..\internal\service\collab.go) 将始终返回 `false`（除非你在数据库里手动插入了数据）。
- 就是没有实现用户状态实时更新的功能，应该是要根据用户前端传过来的某些信息判断这个用户是不是在线
- 缺少`AddClient`/`RemoveClient`/`SetUserConnectionStatus`<span style="color: red;">未完成</span>