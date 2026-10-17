# 205 StreamingText 流式文本

## 用途

显示随流式输出不断增长的文本。可平滑出字（按速率逐字显示），可显示流式光标；减少动画时直接显示全部文本。文本被替换为不以旧文本开头的内容时，从共同前缀继续。

## 最小示例

```go
chat.StreamingText(c, reply, chat.StreamingTextOptions{Streaming: true, Smooth: true})
```

签名：`func StreamingText(c *ui.Context, text string, opts StreamingTextOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `text` | `string` | 必填 | 当前收到的全部文本。 |
| `Streaming` | `bool` | false | 在文本末尾显示 AccentText 色的静态光标 ▍。 |
| `Smooth` | `bool` | false | 新到的文本按 Rate 逐字显示；关闭时立即显示。 |
| `Rate` | `float64` | 80 | 平滑出字速度，字符/秒；负数 panic。 |

## 状态

出字中 / 完成；减少动画时始终为完成。光标不闪烁，避免装饰性循环。

## 事件

无。

## 键盘操作

不获取焦点。辅助技术读到的名称始终是完整文本。

## 限制

- 只渲染纯文本；Markdown 内容请交给 MarkdownView（206），它的增量解析同样适合流式输出。
