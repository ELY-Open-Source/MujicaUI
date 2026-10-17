# 098 StatusIndicator 状态标记

## 用途

以图形加文字标出状态：成功、警告、错误、离线、处理中。每种状态有独立形状（对勾圆、三角、叉圆、断开的 Wi-Fi、旋转菱形），并总是显示文字，不只靠颜色。

## 最小示例

```go
feedback.StatusIndicator(c, feedback.StatusSuccess, feedback.StatusIndicatorOptions{Label: "API online"})
feedback.StatusIndicator(c, feedback.StatusOffline, feedback.StatusIndicatorOptions{})
```

签名：`func StatusIndicator(c *ui.Context, status Status, opts StatusIndicatorOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `status` | `Status` | `StatusSuccess` | `StatusSuccess`、`StatusWarning`、`StatusError`、`StatusOffline`、`StatusBusy`；未知值 panic。 |
| `Label` | `string` | 内置词（Healthy、Degraded、Failing、Offline、Working） | 替换显示的文字。 |

## 状态

处理中使用 Spinner，减少动画时静止。整体为 `RoleStatus`，名称为文字。

## 事件

无。

## 键盘操作

不可聚焦。

## 限制

- 无 Badge 式的纯圆点模式；PRD 要求不能只靠颜色。
