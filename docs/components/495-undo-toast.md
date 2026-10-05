# 495 UndoToast 撤销提示

## 用途

在窗口底部显示可撤销操作的消息与倒计时；撤销时调用回调。悬停或焦点在其中时暂停计时，从不抢夺焦点。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.UndoToast(c, &open, "Note deleted", mujica.UndoToastOptions{
		OnUndo: func() {
			// restore the note
		},
	})
}
```

签名：`func UndoToast(c *ui.Context, open *bool, message string, opts UndoToastOptions) UndoToastResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | 必填 | 显示与否；结束时置 false。 |
| `message` | `string` | 必填 | 文字，出现时播报。 |
| `OnUndo` | `func()` | 必填 | 撤销回调。 |
| `Duration` | `time.Duration` | 5 秒 | 倒计时；负数 panic。 |
| `Action` | `string` | “Undo” | 按钮文字。 |

## 状态

按钮显示剩余秒数，下方细条随时间缩短；减少动画时不画细条，只更新秒数。

## 事件

`Undone`：已撤销（回调已执行）；`Expired`：倒计时结束或用户关闭，操作生效。

## 键盘操作

Tab 到“撤销”与关闭按钮；焦点在其中时 Escape 关闭。

## 限制

每个窗口同时只显示一条。
