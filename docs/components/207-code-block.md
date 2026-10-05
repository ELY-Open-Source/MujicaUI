# 207 CodeBlock 代码块

## 用途

显示代码：语言标签、可选行号、基础语法着色（`internal.Highlight`，支持 go、javascript、typescript、python、json、shell、sql），复制、折叠与自动换行切换。

## 最小示例

```go
chat.CodeBlock(c, "func bow() {}", chat.CodeBlockOptions{Language: "go", LineNumbers: true})
```

签名：`func CodeBlock(c *ui.Context, code string, opts CodeBlockOptions) CopyResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Language` | `string` | ""（显示“文本”） | 着色语言；未知语言按纯文本显示。 |
| `LineNumbers` | `bool` | false | 显示右对齐的行号。 |
| `Wrap` | `bool` | false | 自动换行的初始状态；关闭时横向滚动。 |
| `Collapsed` | `bool` | false | 折叠的初始状态。 |

## 状态

换行开关以 Selection 底色加边框表示开启；复制后按钮在 2 秒内显示对勾“已复制”；折叠时只显示标题栏，展开按钮写明行数。关键字、字符串、数字、类型使用不同颜色，注释同时使用斜体。

## 事件

`Copied()`：本帧把代码写入了剪贴板。

## 键盘操作

Tab 依次到换行开关、复制、折叠按钮；Space 或 Enter 触发。

## 限制

- 着色是逐行的词法着色，不做语法分析；块注释可跨行。
- 不做虚拟化，适合消息中的代码片段，不适合上万行文件。
