# 274 ArtifactPanel 产物面板

## 用途

展示 Agent 生成的产物：代码高亮显示，文档与网页由调用方的预览插槽渲染，可在预览与源码间切换，提供复制与导出。

## 最小示例

```go
art := agent.Artifact{Title: "seating.go", Language: "go", Content: "package court\n"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if agent.ArtifactPanel(c, art, agent.ArtifactPanelOptions{}).Exported {
		fmt.Println("save the file")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ArtifactPanel(c *ui.Context, art Artifact, opts ArtifactPanelOptions) ArtifactPanelResult`。`Title` 为空、`Kind` 越界、文档或网页缺少 `Preview`、`Height` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Preview` | `func()` | `nil` | 预览插槽；文档与网页必填。 |
| `Height` | `float32` | 320 | 内容区高度。 |
| `Disabled` | `bool` | `false` | 禁止复制、导出与切换。 |

## 状态

- `Artifact` 字段：`Title`、`Kind`（`ArtifactCode`、`ArtifactDocument`、`ArtifactWeb`）、`Language`、`Content`。
- 预览与源码切换用选中底色加勾号标出当前项。
- 复制后两秒内按钮显示“已复制”。

## 事件

- `Copied`：内容已写入剪贴板。
- `Exported`：请求调用方导出。

## 键盘操作

- Tab 到切换组，← → 移动，Enter / Space 选择；Tab 到复制、导出。

## 限制

- 不内置网页渲染，网页预览须由调用方提供（无 WebView）。
