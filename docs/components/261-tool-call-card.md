# 261 ToolCallCard 工具调用卡片

## 用途

显示一次工具调用：工具名、状态、耗时，参数与结果可折叠，JSON 自动格式化并高亮。

## 最小示例

```go
call := mujica.ToolCall{ID: "1", Name: "read_file", Args: "{\"path\":\"ledger.md\"}",
	Result: "{\"lines\":40}", State: mujica.AgentDone, Duration: 420 * time.Millisecond}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ToolCallCard(c, call, mujica.ToolCallCardOptions{OpenArgs: true})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ToolCallCard(c *ui.Context, call ToolCall, opts ToolCallCardOptions) ChoiceResult`。`call.Name` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OpenArgs` | `bool` | `false` | 参数区初始展开。 |
| `OpenResult` | `bool` | `false` | 结果区初始展开。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `ToolCall` 字段：`ID`、`Name`、`Args`（JSON 或纯文本）、`Result`、`Error`（非空时替代结果并以错误色显示）、`State`、`Duration`。
- `AgentState`：`AgentPending`、`AgentPlanning`、`AgentRunning`、`AgentWaiting`、`AgentDone`、`AgentFailed`、`AgentSkipped`，越界 panic。状态同时由图标形状和文字表达，不只靠颜色。
- 仅完成或失败时显示耗时；合法 JSON 以两格缩进格式化，非 JSON 原样显示。
- 结果超过 200 行时虚拟化，只构建可见行。

## 事件

- `Changed()`：用户展开或折叠参数、结果。

## 键盘操作

- Tab 到“参数”“结果”折叠按钮，Enter / Space 展开或折叠。

## 限制

- 展开状态保存在卡片元素上；在列表中请用稳定 `Key` 包裹每张卡片。
- 高亮只区分关键字、字符串、数字和标点，不做 JSON 结构校验提示。
