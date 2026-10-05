# TransferQueue

## 用途

显示上传下载队列并回传暂停、继续、重试、移除。

## 最小示例

```go
items := []mujica.QueuedTransfer{{ID: "a", Name: "ledger.zip", Upload: true, Status: mujica.TransferActive, Done: 25, Total: 100}}
var state mujica.TransferQueueState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TransferQueue(c, items, &state, mujica.TransferQueueOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)
```

## 参数

ID 唯一且必填；Status 默认为排队；Revision 原地更新数据时递增；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Action 与 ID 返回对应传输操作。

## 键盘操作

Tab 在可见行操作间移动；Enter / Space 激活。

## 限制

网络由调用方执行；汇总仅在数据切片或 Revision 变化时重算。
