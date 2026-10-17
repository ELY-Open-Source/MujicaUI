# 219 ScrollToBottomButton 回到底部按钮

## 用途

离开底部时出现在容器底部居中的圆角按钮，带未读计数；点击回到底部。MessageList 已内置，也可单独用于其他滚动区域。

## 最小示例

```go
ui.Box(c).Fill().Children(func() {
	if chat.ScrollToBottomButton(c, !list.AtEnd(), 3).Clicked() {
		list.ScrollToEnd()
	}
})
```

签名：`func ScrollToBottomButton(c *ui.Context, away bool, unread int) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `away` | `bool` | 必填 | 是否已离开底部；false 时返回不可见的空元素，其 Clicked 为 false。 |
| `unread` | `int` | 0 | 未读数，大于 0 时显示；负数 panic。 |

## 状态

出现时 160 毫秒淡入（减少动画时直接显示）；悬停时底色加深。

## 事件

`Clicked()`：点击或 Enter、Space；滚动由调用方完成。

## 键盘操作

Tab 聚焦，Enter 或 Space 触发。辅助技术读到“回到最新，n 条新消息”。

## 限制

- 按钮以绝对定位贴在父元素底部，父元素需要占满滚动区域。
