# 406 UserProfileCard 用户资料卡

## 用途

集中展示用户头像、姓名、在线状态、角色、自定义状态与操作插槽。

## 最小示例

```go
profile := messaging.UserProfile{Name: "Sakiko Togawa", Presence: messaging.PresenceBusy, Role: "Keyboard", Status: "In rehearsal"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.UserProfileCard(c, profile, messaging.UserProfileCardOptions{Actions: func() {
		input.Button(c, "Message", input.ButtonOptions{})
	}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`UserProfile.Name` 必填；`Image`、`Presence`、`Role`、`Status`、`About` 提供展示内容。`UserProfileCardOptions.Width` 默认 300 DIP，必须 ≥ 0；`Actions func()` 提供可选操作。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

操作由 `Actions` 插槽中的控件返回。

## 键盘操作

Tab 访问操作插槽；Enter / Space 激活按钮。

## 限制

业务动作、头像加载与个人资料更新由调用方处理。
