# 268 AgentProgress Agent 进度

## 用途

显示整体进度和当前动作文字；时长未知时用不确定模式。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	agent.AgentProgress(c, 0.4, agent.AgentProgressOptions{Action: "Reading ledger.md"})
}
ui.Render(view, 400, 300, 1)
```

签名：`func AgentProgress(c *ui.Context, value float64, opts AgentProgressOptions) ui.Element`。确定模式下 `value` 不在 0 到 1 时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | “进度” | 进度名称。 |
| `Action` | `string` | `""` | 当前动作，如“正在读取 ledger.md”。 |
| `Indeterminate` | `bool` | `false` | 不确定模式，忽略 value。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- 确定模式显示百分比；不确定模式显示流动条，减少动画时为静态条加“进行中”。

## 事件

- 无事件；当前动作以状态角色读出。

## 键盘操作

- 不获取焦点。

## 限制

- 进度值由调用方计算。
