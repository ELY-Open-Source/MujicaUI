# 408 StatusSetter 状态设置

## 用途

编辑自定义状态文字和表情，选择清除期限。

## 最小示例

```go
status := mujica.CustomStatus{Emoji: "🎻", Text: "Rehearsing"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.StatusSetter(c, &status, mujica.StatusSetterOptions{}).Changed() {
		_ = status // publish it; clear it once status.Until passes
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`CustomStatus` 保存 `Emoji`、`Text`、`Until`；零 `Until` 表示持续保留。`Location` 默认 `time.Local`；`Disabled` 默认 false。清除选项为持续保留、30 分钟、1 小时、4 小时、当天结束、本周结束。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Changed()` 表示表情、文字或清除期限改变。调用方载入已有期限时显示实际时间，选择持续保留可清除期限。

## 键盘操作

Tab 访问表情、文字、期限和清除按钮；选择器用方向键移动、Enter 确认、Escape 关闭。

## 限制

Until 作为业务期限返回，调用方负责届时清除状态；组件保持空闲。
