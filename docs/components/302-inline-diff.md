# 302 InlineDiff 行内差异

## 用途

并列显示一行的旧版与新版，标出词级的删除与新增。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.InlineDiff(c, "The court meets at noon", "The court meets at dusk", git.InlineDiffOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func InlineDiff(c *ui.Context, from, to string, opts InlineDiffOptions) ui.Element`

## 参数

`InlineDiffOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Mono | bool | false | 使用等宽字体，适合代码。 |

词按字母数字、空白与单个标点切分；差异只在 from 或 to 改变时重新计算。

## 状态

- 旧行以 − 开头，删除的词带浅色底与删除线；新行以 + 开头，新增的词带浅色底与下划线。
- 两行分别以“原文”“新文”加全文朗读。

## 事件

- 无。

## 键盘操作

- 无：纯展示。

## 限制

- 只比较一行；多行请用 DiffViewer。
