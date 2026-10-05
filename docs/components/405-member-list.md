# 405 MemberList 成员列表

## 用途

将成员分成在线、离线两组，展示头像、状态、角色并提供选择与打开事件。

## 最小示例

```go
var s messaging.MemberListState
members := []messaging.Member{
	{ID: "1", Name: "Sakiko", Presence: messaging.PresenceOnline, Role: "Leader"},
	{ID: "2", Name: "Umiri", Presence: messaging.PresenceOffline},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.MemberList(c, &s, members, messaging.MemberListOptions{}); r.Opened != "" {
		_ = r.Opened // show that member's profile
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`MemberListState` 保存 `Selected` 和滚动位置；`members []Member` 的 ID 唯一、姓名必填。`Image`、`Presence`、`Role` 控制展示。`Height` 默认 360 DIP；`Disabled` 默认 false。修改在线分组时替换成员切片，缓存仅在切片变化时重建。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Changed()` 表示选中成员变化；`Opened` 返回 Enter 或双击打开的成员 ID。

## 键盘操作

Tab 聚焦列表；上下方向键选择；Home / End 首尾；Enter 或双击打开成员。

## 限制

分组缓存依据切片身份与长度；成员数据变化后替换切片。
