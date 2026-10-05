# 483 HoldToConfirm 按住确认

## 用途

危险操作必须按住指针或空格键一段时间才执行，按住期间填充进度；减少动画时改为显示剩余秒数。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.HoldToConfirm(c, "Delete project", input.HoldToConfirmOptions{}).Changed() {
		// delete the project
	}
}
```

签名：`func HoldToConfirm(c *ui.Context, label string, opts HoldToConfirmOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填 | 按钮文字。 |
| `Duration` | `time.Duration` | 1.5 秒 | 需按住的时长；负数 panic。 |
| `Disabled` | `bool` | `false` | 不响应。 |

## 状态

危险配色；按住时从左向右填充；减少动画时文字变为“继续按住… 1.2 秒”，宽度预留。松开或失去焦点即重置。完成后播报“已确认”。

## 事件

`Changed()`：按满时长后报告一次；同一次按住不会重复报告。

## 键盘操作

Tab 聚焦后按住空格；单击、Enter 或短按空格都不会执行。

## 限制

按住空格依赖系统送来的按键抬起；窗口失去焦点时按住状态即重置。
