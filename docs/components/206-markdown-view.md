# 206 MarkdownView Markdown 视图

## 用途

渲染 Markdown 子集：标题、段落、粗体、斜体、删除线、行内代码、有序与无序列表（含缩进层级）、引用、链接、围栏代码块（交给 CodeBlock）、表格（交给 TableBlock）与分隔线（菱形饰线）。解析器在 `internal/markdown.go`，按块增量解析：新文本只是在旧文本后追加时，只重新解析最后一块，已完成的块保持不变。

## 最小示例

```go
if url, ok := chat.MarkdownView(c, "# Decree\n\nSee **the guide** at [the court](https://example.com).").Link(); ok {
	fmt.Println("open", url)
}
```

签名：`func MarkdownView(c *ui.Context, src string) MarkdownResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `src` | `string` | 必填 | Markdown 源文本；流式输出时每帧传入增长后的全文即可。 |

## 状态

围栏未闭合的代码块按已收到的内容显示；未闭合的强调标记按字面显示，闭合后再加样式。

## 事件

`Link()` 返回本帧被点击或用 Enter 激活的链接地址；打开链接由调用方决定。

## 键盘操作

链接可用 Tab 聚焦，Enter 激活；代码块与表格的按钮按 Tab 顺序访问。

## 限制

- 只支持上述子集：不支持 HTML、图片、脚注、任务列表、Setext 标题、链接引用定义，列表项内不支持多段落。
- 表格单元格按纯文本显示，不解析行内样式。
- 源文本被整体替换（不是在末尾追加）时重新解析全部块。
