# 259 ArtifactCard · 产物卡片

## 用途

消息中的产物卡片（代码、文档、网页），点击打开产物面板。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.ArtifactCard(c, chat.ArtifactCardOptions{Title: "seating.go", Detail: "Go · 120 lines"}).Changed() {
		log.Println("open artifact panel")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ArtifactCard(c *ui.Context, opts ArtifactCardOptions) ChoiceResult`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Kind` | `ArtifactKind` | `ArtifactCode` | `ArtifactCode`、`ArtifactDocument`、`ArtifactWeb`；越界 panic。 |
| `Title` | `string` | 必填 | 产物标题。 |
| `Detail` | `string` | "" | 简短说明，如“Go · 120 行”。 |
| `Streaming` | `bool` | false | 仍在生成，显示旋转标记。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、悬停、聚焦、生成中、禁用。种类同时用图标与文字表示。

## 事件

`Changed()`：请求打开产物面板。

## 键盘操作

- Tab 聚焦；Enter/Space 打开。

## 限制

- 卡片不渲染产物内容，面板由调用方实现。
