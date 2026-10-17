# 220 TypingIndicator 正在输入

## 用途

对方正在输入的提示：三个依次明暗变化的圆点；减少动画时改为静态文字“某某 正在输入…”。

## 最小示例

```go
chat.TypingIndicator(c, "Herald")
```

签名：`func TypingIndicator(c *ui.Context, name string) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `name` | `string` | 必填 | 正在输入的一方；为空 panic。 |

## 状态

动画 / 静态文字；状态角色，辅助技术读到“某某 正在输入…”。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 动画只在显示时请求下一帧，移出界面后停止重绘。
