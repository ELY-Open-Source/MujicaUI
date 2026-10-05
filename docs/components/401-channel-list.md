# 401 ChannelList 频道列表

## 用途

展示频道与私信、未读计数、静音和置顶标记，支持折叠分组与频道选择。

## 最小示例

```go
groups := []mujica.ChannelGroup{{ID: "ch", Title: "Channels", Channels: []mujica.Channel{
	{ID: "general", Name: "general", Unread: 3},
	{ID: "rehearsal", Name: "rehearsal", Pinned: true},
}}}
selected := "general"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ChannelList(c, &selected, groups, mujica.ChannelListOptions{}).Changed() {
		_ = selected // show that channel
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`selected *string` 保存频道 ID；`groups []ChannelGroup` 提供分组。分组 ID、频道 ID 必须唯一，名称必填，`Unread` ≥ 0。`Direct`、`Presence`、`Muted`、`Pinned` 控制展示。`Height` 默认 360 DIP；`Disabled` 默认 false；`Revision` 默认 0，原地修改分组或频道数据后递增以刷新缓存。替换分组切片也会刷新。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Changed()` 表示选择 ID 发生变化；分组展开状态在组件内保存。

## 键盘操作

Tab 进入；上下方向键切换行；Enter / Space 选择或折叠；分组标题上的左、右方向键折叠、展开。

## 限制

业务连接、静音设置与置顶持久化由调用方维护。分组排序与校验在数据或 Revision 变化时进行，可见行由 MyGo List 构建。
