# 250 FeedbackForm · 反馈表单

## 用途

在点赞或点踩之后询问原因：原因标签多选与补充文字，然后提交或取消。

## 最小示例

```go
var value chat.FeedbackValue
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.FeedbackForm(c, &value, chat.FeedbackFormOptions{Rating: chat.FeedbackDown}).Submitted {
		log.Println("reasons", value.Reasons)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func FeedbackForm(c *ui.Context, value *FeedbackValue, opts FeedbackFormOptions) FeedbackFormView`。`FeedbackValue`：`Reasons`、`Comment`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Rating` | `FeedbackRating` | `FeedbackUp` | 赞（`FeedbackUp`）或踩（`FeedbackDown`），决定标题与内置原因。 |
| `Reasons` | `[]string` | nil（内置） | 替换内置原因。 |
| `Submitting` | `bool` | false | 调用方正在提交：提交按钮显示加载，表单不可编辑。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

原因标签（按下：Selection 底 + 下划条）、文字框、提交中、禁用。

## 事件

`Submitted`：提交；`Cancelled`：取消或 Escape。

## 键盘操作

- Tab 依次经过原因组、文字框、取消、提交；原因组内 ←/→；Space 切换原因。
- Escape：取消。

## 限制

- 表单不发送任何数据；提交由调用方完成。
