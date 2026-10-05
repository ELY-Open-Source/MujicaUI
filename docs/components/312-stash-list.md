# 312 StashList 贮藏列表

## 用途

贮藏列表，每行提供应用、弹出与删除。

## 最小示例

```go
var state data.ListState[string]
stashes := []git.Stash{{Ref: "stash@{0}", Message: "WIP on main", Branch: "main"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.StashList(c, &state, stashes, git.StashListOptions{})
	r.Element.Height(200)
	if action, ref, ok := r.Action(); ok {
		_, _ = action, ref // run git stash apply, pop or drop
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func StashList(c *ui.Context, s *ListState[string], stashes []Stash, opts StashListOptions) GitListResult`

## 参数

`Stash`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Ref | string | 必填 | 如 stash@{0}；为空时 panic。 |
| Message | string | "" | 贮藏说明。 |
| Branch | string | "" | 创建时所在分支。 |

`StashListOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Status | DataStatus | DataReady | 加载、追加、失败状态。 |
| Error | string | "" | 失败说明。 |

给返回的 `Element` 一个高度。

## 状态

- Ref 为等宽 AccentText；三个操作为带名称与提示的图标按钮。

## 事件

- `Action()` 返回 `"apply"`、`"pop"` 或 `"drop"` 与 Ref。

## 键盘操作

- ↑/↓ 移动选择；操作按钮用 Tab 到达，Enter 或 Space 执行。

## 限制

- 删除不带确认；需要时由调用方在回调里确认。
