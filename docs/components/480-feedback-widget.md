# 480 FeedbackWidget

反馈组件。

## 用途

选择反馈类型、填写内容、附件插槽，提交一次后等待调用方完成。基于 Form。

## 最小示例

```go
var state mujica.FeedbackWidgetState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if f, ok := mujica.FeedbackWidget(c, &state, mujica.FeedbackWidgetOptions{}).Submitted(); ok {
		fmt.Println("send", f.Kind)
		state.EndSubmit("")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*FeedbackWidgetState` | 必填 | Draft 与 Form。 |
| `Attachments` | `func()` | nil | 截图或文件插槽。 |
| `MinLength` | `int` | 10 | 最少字数。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

填写、过短错误、提交中、失败、禁用。提交中锁定类型、正文与附件控件，失败保留原稿。

## 事件

`Submitted()` 返回反馈；完成后调用 `s.EndSubmit(err)`，成功时清空内容。

## 键盘操作

类型分段用方向键切换；Tab 到内容与按钮。

## 限制

- 截图采集由调用方完成。
