# ShareDialog

## 用途

邀请成员、调整权限、移除成员与复制链接回调。

## 最小示例

```go
open := false
members := []mujica.ShareMember{{ID: "ada", Name: "Ada", Permission: mujica.PermissionEdit, Owner: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Button(c, "Share", mujica.ButtonOptions{}).Clicked() {
		open = true
	}
	mujica.ShareDialog(c, &open, &members, mujica.ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger"})
}
ui.Render(view, 600, 400, 1)
```

## 参数

Title 必填；Link 为空隐藏复制；OnCopy 可覆盖剪贴板动作；Disabled 或祖先禁用关闭弹层。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Invited 为邀请请求；Changed、Removed 已应用成员变化；Copied 为复制请求完成。

## 键盘操作

Tab 在对话框内移动；Enter 邀请或激活操作；Escape 关闭并恢复焦点。

## 限制

权限和成员由调用方持有；网络与链接有效性校验由业务负责。
