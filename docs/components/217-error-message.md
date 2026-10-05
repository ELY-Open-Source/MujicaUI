# 217 ErrorMessage 错误消息

## 用途

显示失败的回复：标题、原始错误文本（等宽、可选择复制，不改写、不截断）与重试按钮；限流时显示倒计时，倒计时结束前不能重试。

## 最小示例

```go
opts := mujica.ErrorMessageOptions{RetryAt: time.Now().Add(30 * time.Second)}
if mujica.ErrorMessage(c, "429 Too Many Requests", opts).Retried() {
	fmt.Println("send again")
}
```

签名：`func ErrorMessage(c *ui.Context, text string, opts ErrorMessageOptions) ErrorMessageResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `text` | `string` | 必填 | 原始错误文本，原样显示；为空 panic。 |
| `Title` | `string` | “出错了” | 概括失败原因。 |
| `RetryAt` | `time.Time` | 零值 | 限流解除时间；之前重试按钮禁用并显示“n 秒后重试”与“已限流”。 |

## 状态

可重试 / 限流倒计时；Danger 边框、图标与标题文字共同表示错误。倒计时只在秒数变化时重绘。

## 事件

`Retried()`：按下重试。

## 键盘操作

Tab 到错误文本（可用 Shift+方向键选择、Cmd/Ctrl+C 复制）与重试按钮。

## 限制

- 倒计时依据窗口时钟；调用方的 RetryAt 应来自服务端的 Retry-After。
