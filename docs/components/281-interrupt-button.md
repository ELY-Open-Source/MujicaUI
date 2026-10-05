# 281 InterruptButton 中断按钮

## 用途

中断正在运行的智能体：可选确认，中断进行中显示“正在停止”直到运行结束。

## 最小示例

```go
running, stopping := true, false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.InterruptButton(c, running, feedback.InterruptButtonOptions{Confirm: true, Interrupting: stopping})
	if r.Changed() {
		stopping = true // cancel the run, then clear running and stopping
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func InterruptButton(c *ui.Context, running bool, opts InterruptButtonOptions) ChoiceResult`

## 参数

`InterruptButtonOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Confirm | bool | false | 点击后先弹出确认，确认后才报告中断。 |
| Interrupting | bool | false | 中断进行中：按钮保留尺寸并显示加载；running 为 false 时设为 true 会 panic。 |
| Label | string | 内置“停止” | 按钮文字，也是确认按钮的文字。 |

`running` 为 false 时按钮禁用。

## 状态

- 危险按钮外观：DangerFill 底与 OnDanger 文字，图标为叉号。
- 中断中：显示加载标记与“正在停止”，不再响应点击。
- 确认弹层默认聚焦“继续运行”，Enter 不会误触中断。

## 事件

- `Changed()`：用户要求中断（需要确认时在确认后）。

## 键盘操作

- Tab 聚焦，Enter 或 Space 按下；确认弹层中 Tab 到“停止”再 Enter，Escape 取消。

## 限制

- 只报告中断请求，取消运行由调用方完成。
