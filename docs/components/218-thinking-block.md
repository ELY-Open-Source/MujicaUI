# 218 ThinkingBlock 思考过程

## 用途

可折叠的思考过程：思考中显示加载标记与实时秒数，结束后显示“思考了 n 秒”。默认折叠或展开由调用方传入的 `*open` 初值决定。基于 Collapsible。

## 最小示例

```go
mujica.ThinkingBlock(c, &open, mujica.ThinkingBlockOptions{Elapsed: 12 * time.Second}, func() {
	ui.Text(c, "Weigh the guest list.")
})
```

签名：`func ThinkingBlock(c *ui.Context, open *bool, opts ThinkingBlockOptions, content func()) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | 必填 | 是否展开；其初值即默认状态。 |
| `Thinking` | `bool` | false | 思考中：显示加载标记与从 Started 起的秒数。 |
| `Started` | `time.Time` | 零值 | 思考开始时间；Thinking 时必填，否则 panic。 |
| `Elapsed` | `time.Duration` | 0 | 思考结束后的总耗时；负数 panic。 |

## 状态

思考中 / 已完成；展开时内容带左侧细线、TextMuted 色。减少动画时加载标记静止，秒数仍每秒更新。

## 事件

`Changed()`：展开状态改变。

## 键盘操作

Tab 聚焦标题，Enter 或 Space 展开与折叠。

## 限制

- 思考中每秒重绘一次以更新秒数；折叠且完成后不再重绘。
