# 267 AgentStatus Agent 状态

## 用途

显示 Agent 当前状态：规划中、执行中、等待、完成、失败，用图标形状和文字表达，不只靠颜色。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	agent.AgentStatus(c, agent.AgentWaiting, agent.AgentStatusOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func AgentStatus(c *ui.Context, state AgentState, opts AgentStatusOptions) ui.Element`。未知状态 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 内置状态词 | 替换显示的文字。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `AgentState`：`AgentPending`、`AgentPlanning`、`AgentRunning`、`AgentWaiting`、`AgentDone`、`AgentFailed`、`AgentSkipped`，越界 panic。状态同时由图标形状和文字表达，不只靠颜色。
- 执行中显示旋转的菱形；`Still` 或系统减少动画时改为静态图标。

## 事件

- 无事件；屏幕阅读器以状态角色读出文字。

## 键盘操作

- 不获取焦点。

## 限制

- 状态动画只在显示时重绘，隐藏后停止。
