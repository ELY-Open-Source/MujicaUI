# 266 AgentPlan 执行计划

## 用途

显示 Agent 的计划清单与进度：随 Agent 实时勾选，用户可勾选、跳过、恢复和新增项。

## 最小示例

```go
items := []mujica.PlanItem{
	{ID: "a", Text: "Read the guest list", State: mujica.AgentDone},
	{ID: "b", Text: "Assign the seats", State: mujica.AgentRunning},
	{ID: "c", Text: "Order the flowers"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AgentPlan(c, &items, mujica.AgentPlanOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func AgentPlan(c *ui.Context, items *[]PlanItem, opts AgentPlanOptions) ChoiceResult`。`items` 为 nil、项 `ID` 重复或状态不是待办、进行中、完成、跳过之一时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | “计划” | 标题。 |
| `Disabled` | `bool` | `false` | 只读显示，隐藏跳过与新增。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `PlanItem` 字段：`ID`、`Text`、`State`（`AgentPending`、`AgentRunning`、`AgentDone`、`AgentSkipped`）。
- 进度为“已完成 / 未跳过总数”，以文字和细进度条显示。
- 跳过的项加删除线并禁用勾选；进行中的项显示“进行中”。
- 新增项的 ID 自动生成为未使用的 `item-N`。

## 事件

- `Changed()`：用户勾选、跳过、恢复或新增了项，`*items` 已更新。

## 键盘操作

- Tab 到每项的勾选框（Space 切换）、跳过或恢复按钮；在新增输入框中按 Enter 添加。

## 限制

- 空白内容不会被添加。
