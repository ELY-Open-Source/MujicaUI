# 262 ToolCallGroup 工具调用组

## 用途

把连续的工具调用折叠成一组，标题行显示调用次数、总耗时与失败数。

## 最小示例

```go
open := false
calls := []agent.ToolCall{
	{ID: "a", Name: "list_dir", State: agent.AgentDone, Duration: 300 * time.Millisecond},
	{ID: "b", Name: "read_file", State: agent.AgentFailed, Error: "not found"},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	agent.ToolCallGroup(c, &open, calls, agent.ToolCallGroupOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ToolCallGroup(c *ui.Context, open *bool, calls []ToolCall, opts ToolCallGroupOptions) ChoiceResult`。`open` 为 nil、调用 `ID` 为空或重复时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | “N 次工具调用” | 替换内置标题。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- 组的状态：有调用执行中为执行中，否则有失败为失败，否则为完成。
- 失败数以文字“N 次失败”显示。
- 展开后每次调用用 ToolCallCard 显示。

## 事件

- `Changed()`：用户展开或折叠该组。

## 键盘操作

- Tab 到组标题，Enter / Space 展开或折叠；组内卡片各自可折叠。

## 限制

- 总耗时是各调用耗时之和，不扣除并行重叠。
