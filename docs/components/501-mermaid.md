# 501 Mermaid Mermaid 图

## 用途

把 Mermaid 源码画成图：流程图、时序图、饼图等 33 种。颜色取自窗口主题，文字用绘制它的字体量宽。

## 最小示例

```go
source := "flowchart LR\n  Score[Write score] --> Rehearse[Rehearse]\n  Rehearse --> Stage([Take the stage])"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	diagram.Mermaid(c, source, diagram.MermaidOptions{Label: "Concert plan"})
}
```

签名：`func Mermaid(c *ui.Context, source string, opts MermaidOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `source` | `string` | 必填 | Mermaid 源码。 |
| `Label` | `string` | 源码的 `accTitle`，再无则“图表” | 读给辅助技术的名称。 |

## 状态

图按自身尺寸显示，容器窄时等比缩小，不放大。源码有错时显示错误提示：Pele 的原文消息，前面带行号与列号。

## 事件

无。

## 键盘操作

无；图是一张静态图片，角色为 image。

## 限制

引擎是 pele-go v0.2.0（Pele 1.10.0 的 Go 移植）。链接、远程图片与数学公式不绘制：`click` 不产生链接，`$$…$$` 按原文显示。只有导入 `diagram` 的程序才链接引擎，体积增加约 6 MB。场景在源码、主题色或字号变化时才重算。
