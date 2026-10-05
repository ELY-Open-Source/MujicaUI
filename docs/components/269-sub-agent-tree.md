# 269 SubAgentTree 子 Agent 树

## 用途

以树显示子 Agent 层级、各自任务与状态，可展开查看下级。

## 最小示例

```go
agents := []agent.SubAgent{{ID: "lead", Name: "Steward", State: agent.AgentRunning, Children: []agent.SubAgent{
	{ID: "cook", Name: "Cook", Task: "Plan the menu", State: agent.AgentWaiting},
}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	agent.SubAgentTree(c, agents, agent.SubAgentTreeOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func SubAgentTree(c *ui.Context, roots []SubAgent, opts SubAgentTreeOptions) ChoiceResult`。`ID` 为空或重复、`Name` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | `false` | 禁止展开与折叠。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `SubAgent` 字段：`ID`、`Name`、`Task`、`State`、`Children`。
- `AgentState`：`AgentPending`、`AgentPlanning`、`AgentRunning`、`AgentWaiting`、`AgentDone`、`AgentFailed`、`AgentSkipped`，越界 panic。状态同时由图标形状和文字表达，不只靠颜色。
- 有下级的行显示箭头和“N 个子智能体”。

## 事件

- `Changed()`：用户展开或折叠了某个 Agent。

## 键盘操作

- Tab 进入树；↑ ↓ 移动；→ 展开，← 折叠；Enter / Space 切换。

## 限制

- 树不虚拟化，适合几十到几百个 Agent。
