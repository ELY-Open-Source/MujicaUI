# 284 TraceViewer 调用链查看器

## 用途

以树显示调用链跨度，每行带耗时与时间轴条；选中跨度在下方查看属性。

## 最小示例

```go
var state mujica.TraceViewerState
spans := []mujica.TraceSpan{
	{ID: "run", Name: "agent.run", End: 900 * time.Millisecond},
	{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond,
		Attrs: []mujica.TraceAttr{{Key: "model", Value: "opus"}}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TraceViewer(c, &state, spans, mujica.TraceViewerOptions{}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)
```

签名：`func TraceViewer(c *ui.Context, s *TraceViewerState, spans []TraceSpan, opts TraceViewerOptions) DataTableView`

## 参数

`TraceSpan`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| ID | string | 必填 | 唯一标识；重复 panic。 |
| Parent | string | "" | 父跨度；不存在时作为根。 |
| Name | string | "" | 名称。 |
| Start, End | time.Duration | 0 | 相对调用链开始的偏移；需 0 ≤ Start ≤ End。 |
| Failed | bool | false | 失败。 |
| Attrs | []TraceAttr | nil | 属性键值。 |

`TraceViewerOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Status | DataStatus | DataReady | 加载、失败状态。 |
| Error | string | "" | 失败说明。 |
| Version | int | 0 | 原地修改数据时递增；新切片或长度变化自动刷新。 |

`TraceViewerState.Tree` 保存展开、排序与选择。索引只在 spans 变化时重建；父链成环时 panic。给返回的 `Element` 一个高度。

## 状态

- 时间轴条按整条调用链的跨度缩放，位置吸附像素。
- 失败跨度：叉号图标、Danger 文字与空心时间条，不只靠颜色。
- 属性面板显示开始、耗时与全部属性。

## 事件

- 同 TreeTable：`Changed()` 选择变化。

## 键盘操作

- ↑/↓ 选择，→ 展开，← 收起；输入名称开头字母跳转。

## 限制

- 只构建可见行（TreeTable 虚拟化）；属性面板也按可见行虚拟化，独立滚动。
