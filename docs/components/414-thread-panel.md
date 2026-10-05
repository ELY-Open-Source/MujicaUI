# 414 ThreadPanel 话题侧栏

## 用途

在侧栏展示原消息与虚拟化回复列表，提供回复输入、发送和关闭。

## 最小示例

```go
open := true
root := mujica.ChatEntry{ID: "r", Author: "Sakiko", Text: "Set list for Friday?"}
var replies []mujica.ChatEntry
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.ThreadPanel(c, &open, root, replies, mujica.ThreadPanelOptions{}); r.Sent != "" {
		replies = append(replies, mujica.ChatEntry{ID: r.Sent, Author: "Me", Text: r.Sent, Own: true})
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`open *bool` 控制显示；`root ChatEntry` 提供原消息，`replies []ChatEntry` 提供回复。`Width` 默认 360 DIP并受窗口宽度约束；`Disabled` 默认 false，禁用时关闭面板。回复 ID 必须稳定唯一。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Sent` 返回本帧发送的非空回复，发送后清空输入；关闭时设置 `*open=false`。

## 键盘操作

Tab 访问关闭、回复和发送；Enter 发送回复；Escape 关闭面板。

## 限制

回复发送和持久化由调用方完成；原消息区限高滚动，回复列表独立滚动。
