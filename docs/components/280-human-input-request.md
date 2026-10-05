# 280 HumanInputRequest 人工输入请求

## 用途

Agent 向用户提问：显示问题、可选答案与自由回答，并提示剩余时间与超时。

## 最小示例

```go
answer := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := agent.HumanInputRequest(c, &answer, agent.HumanInputRequestOptions{
		Question: "Which hall?", Options: []string{"East wing", "Great hall"}, AllowFree: true,
	})
	if r.Submitted {
		fmt.Println("reply", answer)
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func HumanInputRequest(c *ui.Context, answer *string, opts HumanInputRequestOptions) HumanInputResult`。`answer` 为 nil、`Question` 为空、既无选项又不允许自由回答时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Question` | `string` | 必填 | 问题。 |
| `Options` | `[]string` | `nil` | 可选答案。 |
| `AllowFree` | `bool` | `false` | 允许自由回答。 |
| `Deadline` | `time.Time` | 零值 | 截止时间；零值不限时。 |
| `Disabled` | `bool` | `false` | 禁止回答。 |

## 状态

- 选择选项会清空自由回答；输入自由回答会取代所选选项。
- 剩余时间每秒刷新；到时显示“智能体已停止等待”并禁用。
- 回答为空白时“发送回答”禁用。

## 事件

- `Submitted`：用户发送了 `*answer`。
- `TimedOut`：截止时间在本帧到达（只报告一次）。

## 键盘操作

- Tab 进入选项组，方向键选择；Tab 到输入框，Enter 发送；或 Tab 到“发送回答”。

## 限制

- 超时只在界面帧中检测，窗口隐藏期间推迟到重新显示。
